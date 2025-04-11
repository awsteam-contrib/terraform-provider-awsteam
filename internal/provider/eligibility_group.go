package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/YakDriver/regexache"
	"github.com/aws/smithy-go/ptr"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/names"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	EligibilityGroupType = "Group"
)

var _ resource.Resource = &EligibilityGroupResource{}
var _ resource.ResourceWithImportState = &EligibilityGroupResource{}
var _ resource.ResourceWithValidateConfig = &EligibilityGroupResource{}

func NewEligibilityGroupResource() resource.Resource {
	return &EligibilityGroupResource{}
}

type EligibilityGroupResource struct {
	client *awsteam.Client
}

type EligibilityGroupModel struct {
	Id               types.String `tfsdk:"id"`
	GroupName        types.String `tfsdk:"group_name"`
	GroupId          types.String `tfsdk:"group_id"`
	Accounts         types.Set    `tfsdk:"accounts"`
	OUs              types.Set    `tfsdk:"ous"`
	Permissions      types.Set    `tfsdk:"permissions"`
	TicketNo         types.String `tfsdk:"ticket_no"`
	ApprovalRequired types.Bool   `tfsdk:"approval_required"`
	Duration         types.Int64  `tfsdk:"duration"`
	ModifiedBy       types.String `tfsdk:"modified_by"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

func (r *EligibilityGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eligibility_group"
}

func (r *EligibilityGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Allows configuration of eligibility policies for an aws iam identity center group account within an AWS TEAM deployment." +
			accountsAndOUsRequiredMessageMarkdown,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The UUID of the eligibility.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"approval_required": schema.BoolAttribute{
				MarkdownDescription: "Determines if approval is required for elevated access",
				Required:            true,
			},
			"group_name": schema.StringAttribute{
				MarkdownDescription: "Name of the AWS iam identity center group the eligibility policy will be applied to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexache.MustCompile(`[\s\S]*`),
						"value must be a valid aws group name.",
					),
				},
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "Id of the AWS iam identity center group the eligibility policy will be applied to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"duration": schema.Int64Attribute{
				MarkdownDescription: "The maximum elevated access request duration in hours.",
				Required:            true,
			},
			"ticket_no": schema.StringAttribute{
				MarkdownDescription: "The Change Management system ticket system number.",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Computed:            true,
			},
			names.AttrAccountSet:    AccountAttributeSet(),
			names.AttrOUSet:         OUAttributeSet(),
			names.AttrPermissionSet: PermissionAttributeSet(),
			names.AttrModifiedBy:    ModifiedByAttribute(),
			names.AttrCreatedAt:     CreatedAtAttribute(),
			names.AttrUpdatedAt:     UpdatedAtAttribute(),
		},
	}
}

func (r *EligibilityGroupResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config EligibilityGroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(config.Accounts.Elements()) == 0 && len(config.OUs.Elements()) == 0 {
		resp.Diagnostics.AddError(
			accountsAndOUsRequiredMessageSummary,
			accountsAndOUsRequiredMessageDetail,
		)
	}
}

func (r *EligibilityGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*awsteam.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *awsteam.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *EligibilityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EligibilityGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var accounts []*EligibilityAccount
	var ous []*EligibilityOU
	var permissions []*EligibilityPermission

	resp.Diagnostics.Append(data.Accounts.ElementsAs(ctx, &accounts, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.OUs.ElementsAs(ctx, &ous, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(data.Permissions.ElementsAs(ctx, &permissions, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.CreateEligibilityInput{
		Id:               data.GroupId.ValueStringPointer(),
		Type:             ptr.String(EligibilityGroupType),
		Name:             data.GroupName.ValueStringPointer(),
		ApprovalRequired: data.ApprovalRequired.ValueBoolPointer(),
		Duration:         data.Duration.ValueInt64Pointer(),
		TicketNo:         data.TicketNo.ValueStringPointer(),
		ModifiedBy:       data.ModifiedBy.ValueStringPointer(),
		Accounts:         expandEligibilityAccounts(accounts),
		OUs:              expandEligibilityOUs(ous),
		Permissions:      expandEligibilityPermissions(permissions),
	}

	out, err := r.client.CreateEligibility(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create eligibility group, got error: %s", err))
		return
	}

	if out == nil {
		resp.Diagnostics.AddError("Create Error", "Received empty Eligibility.")
		return
	}

	if out.Eligibility == nil {
		resp.Diagnostics.AddError("Create Error", "Received empty Eligibility.")
		return
	}

	diags := data.flatten(data, out.Eligibility)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EligibilityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EligibilityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.GetEligibilityInput{
		Id: data.Id.ValueStringPointer(),
	}

	out, err := r.client.GetEligibility(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read eligibility group policy, got error: %s", err))
		return
	}

	if out == nil {
		resp.Diagnostics.AddWarning("Read Error", "Received empty Eligibility. Removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	if out.Eligibility == nil {
		resp.Diagnostics.AddWarning("Read Error", "Received empty Eligibility. Removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	diags := data.flatten(data, out.Eligibility)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read eligibility group resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EligibilityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config, plan, state EligibilityGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateRequired := false

	if !reflect.DeepEqual(state, plan) {
		updateRequired = true
	}

	if updateRequired {
		var accounts []*EligibilityAccount
		var ous []*EligibilityOU
		var permissions []*EligibilityPermission

		resp.Diagnostics.Append(plan.Accounts.ElementsAs(ctx, &accounts, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(plan.OUs.ElementsAs(ctx, &ous, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(plan.Permissions.ElementsAs(ctx, &permissions, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		in := &awsteam.UpdateEligibilityInput{
			Id:               state.Id.ValueStringPointer(),
			Type:             ptr.String(EligibilityGroupType),
			Name:             plan.GroupName.ValueStringPointer(),
			ApprovalRequired: plan.ApprovalRequired.ValueBoolPointer(),
			Duration:         plan.Duration.ValueInt64Pointer(),
			TicketNo:         plan.TicketNo.ValueStringPointer(),
			ModifiedBy:       plan.ModifiedBy.ValueStringPointer(),
			Accounts:         expandEligibilityAccounts(accounts),
			OUs:              expandEligibilityOUs(ous),
			Permissions:      expandEligibilityPermissions(permissions),
		}

		out, err := r.client.UpdateEligibility(ctx, in)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update eligibility group, got error: %s", err))
			return
		}

		if out == nil {
			resp.Diagnostics.AddError("Refresh Error", "Received empty Eligibility.")
			return
		}

		if out.Eligibility == nil {
			resp.Diagnostics.AddError("Refresh Error", "Received empty Eligibility.")
			return
		}

		diags := plan.flatten(config, out.Eligibility)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Trace(ctx, "updated eligibility group resource")

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EligibilityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EligibilityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.DeleteEligibilityInput{
		Id: data.Id.ValueStringPointer(),
	}

	_, err := r.client.DeleteEligibility(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete eligibility group, got error: %s", err))
		return
	}
}

func (r *EligibilityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (d *EligibilityGroupModel) flatten(config EligibilityGroupModel, out *awsteam.Eligibility) diag.Diagnostics {
	var diags diag.Diagnostics

	accountsSet, diags := flattenEligibilityAccounts(out.Accounts)
	diags.Append(diags...)
	if diags.HasError() {
		return diags
	}

	ousSet, diags := flattenEligibilityOUs(out.OUs)
	diags.Append(diags...)
	if diags.HasError() {
		return diags
	}

	permissionsSet, diags := flattenEligibilityPermissions(out.Permissions)
	diags.Append(diags...)
	if diags.HasError() {
		return diags
	}

	d.Id = types.StringPointerValue(out.Id)
	d.GroupName = types.StringPointerValue(out.Name)
	d.GroupId = types.StringPointerValue(out.Id)
	if !config.Accounts.IsNull() || len(accountsSet.Elements()) != 0 {
		d.Accounts = accountsSet
	}
	if !config.OUs.IsNull() || len(ousSet.Elements()) != 0 {
		d.OUs = ousSet
	}
	d.Permissions = permissionsSet
	d.ApprovalRequired = types.BoolPointerValue(out.ApprovalRequired)
	d.Duration = types.Int64PointerValue(out.Duration)
	d.TicketNo = types.StringPointerValue(out.TicketNo)
	d.ModifiedBy = types.StringPointerValue(out.ModifiedBy)
	d.UpdatedAt = types.StringPointerValue(out.UpdatedAt)
	d.CreatedAt = types.StringPointerValue(out.UpdatedAt)

	return diags
}
