package provider

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccSettingsDataSource_basic(t *testing.T) {
	testAccSettingsSaveAndRestore(t)

	resourceName := "awsteam_settings.test"
	dataSourceName := "data.awsteam_settings.test"
	teamAdminGroup := "Team-Admin-Group"
	teamAuditorGroup := "Team-Auditor-Group"
	duration := rand.Intn(10)
	expiry := rand.Intn(10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsDataSourceConfig(teamAdminGroup, teamAuditorGroup, duration, expiry),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "team_admin_group", dataSourceName, "team_admin_group"),
					resource.TestCheckResourceAttrPair(resourceName, "team_auditor_group", dataSourceName, "team_auditor_group"),
					resource.TestCheckResourceAttrPair(resourceName, "approval", dataSourceName, "approval"),
					resource.TestCheckResourceAttrPair(resourceName, "comments", dataSourceName, "comments"),
					resource.TestCheckResourceAttrPair(resourceName, "ses_notifications_enabled", dataSourceName, "ses_notifications_enabled"),
					resource.TestCheckResourceAttrPair(resourceName, "sns_notifications_enabled", dataSourceName, "sns_notifications_enabled"),
					resource.TestCheckResourceAttrPair(resourceName, "slack_notifications_enabled", dataSourceName, "slack_notifications_enabled"),
					resource.TestCheckResourceAttrPair(resourceName, "ticket_no", dataSourceName, "ticket_no"),
					resource.TestCheckResourceAttrPair(resourceName, "created_at", dataSourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "updated_at", dataSourceName, "updated_at"),
				),
			},
		},
	})
}
func testAccSettingsDataSourceConfig(teamAdminGroup string, teamAuditorGroup string, duration int, expiry int) string {
	return fmt.Sprintf(`
resource "awsteam_settings" "test" {
	team_admin_group = %[1]q
	team_auditor_group = %[2]q
	duration = %d
	expiry = %d
  }

data "awsteam_settings" "test" {
	depends_on = [
		awsteam_settings.test
	]
}
`, teamAdminGroup, teamAuditorGroup, duration, expiry)
}
