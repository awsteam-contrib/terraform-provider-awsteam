data "awsteam_accounts" "all" {}

// Build a map of Account Name to Account id
locals {
  account_name_to_id_map = { for account in data.awsteam_accounts.all.accounts : account.name => account.id }
}

// How to access the account id from the account name mapping
output "test_account_id" {
  value = local.account_name_to_id_map["my-test-aws-account-name"]
}

// Access all account names from the data source
output "account_names" {
  value = data.awsteam_accounts.all.accounts.*.name
}

// Access all account ids from the data source
output "account_ids" {
  value = data.awsteam_accounts.all.accounts.*.id
}
