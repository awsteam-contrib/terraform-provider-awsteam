resource "awsteam_settings" "example" {
  duration           = 5
  expiry             = 3
  team_admin_group   = "My-Team-Admin-Group"
  team_auditor_group = "My-Team-Auditor-Group"
}

// Import the existing settings on a fresh install of AWS TEAM
import {
  to = awsteam_settings.example
  id = "settings"
}

// Enable OU caching to improve performance when resolving AWS Organizations OU data
resource "awsteam_settings" "example_with_ou_cache" {
  duration           = 5
  expiry             = 3
  team_admin_group   = "My-Team-Admin-Group"
  team_auditor_group = "My-Team-Auditor-Group"
  use_ou_cache       = true
}

