package provider

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"testing"

	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/acctest"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSettings_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":        testAccSettingsResource_basic,
			"duration":     testAccSettingsResource_duration,
			"use_ou_cache": testAccSettingsResource_useOUCache,
		},
		"DataSource": {
			"basic": testAccSettingsDataSource_basic,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}

// testAccSettingsSaveAndRestore reads any existing settings, registers a t.Cleanup
// to recreate them after the test's terraform destroy runs, and then deletes them so
// the test's Create call succeeds (TEAM only permits one settings record at a time).
// Call this at the top of every settings acceptance test.
func testAccSettingsSaveAndRestore(t *testing.T) {
	t.Helper()

	ctx := context.Background()
	client := acctest.NewAWSTeamClient(ctx)

	existing, err := client.GetSettings(ctx, &awsteam.GetSettingsInput{})
	if err != nil {
		t.Fatalf("failed to read existing settings before test: %s", err)
	}

	if existing == nil || existing.Settings == nil {
		return
	}

	s := existing.Settings

	t.Cleanup(func() {
		_, restoreErr := client.CreateSettings(ctx, &awsteam.CreateSettingsInput{
			Approval:                  s.Approval,
			Comments:                  s.Comments,
			Duration:                  s.Duration,
			Expiry:                    s.Expiry,
			SesNotificationsEnabled:   s.SesNotificationsEnabled,
			SnsNotificationsEnabled:   s.SnsNotificationsEnabled,
			SlackNotificationsEnabled: s.SlackNotificationsEnabled,
			SesSourceEmail:            s.SesSourceEmail,
			SesSourceArn:              s.SesSourceArn,
			SlackToken:                s.SlackToken,
			TeamAdminGroup:            s.TeamAdminGroup,
			TeamAuditorGroup:          s.TeamAuditorGroup,
			TicketNo:                  s.TicketNo,
			UseOUCache:                s.UseOUCache,
			ModifiedBy:                s.ModifiedBy,
		})
		if restoreErr != nil {
			t.Logf("warning: failed to restore settings after test: %s", restoreErr)
		}
	})

	if _, delErr := client.DeleteSettings(ctx, &awsteam.DeleteSettingsInput{}); delErr != nil {
		t.Fatalf("failed to delete existing settings before test: %s", delErr)
	}
}

func testAccSettingsResource_basic(t *testing.T) {
	testAccSettingsSaveAndRestore(t)

	resourceName := "awsteam_settings.test"
	teamAdminGroup1 := "Team-Admin-Group"
	teamAuditorGroup1 := "Team-Auditor-Group"
	teamAdminGroup2 := "Team-Admin-Group"
	teamAuditorGroup2 := "Team-Auditor-Group"
	duration := rand.Intn(10)
	expiry := rand.Intn(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfig(teamAdminGroup1, teamAuditorGroup1, duration, expiry),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "team_admin_group", teamAdminGroup1),
					resource.TestCheckResourceAttr(resourceName, "team_auditor_group", teamAuditorGroup1),
					resource.TestCheckResourceAttr(resourceName, "approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "comments", "false"),
					resource.TestCheckResourceAttr(resourceName, "ses_notifications_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "sns_notifications_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "slack_notifications_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ticket_no", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSettingsResourceConfig(teamAdminGroup2, teamAuditorGroup2, duration, expiry),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "team_admin_group", teamAdminGroup2),
					resource.TestCheckResourceAttr(resourceName, "team_auditor_group", teamAuditorGroup2)),
			},
		},
	})
}

func testAccSettingsResource_duration(t *testing.T) {
	testAccSettingsSaveAndRestore(t)

	resourceName := "awsteam_settings.test"
	duration := rand.Intn(10)
	duration2 := rand.Intn(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigDuration(duration),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "duration", fmt.Sprint(duration)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSettingsResourceConfigDuration(duration2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "duration", fmt.Sprint(duration2)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
		},
	})
}

func testAccSettingsResource_useOUCache(t *testing.T) {
	testAccSettingsSaveAndRestore(t)

	ctx := context.Background()
	client := acctest.NewAWSTeamClient(ctx)

	// When the target environment does not support the use_ou_cache setting
	// verify we get the expected error
	if !client.SettingsCapabilities.UseOUCacheSupported {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccSettingsResourceConfigUseOUCache(true),
					ExpectError: regexp.MustCompile(`use_ou_cache requires AWS TEAM`),
				},
			},
		})
		return
	}

	resourceName := "awsteam_settings.test"

	// When the target environment does support the use_ou_cache setting
	// verify functionality
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigUseOUCache(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "use_ou_cache", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSettingsResourceConfigUseOUCache(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "use_ou_cache", "false"),
				),
			},
		},
	})
}

func testAccSettingsResourceConfig(teamAdminGroup string, teamAuditorGroup string, duration int, expiry int) string {
	return fmt.Sprintf(`
resource "awsteam_settings" "test" {
  team_admin_group   = %[1]q
  team_auditor_group = %[2]q
  duration           = %d
  expiry             = %d
}
`, teamAdminGroup, teamAuditorGroup, duration, expiry)
}

func testAccSettingsResourceConfigDuration(duration int) string {
	return fmt.Sprintf(`
resource "awsteam_settings" "test" {
  team_admin_group   = "Team-Admin-Group"
  team_auditor_group = "Team-Auditor-Group"
  duration           = %d
  expiry             = 5
}
`, duration)
}

func testAccSettingsResourceConfigUseOUCache(useOUCache bool) string {
	return fmt.Sprintf(`
resource "awsteam_settings" "test" {
  team_admin_group   = "Team-Admin-Group"
  team_auditor_group = "Team-Auditor-Group"
  duration           = 5
  expiry             = 5
  use_ou_cache       = %t
}
`, useOUCache)
}
