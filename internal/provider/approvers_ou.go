package provider

import (
	"context"
	"fmt"

	"github.com/YakDriver/regexache"
	"github.com/aws/smithy-go/ptr"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/names"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	ApproversOUType = "OU"
)

var _ resource.Resource = &ApproversOUResource{}
var _ resource.ResourceWithImportState = &ApproversOUResource{}

func NewApproversOUResource() resource.Resource {
	return &ApproversOUResource{}
}

type ApproversOUResource struct {
	client *awsteam.Client
}

type ApproversOUModel struct {
	Id         types.String `tfsdk:"id"`
	OUName     types.String `tfsdk:"ou_name"`
	Approvers  types.Set    `tfsdk:"approvers"`
	GroupIds   types.Set    `tfsdk:"group_ids"`
	OUId       types.String `tfsdk:"ou_id"`
	TicketNo   types.String `tfsdk:"ticket_no"`
	ModifiedBy types.String `tfsdk:"modified_by"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *ApproversOUResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_approvers_ou"
}

func (r *ApproversOUResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Allows configuration of approval policies for an OU within an AWS TEAM deployment.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The approvers ou identifier. This is the same as the ou_id.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ou_name": schema.StringAttribute{
				MarkdownDescription: "Name of the OU the approvers policy will be applied to. This needs to match the name of the id provided in ou_id.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexache.MustCompile(`[\s\S]*`),
						"value must be a valid aws ou name.",
					),
				},
			},
			"ou_id": schema.StringAttribute{
				MarkdownDescription: "Id of the OU the approvers policy will be applied to. This needs to match the id of the name provided in ou_name.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexache.MustCompile(`^(r-[0-9a-z]{4,32})|(ou-[0-9a-z]{4,32}-[a-z0-9]{8,32})$`),
						"value must be a valid aws ou id.",
					),
				},
			},
			"approvers": schema.SetAttribute{
				MarkdownDescription: "The list of group names that will be approvers for the OU. This needs to match the names of the ids provided in group_ids.",
				ElementType:         types.StringType,
				Required:            true,
			},
			"group_ids": schema.SetAttribute{
				MarkdownDescription: "The list of group names that will be approvers for the OU. This needs to match the ids of the names provided in approvers.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"ticket_no": schema.StringAttribute{
				MarkdownDescription: "The Change Management system ticket system number.",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Computed:            true,
			},
			names.AttrModifiedBy: ModifiedByAttribute(),
			names.AttrCreatedAt:  CreatedAtAttribute(),
			names.AttrUpdatedAt:  UpdatedAtAttribute(),
		},
	}
}

func (r *ApproversOUResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ApproversOUResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApproversOUModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var approvers []*string
	resp.Diagnostics.Append(data.Approvers.ElementsAs(ctx, &approvers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groupIds []*string
	resp.Diagnostics.Append(data.GroupIds.ElementsAs(ctx, &groupIds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.CreateApproversInput{
		Id:         data.OUId.ValueStringPointer(),
		Type:       ptr.String(ApproversOUType),
		Name:       data.OUName.ValueStringPointer(),
		Approvers:  approvers,
		GroupIds:   groupIds,
		TicketNo:   data.TicketNo.ValueStringPointer(),
		ModifiedBy: data.ModifiedBy.ValueStringPointer(),
	}

	out, err := r.client.CreateApprovers(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create approvers ou, got error: %s", err))
		return
	}

	if out == nil {
		resp.Diagnostics.AddError("Create Error", "Received empty Approvers.")
		return
	}

	if out.Approvers == nil {
		resp.Diagnostics.AddError("Create Error", "Received empty Approvers.")
		return
	}

	diags := data.flatten(ctx, out.Approvers)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApproversOUResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApproversOUModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.GetApproversInput{
		Id: data.Id.ValueStringPointer(),
	}

	out, err := r.client.GetApprovers(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read approvers ou policy, got error: %s", err))
		return
	}

	if out == nil {
		resp.Diagnostics.AddWarning("Read Error", "Received empty Approvers. Removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	if out.Approvers == nil {
		resp.Diagnostics.AddWarning("Read Error", "Received empty Approvers. Removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	diags := data.flatten(ctx, out.Approvers)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "read approvers ou resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApproversOUResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config, plan, state ApproversOUModel

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

	var approvers []*string
	resp.Diagnostics.Append(plan.Approvers.ElementsAs(ctx, &approvers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groupIds []*string
	resp.Diagnostics.Append(plan.GroupIds.ElementsAs(ctx, &groupIds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.UpdateApproversInput{}

	if !plan.Approvers.IsUnknown() && !state.Approvers.Equal(plan.Approvers) {
		updateRequired = true
	}

	if !plan.GroupIds.IsUnknown() && !state.GroupIds.Equal(plan.GroupIds) {
		updateRequired = true
	}

	if !plan.TicketNo.IsUnknown() && !state.TicketNo.Equal(plan.TicketNo) {
		updateRequired = true
	}

	if !plan.ModifiedBy.IsUnknown() && !state.ModifiedBy.Equal(plan.ModifiedBy) {
		updateRequired = true
	}

	if updateRequired {
		in.Id = plan.OUId.ValueStringPointer()
		in.Type = ptr.String(ApproversOUType)
		in.Name = plan.OUName.ValueStringPointer()
		in.Approvers = approvers
		in.GroupIds = groupIds
		in.TicketNo = plan.TicketNo.ValueStringPointer()
		in.ModifiedBy = plan.ModifiedBy.ValueStringPointer()

		out, err := r.client.UpdateApprovers(ctx, in)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update approvers ou, got error: %s", err))
			return
		}

		if out == nil {
			resp.Diagnostics.AddError("Refresh Error", "Received empty Approvers.")
			return
		}

		if out.Approvers == nil {
			resp.Diagnostics.AddError("Refresh Error", "Received empty Approvers.")
			return
		}

		diags := plan.flatten(ctx, out.Approvers)

		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Trace(ctx, "updated approvers ou resource")

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApproversOUResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApproversOUModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.DeleteApproversInput{
		Id: data.Id.ValueStringPointer(),
	}

	_, err := r.client.DeleteApprovers(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete approvers ou, got error: %s", err))
		return
	}
}

func (r *ApproversOUResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (d *ApproversOUModel) flatten(ctx context.Context, out *awsteam.Approvers) diag.Diagnostics {
	var diags diag.Diagnostics

	approversSet, diag := types.SetValueFrom(ctx, types.StringType, out.Approvers)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	groupIdsSet, diag := types.SetValueFrom(ctx, types.StringType, out.GroupIds)
	diags.Append(diag...)
	if diags.HasError() {
		return diags
	}

	d.Id = types.StringPointerValue(out.Id)
	d.OUName = types.StringPointerValue(out.Name)
	d.OUId = types.StringPointerValue(out.Id)
	d.Approvers = approversSet
	d.GroupIds = groupIdsSet
	d.TicketNo = types.StringPointerValue(out.TicketNo)
	d.ModifiedBy = types.StringPointerValue(out.ModifiedBy)
	d.UpdatedAt = types.StringPointerValue(out.UpdatedAt)
	d.CreatedAt = types.StringPointerValue(out.UpdatedAt)

	return diags
}
