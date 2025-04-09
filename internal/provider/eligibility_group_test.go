package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEligibilityGroupResource_basic(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_group.test"
	group1 := gofakeit.Email()
	group2 := gofakeit.Email()
	groupId1 := gofakeit.UUID()
	groupId2 := gofakeit.UUID()
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
				Config: testAccEligibilityGroupResourceConfig(group1, groupId1, approval1, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_name", group1),
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
				Config: testAccEligibilityGroupResourceConfig(group2, groupId2, approval2, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "group_name", group2),
					resource.TestCheckResourceAttr(resourceName, "approval_required", fmt.Sprint(approval2)),
				),
			},
		},
	})
}

func TestAccEligibilityGroupResource_missingAccountsAndOUs(t *testing.T) {
	group1 := gofakeit.Email()
	groupId1 := gofakeit.UUID()
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
				Config: testAccEligibilityGroupMissingAccountsAndOUsResourceConfig(group1, groupId1, approval1, duration, ticketNo, permissionArn, permissionName),
				ExpectError: regexp.MustCompile(
					"At Least One Account or OU Must Be Specified.",
				),
			},
			{
				Config: testAccEligibilityGroupEmptyAccountsAndOUsResourceConfig(group1, groupId1, approval1, duration, ticketNo, permissionArn, permissionName),
				ExpectError: regexp.MustCompile(
					"At Least One Account or OU Must Be Specified.",
				),
			},
		},
	})
}

func TestAccEligibilityGroupResource_missingAccounts(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_group.test"
	group1 := gofakeit.Email()
	groupId1 := gofakeit.UUID()
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
				Config: testAccEligibilityGroupMissingAccountsResourceConfig(group1, groupId1, approval1, duration, ticketNo, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_name", group1),
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

func TestAccEligibilityGroupResource_missingOUs(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_group.test"
	group1 := gofakeit.Email()
	groupId1 := gofakeit.UUID()
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
				Config: testAccEligibilityGroupMissingOUsResourceConfig(group1, groupId1, approval1, duration, ticketNo, accountId, accountName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_name", group1),
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

func TestAccEligibilityGroupResource_emptyAccounts(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_group.test"
	group1 := gofakeit.Email()
	groupId1 := gofakeit.UUID()
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
				Config: testAccEligibilityGroupEmptyAccountsResourceConfig(group1, groupId1, approval1, duration, ticketNo, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_name", group1),
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

func TestAccEligibilityGroupResource_emptyOUs(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_group.test"
	group1 := gofakeit.Email()
	groupId1 := gofakeit.UUID()
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
				Config: testAccEligibilityGroupEmptyOUsResourceConfig(group1, groupId1, approval1, duration, ticketNo, accountId, accountName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_name", group1),
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

func TestAccEligibilityGroupResource_Accounts(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_group.test"
	group1 := gofakeit.Email()
	groupId1 := gofakeit.UUID()
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
				Config: testAccEligibilityGroupResourceConfigNoAccounts(group1, groupId1, approval1, duration, ticketNo, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_id", groupId1),
					resource.TestCheckResourceAttr(resourceName, "group_name", group1),
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

func TestAccEligibilityGroupResource_disappears(t *testing.T) {
	ctx := context.Background()
	resourceName := "awsteam_eligibility_group.test"
	group1 := gofakeit.Email()
	groupId := gofakeit.UUID()
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
				Config: testAccEligibilityGroupResourceConfig(group1, groupId, approval1, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName),
				Check: resource.ComposeTestCheckFunc(
					testAccEligibilityResourceExists(ctx, resourceName),
					testAccEligibilityResourceDisappears(ctx, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccEligibilityGroupResourceConfig(group string, groupId string, approvalRequired bool, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
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
}`, group, groupId, approvalRequired, duration, ticketNo, accountId, accountName, ouId, ouName, permissionArn, permissionName)
}

func testAccEligibilityGroupResourceConfigNoAccounts(group string, groupId string, approvalRequired bool, duration, ticketNo, ouId, ouName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
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
}`, group, groupId, approvalRequired, duration, ticketNo, ouId, ouName, permissionArn, permissionName)
}

func testAccEligibilityGroupMissingAccountsAndOUsResourceConfig(group string, groupId string, approvalRequired bool, duration, ticketNo, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
	approval_required = %t
	duration          = %s
	ticket_no = "%s"
	permissions = [
		{
		permission_arn   = "%s"
		permission_name = "%s"
		}
	]
}`, group, groupId, approvalRequired, duration, ticketNo, permissionArn, permissionName)
}

func testAccEligibilityGroupEmptyAccountsAndOUsResourceConfig(group string, groupId string, approvalRequired bool, duration, ticketNo, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
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
}`, group, groupId, approvalRequired, duration, ticketNo, permissionArn, permissionName)
}

func testAccEligibilityGroupMissingAccountsResourceConfig(group string, groupId string, approvalRequired bool, duration, ticketNo, ouId, ouName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
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
}`, group, groupId, approvalRequired, duration, ticketNo, ouId, ouName, permissionArn, permissionName)
}

func testAccEligibilityGroupMissingOUsResourceConfig(group string, groupId string, approvalRequired bool, duration, ticketNo, accountId, accountName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
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
}`, group, groupId, approvalRequired, duration, ticketNo, accountId, accountName, permissionArn, permissionName)
}

func testAccEligibilityGroupEmptyAccountsResourceConfig(group string, groupId string, approvalRequired bool, duration, ticketNo, ouId, ouName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
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
}`, group, groupId, approvalRequired, duration, ticketNo, ouId, ouName, permissionArn, permissionName)
}

func testAccEligibilityGroupEmptyOUsResourceConfig(group string, groupId string, approvalRequired bool, duration, ticketNo, accountId, accountName, permissionArn, permissionName string) string {
	return fmt.Sprintf(`
resource "awsteam_eligibility_group" "test" {
	group_name         = "%s"
	group_id = "%s"
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
}`, group, groupId, approvalRequired, duration, ticketNo, accountId, accountName, permissionArn, permissionName)
}
