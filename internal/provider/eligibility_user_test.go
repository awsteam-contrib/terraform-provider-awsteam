package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEligibilityUserResource_basic(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_user.test"
	user1 := gofakeit.Email()
	user2 := gofakeit.Email()
	userId1 := gofakeit.UUID()
	userId2 := gofakeit.UUID()
	approval1 := true
	approval2 := false
	duration := fmt.Sprint(gofakeit.Number(1, 10))
	ticketNo := gofakeit.BS()
	accountId := gofakeit.DigitN(12)
	accountName := gofakeit.BS()
	ouId := "ou-cxt3-2782ty5g" // hard coded fake ou id
	ouName := gofakeit.BS()
	permissionArn := "arn:aws:sso:::permissionSet/ssoins-4334d1f197f50907/ps-f5ge203d3d2428d3" // hard coded fake arn
	permissionName := "elevated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEligibilityUserResourceConfig(user1, userId1, approval1, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_name", user1),
					resource.TestCheckResourceAttr(resourceName, "approval_required", fmt.Sprint(approval1)),
					resource.TestCheckResourceAttr(resourceName, "duration", duration),
					resource.TestCheckResourceAttr(resourceName, "ticket_no", ticketNo),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "accounts.*",
						map[string]string{
							"account_id":   accountId,
							"account_name": accountName,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ous.*",
						map[string]string{
							"ou_id":   ouId,
							"ou_name": ouName,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "permissions.*",
						map[string]string{
							"permission_arn":  permissionArn,
							"permission_name": permissionName,
						}),
					// resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupId1),
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
				Config: testAccEligibilityUserResourceConfig(user2, userId2, approval2, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", userId2),
					resource.TestCheckResourceAttr(resourceName, "user_name", user2),
					resource.TestCheckResourceAttr(resourceName, "approval_required", fmt.Sprint(approval2)),
				),
			},
		},
	})
}

func TestAccEligibilityUserResource_missingAccountsAndOUs(t *testing.T) {
	user1 := gofakeit.Email()
	userId1 := gofakeit.UUID()
	approval1 := true
	duration := fmt.Sprint(gofakeit.Number(1, 10))
	ticketNo := gofakeit.BS()
	permissionArn := "arn:aws:sso:::permissionSet/ssoins-4334d1f197f50907/ps-f5ge203d3d2428d3" // hard coded fake arn
	permissionName := "elevated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEligibilityUserMissingAccountsAndOUsResourceConfig(user1, userId1, approval1, duration, ticketNo, permissionArn, permissionName),
				ExpectError: regexp.MustCompile(
					"At Least One Account or OU Must Be Specified.",
				),
			},
			{
				Config: testAccEligibilityUserEmptyAccountsAndOUsResourceConfig(user1, userId1, approval1, duration, ticketNo, permissionArn, permissionName),
				ExpectError: regexp.MustCompile(
					"At Least One Account or OU Must Be Specified.",
				),
			},
		},
	})
}

