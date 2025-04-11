package provider

import (
	"context"
	"fmt"

	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	accountsAttrTypes = map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
)
var _ datasource.DataSource = &AccountsDataSource{}

func NewAccountsDataSource() datasource.DataSource {
	return &AccountsDataSource{}
}

type AccountsDataSource struct {
	client *awsteam.Client
}

type AccountsModel struct {
	Id       types.String `tfsdk:"id"`
	Accounts types.Set    `tfsdk:"accounts"`
}

func (d *AccountsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

func (d *AccountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a data source for AWS TEAMs Accounts",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Accounts Identifier. This is a static value of `accounts` as it contains all accounts.",
				Computed:            true,
			},
			"accounts": schema.SetNestedAttribute{
				MarkdownDescription: "A set of AWS accounts.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The AWS account id",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the AWS account.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AccountsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*awsteam.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *awsteam.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *AccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AccountsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.GetAccountsInput{}

	out, err := d.client.GetAccounts(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}

	if out.Accounts == nil {
		resp.Diagnostics.AddError("Client Error", "Accounts does not exist")
		return
	}

	data.flatten(out)
	tflog.Trace(ctx, "read settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *AccountsModel) flatten(out *awsteam.GetAccountsOutput) diag.Diagnostics {
	var diags diag.Diagnostics

	accountsSet, diags := flattenAccounts(out.Accounts)
	diags.Append(diags...)
	if diags.HasError() {
		return diags
	}

	d.Id = types.StringValue("accounts")
	d.Accounts = accountsSet

	return diags
}

func flattenAccounts(apiObject []*awsteam.Account) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics
	elemType := types.ObjectType{AttrTypes: accountsAttrTypes}
	elems := []attr.Value{}

	if len(apiObject) == 0 {
		return types.SetValueMust(elemType, []attr.Value{}), diags
	}

	for _, account := range apiObject {
		obj := map[string]attr.Value{
			"id":   types.StringPointerValue(account.Id),
			"name": types.StringPointerValue(account.Name),
		}
		objVal, d := types.ObjectValue(accountsAttrTypes, obj)
		diags.Append(d...)

		elems = append(elems, objVal)
	}
	setVal, d := types.SetValue(elemType, elems)
	diags.Append(d...)

	return setVal, diags
}
