package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/YakDriver/regexache"
	"github.com/aws/smithy-go/ptr"
	"github.com/brittandeyoung/terraform-provider-awsteam/internal/names"
	"github.com/brittandeyoung/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/google/uuid"
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
	EligibilityUserType = "User"
)

var _ resource.Resource = &EligibilityUserResource{}
var _ resource.ResourceWithImportState = &EligibilityUserResource{}

func NewEligibilityUserResource() resource.Resource {
	return &EligibilityUserResource{}
}

type EligibilityUserResource struct {
	client *awsteam.Client
}

type EligibilityUserModel struct {
	Id               types.String `tfsdk:"id"`
	UserName         types.String `tfsdk:"user_name"`
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

func (r *EligibilityUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eligibility_user"
}

func (r *EligibilityUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Allows configuration of eligibility policies for an aws iam identity center user account within an AWS TEAM deployment.",

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
			"user_name": schema.StringAttribute{
				MarkdownDescription: "Name of the AWS iam identity center user the eligibility policy will be applied to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexache.MustCompile(`[\s\S]*`),
						"value must be a valid aws user name.",
					),
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

func (r *EligibilityUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EligibilityUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EligibilityUserModel

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
		Id:               ptr.String(fmt.Sprint(uuid.New())),
		Type:             ptr.String(EligibilityUserType),
		Name:             data.UserName.ValueStringPointer(),
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
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create eligibility user, got error: %s", err))
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

	diags := data.flatten(out.Eligibility)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EligibilityUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EligibilityUserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.GetEligibilityInput{
		Id: data.Id.ValueStringPointer(),
	}

	out, err := r.client.GetEligibility(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read eligibility user policy, got error: %s", err))
		return
	}

	if out == nil {
		resp.Diagnostics.AddError("Read Error", "Received empty Eligibility.")
		return
	}

	if out.Eligibility == nil {
		resp.Diagnostics.AddError("Read Error", "Received empty Eligibility.")
		return
	}

	diags := data.flatten(out.Eligibility)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read eligibility user resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EligibilityUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config, plan, state EligibilityUserModel

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
			Type:             ptr.String(EligibilityUserType),
			Name:             plan.UserName.ValueStringPointer(),
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
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update eligibility user, got error: %s", err))
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

		diags := plan.flatten(out.Eligibility)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Trace(ctx, "updated eligibility user resource")

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EligibilityUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EligibilityUserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.DeleteEligibilityInput{
		Id: data.Id.ValueStringPointer(),
	}

	_, err := r.client.DeleteEligibility(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete eligibility user, got error: %s", err))
		return
	}
}

func (r *EligibilityUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (d *EligibilityUserModel) flatten(out *awsteam.Eligibility) diag.Diagnostics {
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
	d.UserName = types.StringPointerValue(out.Name)
	d.Accounts = accountsSet
	d.OUs = ousSet
	d.Permissions = permissionsSet
	d.ApprovalRequired = types.BoolPointerValue(out.ApprovalRequired)
	d.Duration = types.Int64PointerValue(out.Duration)
	d.TicketNo = types.StringPointerValue(out.TicketNo)
	d.ModifiedBy = types.StringPointerValue(out.ModifiedBy)
	d.UpdatedAt = types.StringPointerValue(out.UpdatedAt)
	d.CreatedAt = types.StringPointerValue(out.UpdatedAt)

	return diags
}