func TestAccEligibilityUserResource_missingAccounts(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_user.test"
	user1 := gofakeit.Email()
	userId1 := gofakeit.UUID()
	approval1 := true
	duration := fmt.Sprint(gofakeit.Number(1, 10))
	ticketNo := gofakeit.BS()
	ouId := "ou-cxt3-2782ty5g" // hard coded fake ou id
	ouName := gofakeit.BS()
	permissionArn := "arn:aws:sso:::permissionSet/ssoins-4334d1f197f50907/ps-f5ge203d3d2428d3" // hard coded fake arn
	permissionName := "elevated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEligibilityUserMissingAccountsResourceConfig(user1, userId1, approval1, duration, ticketNo, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_name", user1),
					resource.TestCheckResourceAttr(resourceName, "approval_required", fmt.Sprint(approval1)),
					resource.TestCheckResourceAttr(resourceName, "duration", duration),
					resource.TestCheckResourceAttr(resourceName, "ticket_no", ticketNo),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ous.*",
						map[string]string{
							"ou_id":   ouId,
							"ou_name": ouName,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "permissions.*",
						map[string]string{
							"permission_arn":  permissionArn,
							"permission_name": permissionName,
						}),
					// resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupId1),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEligibilityUserResource_missingOUs(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_user.test"
	user1 := gofakeit.Email()
	userId1 := gofakeit.UUID()
	approval1 := true
	duration := fmt.Sprint(gofakeit.Number(1, 10))
	ticketNo := gofakeit.BS()
	accountId := gofakeit.DigitN(12)
	accountName := gofakeit.BS()
	permissionArn := "arn:aws:sso:::permissionSet/ssoins-4334d1f197f50907/ps-f5ge203d3d2428d3" // hard coded fake arn
	permissionName := "elevated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEligibilityUserMissingOUsResourceConfig(user1, userId1, approval1, duration, ticketNo, accountId, accountName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_name", user1),
					resource.TestCheckResourceAttr(resourceName, "approval_required", fmt.Sprint(approval1)),
					resource.TestCheckResourceAttr(resourceName, "duration", duration),
					resource.TestCheckResourceAttr(resourceName, "ticket_no", ticketNo),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "accounts.*",
						map[string]string{
							"account_id":   accountId,
							"account_name": accountName,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "permissions.*",
						map[string]string{
							"permission_arn":  permissionArn,
							"permission_name": permissionName,
						}),
					// resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupId1),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEligibilityUserResource_emptyAccounts(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_user.test"
	user1 := gofakeit.Email()
	userId1 := gofakeit.UUID()
	approval1 := true
	duration := fmt.Sprint(gofakeit.Number(1, 10))
	ticketNo := gofakeit.BS()
	ouId := "ou-cxt3-2782ty5g" // hard coded fake ou id
	ouName := gofakeit.BS()
	permissionArn := "arn:aws:sso:::permissionSet/ssoins-4334d1f197f50907/ps-f5ge203d3d2428d3" // hard coded fake arn
	permissionName := "elevated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEligibilityUserEmptyAccountsResourceConfig(user1, userId1, approval1, duration, ticketNo, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_name", user1),
					resource.TestCheckResourceAttr(resourceName, "approval_required", fmt.Sprint(approval1)),
					resource.TestCheckResourceAttr(resourceName, "duration", duration),
					resource.TestCheckResourceAttr(resourceName, "ticket_no", ticketNo),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "ous.*",
						map[string]string{
							"ou_id":   ouId,
							"ou_name": ouName,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "permissions.*",
						map[string]string{
							"permission_arn":  permissionArn,
							"permission_name": permissionName,
						}),
					// resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupId1),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEligibilityUserResource_emptyOUs(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_user.test"
	user1 := gofakeit.Email()
	userId1 := gofakeit.UUID()
	approval1 := true
	duration := fmt.Sprint(gofakeit.Number(1, 10))
	ticketNo := gofakeit.BS()
	accountId := gofakeit.DigitN(12)
	accountName := gofakeit.BS()
	permissionArn := "arn:aws:sso:::permissionSet/ssoins-4334d1f197f50907/ps-f5ge203d3d2428d3" // hard coded fake arn
	permissionName := "elevated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEligibilityUserEmptyOUsResourceConfig(user1, userId1, approval1, duration, ticketNo, accountId, accountName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_id", userId1),
					resource.TestCheckResourceAttr(resourceName, "user_name", user1),
					resource.TestCheckResourceAttr(resourceName, "approval_required", fmt.Sprint(approval1)),
					resource.TestCheckResourceAttr(resourceName, "duration", duration),
					resource.TestCheckResourceAttr(resourceName, "ticket_no", ticketNo),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "accounts.*",
						map[string]string{
							"account_id":   accountId,
							"account_name": accountName,
						}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "permissions.*",
						map[string]string{
							"permission_arn":  permissionArn,
							"permission_name": permissionName,
						}),
					// resource.TestCheckTypeSetElemAttr(resourceName, "group_ids.*", groupId1),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEligibilityUserResource_disappears(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_user.test"
	user1 := gofakeit.Email()
	userId1 := gofakeit.UUID()
	approval1 := true
	duration := fmt.Sprint(gofakeit.Number(1, 10))
	ticketNo := gofakeit.BS()
	accountId := gofakeit.DigitN(12)
	accountName := gofakeit.BS()
	ouId := "ou-cxt3-2782ty5g" // hard coded fake ou id
	ouName := gofakeit.BS()
	permissionArn := "arn:aws:sso:::permissionSet/ssoins-4334d1f197f50907/ps-f5ge203d3d2428d3" // hard coded fake arn
	permissionName := "elevated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEligibilityUserResourceConfig(user1, userId1, approval1, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					testAccEligibilityResourceDisappears(ctx, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccEligibilityUserResourceConfig(user string, userId string, approvalRequired bool, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_user" "test" {
	user_name         = "%s"
	user_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	accounts = [
		{
		account_id   = "%s"
		account_name = "%s"
		}
	]
	ous = [
		{
		ou_id   = "%s"
		ou_name = "%s"
		}
	]
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, user, userId, approvalRequired, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName)
}

func testAccEligibilityUserMissingAccountsAndOUsResourceConfig(user string, userId string, approvalRequired bool, duration, ticketNo, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_user" "test" {
	user_name         = "%s"
	user_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, user, userId, approvalRequired, duration, ticketNo, permissionArn, permissionName)
}

func testAccEligibilityUserEmptyAccountsAndOUsResourceConfig(user string, userId string, approvalRequired bool, duration, ticketNo, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_user" "test" {
	user_name         = "%s"
	user_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	accounts = []
	ous = []
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, user, userId, approvalRequired, duration, ticketNo, permissionArn, permissionName)
}

func testAccEligibilityUserMissingAccountsResourceConfig(user string, userId string, approvalRequired bool, duration, ticketNo, ouId, ouName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_user" "test" {
	user_name         = "%s"
	user_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	ous = [
		{
		ou_id   = "%s"
		ou_name = "%s"
		}
	]
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, user, userId, approvalRequired, duration, ticketNo, ouId, ouName, permissionArn, permissionName)
}

func testAccEligibilityUserMissingOUsResourceConfig(user string, userId string, approvalRequired bool, duration, ticketNo, accountId, accountName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_user" "test" {
	user_name         = "%s"
	user_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	accounts = [
		{
		account_id   = "%s"
		account_name = "%s"
		}
	]
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, user, userId, approvalRequired, duration, ticketNo, accountId, accountName, permissionArn, permissionName)
}

func testAccEligibilityUserEmptyAccountsResourceConfig(user string, userId string, approvalRequired bool, duration, ticketNo, ouId, ouName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_user" "test" {
	user_name         = "%s"
	user_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	accounts = []
	ous = [
		{
		ou_id   = "%s"
		ou_name = "%s"
		}
	]
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, user, userId, approvalRequired, duration, ticketNo, ouId, ouName, permissionArn, permissionName)
}

func testAccEligibilityUserEmptyOUsResourceConfig(user string, userId string, approvalRequired bool, duration, ticketNo, accountId, accountName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_user" "test" {
	user_name         = "%s"
	user_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	accounts = [
		{
		account_id   = "%s"
		account_name = "%s"
		}
	]
	ous = []
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, user, userId, approvalRequired, duration, ticketNo, accountId, accountName, permissionArn, permissionName)
}
