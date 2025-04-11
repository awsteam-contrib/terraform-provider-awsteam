package provider

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// AWS TEAM Only allows one settings to be defined at a time.
func TestAccSettings_serial(t *testing.T) {
	t.Parallel()

	// AWS TEAM only allows settings to be defined once, Running these tests requires the settings to be deleted from the environment before running.
	// Added an environment variable that will enable us to control when these tests run.
	// To run these tests set the environment variable like so: export AWSTEAM_RUN_SETTINGS_TESTS="true"
	key := "AWSTEAM_RUN_SETTINGS_TESTS"
	vifId := os.Getenv(key)
	if vifId != "true" {
		t.Skipf("Skipping Settings Tests, Environment variable %s is not set to true", key)
	}

	testCases := map[string]map[string]func(t *testing.T){
		"Resource": {
			"basic":    testAccSettingsResource_basic,
			"duration": testAccSettingsResource_duration,
		},
		"DataSource": {
			"basic": testAccSettingsDataSource_basic,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}

func testAccSettingsResource_basic(t *testing.T) {
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

func testAccSettingsResourceConfig(teamAdminGroup string, teamAuditorGroup string, duration int, expiry int) string {
	return fmt.Sprintf(`
resource "awsteam_settings" "test" {
  team_admin_group = %[1]q
  team_auditor_group = %[2]q
  duration = %d
  expiry = %d
}
`, teamAdminGroup, teamAuditorGroup, duration, expiry)
}

func testAccSettingsResourceConfigDuration(duration int) string {
	return fmt.Sprintf(`
resource "awsteam_settings" "test" {
  team_admin_group = "Team-Admin-Group"
  team_auditor_group = "Team-Auditor-Group"
  duration = %d
  expiry = 5
}
`, duration)
}
