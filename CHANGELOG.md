# CHANGELOG

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/) and [Keep a Changelog](http://keepachangelog.com/).

## Unreleased

### New

### Changes

### Fixes

### Breaks

## 1.1.2 - (2025-04-16)

### Fixes

* Resource: `awsteam_eligibility_group` - Validation rule to ensure a single value is provided for `accounts` or `ous` now allows `unknown` values.
* Resource: `awsteam_eligibility_user` - Validation rule to ensure a single value is provided for `accounts` or `ous` now allows `unknown` values.

## 1.1.1 - (2025-04-15)

### Changes

* Resource: `awsteam_eligibility_group` - The `accounts` and `ous` fields have been made *optional*. A configuration validation rule now verifies that at least one `account` or `ou` is provided. It is no longer required to pass an empty list for `accounts` or `ous` as long as you provide at least one value for either of these fields. [70](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/70)
* Resource: `awsteam_eligibility_user` - The `accounts` and `ous` fields have been made *optional*. A configuration validation rule now verifies that at least one `account` or `ou` is provided. It is no longer required to pass an empty list for `accounts` or `ous` as long as you provide at least one value for either of these fields. [70](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/70)

## 1.1.0 - (2024-03-26)

### New

* DataSource: `awsteam_accounts` [#44](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/44)

### Fixes

* Resource: `awsteam_settings` - resolved update flow failing to set computed values [#41](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/41)

## 1.0.1 - (2024-03-05)
---

### Fixes

* Resource: `awsteam_eligibility_group` now requires the `group_id` field. This will be used as the resource `id`. [36](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/36)
* Resource: `awsteam_eligibility_user` now requires the `user_id` field. This will be used as the resource `id`. [36](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/36)
* Resources: `awsteam_eligibility_group`, `awsteam_eligibility_user`, `awsteam_approvers_account`, `awsteam_approvers_ou`, and `awsteam_settings` Deleting outside of terraform will no longer cause terraform to error.

## 1.0.0 - (2024-02-27)

### Changes

* Removed validators on `group_id` schema fields as they can be values other than UUID [26](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/26)

### Fixes

* Provider: Corrected the provider schema attribute mapping to configuration fields for `client_secret`, `graph_endpoint`, and `token_endpoint` [24](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/24)

### Breaks

* Resources: `awsteam_eligibility_group`, `awsteam_eligibility_user`, and `awsteam_approvers_account` AWS account numbers are now string values [26](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/26)


## 0.2.0 - (2024-01-17)

### New

* Resource: `awsteam_eligibility_group` [#8](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/8)
* Resource: `awsteam_eligibility_user` [#8](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/8)

## 0.1.0 - (2024-01-12)

### New

* Resource: `awsteam_approvers_account` [#7](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/7)
* Resource: `awsteam_approvers_ou` [#7](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/7)
* DataSource: `awsteam_settings` [#4](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/4)
* Resource: `awsteam_settings` [#4](https://github.com/awsteam-contrib/terraform-provider-awsteam/issues/4)
