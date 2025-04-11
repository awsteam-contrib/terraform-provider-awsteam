package provider

import (
	"github.com/YakDriver/regexache"
	"github.com/brittandeyoung/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	accountsAndOUsRequiredMessageSummary  = "At Least One Account or OU Must Be Specified."
	accountsAndOUsRequiredMessageDetail   = "Either 'Accounts' or 'OUs' must have at least one element. Both cannot be empty."
	accountsAndOUsRequiredMessageMarkdown = "\n> **NOTE:** The `ous` and `accounts` fields are labeled as *Optional*; however, at least one `ou` or `account` must be provided."
)

var (
	eligibilityAccountAttrTypes = map[string]attr.Type{
		"account_id":   types.StringType,
		"account_name": types.StringType,
	}
	eligibilityOUAttrTypes = map[string]attr.Type{
		"ou_id":   types.StringType,
		"ou_name": types.StringType,
	}
	eligibilityPermissionAttrTypes = map[string]attr.Type{
		"permission_arn":  types.StringType,
		"permission_name": types.StringType,
	}
)

type EligibilityAccount struct {
	AccountId   types.String `tfsdk:"account_id"`
	AccountName types.String `tfsdk:"account_name"`
}

type EligibilityOU struct {
	OUId   types.String `tfsdk:"ou_id"`
	OUName types.String `tfsdk:"ou_name"`
}

type EligibilityPermission struct {
	PermissionId   types.String `tfsdk:"permission_arn"`
	PermissionName types.String `tfsdk:"permission_name"`
}

func AccountAttributeSet() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "A list of AWS accounts the eligibility will apply to. " + accountsAndOUsRequiredMessageDetail,
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"account_id": schema.StringAttribute{
					MarkdownDescription: "The AWS account id the eligibility policy will be applied to. This needs to match the account id of the name provided in account_name.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexache.MustCompile(`\d{12}`),
							"value must be a valid aws account id.",
						),
					},
				},
				"account_name": schema.StringAttribute{
					MarkdownDescription: "Name of the AWS account the eligibility policy will be applied to. This needs to match the name of the account number provided in account_id.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexache.MustCompile(`[\s\S]*`),
							"value must be a valid aws account name.",
						),
					},
				},
			},
		},
	}
}

func OUAttributeSet() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "A list of AWS OUs the eligibility will apply to. " + accountsAndOUsRequiredMessageDetail,
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"ou_id": schema.StringAttribute{
					MarkdownDescription: "Id of the OU the eligibility policy will be applied to. This needs to match the id of the name provided in ou_name.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexache.MustCompile(`^(r-[0-9a-z]{4,32})|(ou-[0-9a-z]{4,32}-[a-z0-9]{8,32})$`),
							"value must be a valid aws ou id.",
						),
					},
				},
				"ou_name": schema.StringAttribute{
					MarkdownDescription: "Name of the OU the eligibility policy will be applied to. This needs to match the name of the id provided in ou_id.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexache.MustCompile(`[\s\S]*`),
							"value must be a valid aws ou name.",
						),
					},
				},
			},
		},
	}
}

func PermissionAttributeSet() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "A list of AWS permission sets for the eligibility policy.",
		Required:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"permission_arn": schema.StringAttribute{
					MarkdownDescription: "The ARN of the permission for the eligibility policy. This needs to match the ARN of the name provided in name.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexache.MustCompile(`^arn:(aws|aws-us-gov|aws-cn|aws-iso|aws-iso-b):sso:::permissionSet/(sso)?ins-[a-zA-Z0-9-.]{16}/ps-[a-zA-Z0-9-./]{16}$`),
							"value must be a valid AWS permissionSet ARN.",
						),
					},
				},
				"permission_name": schema.StringAttribute{
					MarkdownDescription: "Name of the permission for the eligibility policy. This needs to match the name of the ARN provided in ARN.",
					Required:            true,
				},
			},
		},
	}
}

func expandEligibilityAccounts(raw []*EligibilityAccount) []*awsteam.EligibilityAccount {
	var accounts []*awsteam.EligibilityAccount

	if len(raw) == 0 {
		return accounts
	}

	for _, v := range raw {
		account := &awsteam.EligibilityAccount{}
		account.Id = v.AccountId.ValueStringPointer()
		account.Name = v.AccountName.ValueStringPointer()
		accounts = append(accounts, account)
	}

	return accounts
}

func expandEligibilityOUs(raw []*EligibilityOU) []*awsteam.EligibilityOU {
	var ous []*awsteam.EligibilityOU

	if len(raw) == 0 {
		return ous
	}

	for _, v := range raw {
		ou := &awsteam.EligibilityOU{}
		ou.Id = v.OUId.ValueStringPointer()
		ou.Name = v.OUName.ValueStringPointer()
		ous = append(ous, ou)
	}

	return ous
}

func expandEligibilityPermissions(raw []*EligibilityPermission) []*awsteam.EligibilityPermission {
	var permissions []*awsteam.EligibilityPermission

	if len(raw) == 0 {
		return permissions
	}

	for _, v := range raw {
		permission := &awsteam.EligibilityPermission{}
		permission.Id = v.PermissionId.ValueStringPointer()
		permission.Name = v.PermissionName.ValueStringPointer()
		permissions = append(permissions, permission)
	}

	return permissions
}

func flattenEligibilityAccounts(apiObject []*awsteam.EligibilityAccount) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: eligibilityAccountAttrTypes}
	elems := []attr.Value{}

	if len(apiObject) == 0 {
		return types.SetValueMust(elemType, []attr.Value{}), diags
	}

	for _, account := range apiObject {
		obj := map[string]attr.Value{
			"account_id":   types.StringPointerValue(account.Id),
			"account_name": types.StringPointerValue(account.Name),
		}
		objVal, d := types.ObjectValue(eligibilityAccountAttrTypes, obj)
		diags.Append(d...)

		elems = append(elems, objVal)
	}
	setVal, d := types.SetValue(elemType, elems)
	diags.Append(d...)

	return setVal, diags
}

func flattenEligibilityOUs(apiObject []*awsteam.EligibilityOU) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: eligibilityOUAttrTypes}
	elems := []attr.Value{}

	if len(apiObject) == 0 {
		return types.SetValueMust(elemType, []attr.Value{}), diags
	}

	for _, ou := range apiObject {
		obj := map[string]attr.Value{
			"ou_id":   types.StringPointerValue(ou.Id),
			"ou_name": types.StringPointerValue(ou.Name),
		}
		objVal, d := types.ObjectValue(eligibilityOUAttrTypes, obj)
		diags.Append(d...)

		elems = append(elems, objVal)
	}
	setVal, d := types.SetValue(elemType, elems)
	diags.Append(d...)

	return setVal, diags
}

func flattenEligibilityPermissions(apiObject []*awsteam.EligibilityPermission) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: eligibilityPermissionAttrTypes}

	if len(apiObject) == 0 {
		return types.SetValueMust(elemType, []attr.Value{}), diags
	}

	elems := []attr.Value{}
	for _, permission := range apiObject {
		obj := map[string]attr.Value{
			"permission_arn":  types.StringPointerValue(permission.Id),
			"permission_name": types.StringPointerValue(permission.Name),
		}
		objVal, d := types.ObjectValue(eligibilityPermissionAttrTypes, obj)
		diags.Append(d...)

		elems = append(elems, objVal)
	}
	setVal, d := types.SetValue(elemType, elems)
	diags.Append(d...)

	return setVal, diags
}
