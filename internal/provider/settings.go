package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/names"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &SettingsResource{}
var _ resource.ResourceWithImportState = &SettingsResource{}

func NewSettingsResource() resource.Resource {
	return &SettingsResource{}
}

type SettingsResource struct {
	client *awsteam.Client
}

type SettingsModel struct {
	Approval                  types.Bool   `tfsdk:"approval"`
	Comments                  types.Bool   `tfsdk:"comments"`
	Id                        types.String `tfsdk:"id"`
	Duration                  types.Int64  `tfsdk:"duration"`
	Expiry                    types.Int64  `tfsdk:"expiry"`
	SesNotificationsEnabled   types.Bool   `tfsdk:"ses_notifications_enabled"`
	SnsNotificationsEnabled   types.Bool   `tfsdk:"sns_notifications_enabled"`
	SlackNotificationsEnabled types.Bool   `tfsdk:"slack_notifications_enabled"`
	SesSourceEmail            types.String `tfsdk:"ses_source_email"`
	SesSourceArn              types.String `tfsdk:"ses_source_arn"`
	SlackToken                types.String `tfsdk:"slack_token"`
	TeamAdminGroup            types.String `tfsdk:"team_admin_group"`
	TeamAuditorGroup          types.String `tfsdk:"team_auditor_group"`
	TicketNo                  types.Bool   `tfsdk:"ticket_no"`
	ModifiedBy                types.String `tfsdk:"modified_by"`
	CreatedAt                 types.String `tfsdk:"created_at"`
	UpdatedAt                 types.String `tfsdk:"updated_at"`
}

func (r *SettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (r *SettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Allows configuration of the settings within an AWS TEAM deployment",
		MarkdownDescription: "Allows configuration of the settings within an AWS TEAM deployment.\n\n" +
			"> **Important:** By default the `settings` resource will already exist on a fresh deployment of AWS TEAM and will cause a create to fail. Use the import block example below when deploying to a fresh instance of AWS TEAM to import the existing `settings`\n",

		Attributes: map[string]schema.Attribute{
			"approval": schema.BoolAttribute{
				MarkdownDescription: "If disabled, approval will not be required for all elevated access requests. If enabled, approval requirement is managed in eligibility policy configuration.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"comments": schema.BoolAttribute{
				MarkdownDescription: "Determines if comment field is mandatory for all elevated access requests.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The settings identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"duration": schema.Int64Attribute{
				MarkdownDescription: "Default maximum request duration in hours.",
				Required:            true,
			},
			"expiry": schema.Int64Attribute{
				MarkdownDescription: "Number of time in hours before an unapproved TEAM request expires.",
				Required:            true,
			},
			"ses_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable sending notifications via Amazon SES.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"ses_source_arn": schema.StringAttribute{
				MarkdownDescription: "ARN of a verified SES identity in another AWS account. Must be configured to authorize sending mail from the TEAM account.",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Computed:            true,
			},
			"ses_source_email": schema.StringAttribute{
				MarkdownDescription: "Email address to send notifications from. Must be verified in SES.",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Computed:            true,
			},
			"sns_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Send notifications via Amazon SNS. Once enabled, create a subscription to the SNS topic (TeamNotifications-main) in the TEAM account.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"slack_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable to send notifications directly to users in Slack via a Slack bot app.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"slack_token": schema.StringAttribute{
				MarkdownDescription: "Slack OAuth token associated with the installed app.",
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Computed:            true,
				Sensitive:           true,
			},
			"team_admin_group": schema.StringAttribute{
				MarkdownDescription: "Group of users responsible for managing TEAM administrative configurations",
				Required:            true,
			},
			"team_auditor_group": schema.StringAttribute{
				MarkdownDescription: "Group of users allowed to audit TEAM elevated access requests",
				Required:            true,
			},
			"ticket_no": schema.BoolAttribute{
				MarkdownDescription: "Determines if ticket number field is mandatory for elevated access requests",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			names.AttrModifiedBy: ModifiedByAttribute(),
			names.AttrCreatedAt:  CreatedAtAttribute(),
			names.AttrUpdatedAt:  UpdatedAtAttribute(),
		},
	}
}

func (r *SettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SettingsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.CreateSettingsInput{
		TeamAdminGroup:            data.TeamAdminGroup.ValueStringPointer(),
		TeamAuditorGroup:          data.TeamAuditorGroup.ValueStringPointer(),
		Approval:                  data.Approval.ValueBoolPointer(),
		Comments:                  data.Comments.ValueBoolPointer(),
		SesNotificationsEnabled:   data.SesNotificationsEnabled.ValueBoolPointer(),
		SnsNotificationsEnabled:   data.SnsNotificationsEnabled.ValueBoolPointer(),
		SlackNotificationsEnabled: data.SlackNotificationsEnabled.ValueBoolPointer(),
		TicketNo:                  data.TicketNo.ValueBoolPointer(),
		Duration:                  data.Duration.ValueInt64Pointer(),
		Expiry:                    data.Expiry.ValueInt64Pointer(),
		ModifiedBy:                data.ModifiedBy.ValueStringPointer(),
	}

	if !data.SesSourceArn.IsNull() {
		in.SesSourceArn = data.SesSourceArn.ValueStringPointer()
	}

	if !data.SesSourceEmail.IsNull() {
		in.SesSourceEmail = data.SesSourceEmail.ValueStringPointer()
	}

	if !data.SlackToken.IsNull() {
		in.SlackToken = data.SlackToken.ValueStringPointer()
	}

	out, err := r.client.CreateSettings(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create settings, got error: %s", err))
		return
	}

	if out == nil {
		resp.Diagnostics.AddError("Create Error", "Received empty Settings.")
		return
	}

	if out.Settings == nil {
		resp.Diagnostics.AddError("Create Error", "Received empty Settings.")
		return
	}

	data.flatten(out.Settings)
	tflog.Trace(ctx, "created settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SettingsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.GetSettingsInput{}

	out, err := r.client.GetSettings(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}

	if out == nil {
		resp.Diagnostics.AddWarning("Read Error", "Received empty Settings. Removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	if out.Settings == nil {
		resp.Diagnostics.AddWarning("Read Error", "Received empty Settings. Removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	data.flatten(out.Settings)
	tflog.Trace(ctx, "read settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var config, plan, state SettingsModel

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
		in := &awsteam.UpdateSettingsInput{
			TeamAdminGroup:            plan.TeamAdminGroup.ValueStringPointer(),
			TeamAuditorGroup:          plan.TeamAuditorGroup.ValueStringPointer(),
			Approval:                  plan.Approval.ValueBoolPointer(),
			Comments:                  plan.Comments.ValueBoolPointer(),
			SesNotificationsEnabled:   plan.SesNotificationsEnabled.ValueBoolPointer(),
			SnsNotificationsEnabled:   plan.SnsNotificationsEnabled.ValueBoolPointer(),
			SlackNotificationsEnabled: plan.SlackNotificationsEnabled.ValueBoolPointer(),
			TicketNo:                  plan.TicketNo.ValueBoolPointer(),
			Duration:                  plan.Duration.ValueInt64Pointer(),
			Expiry:                    plan.Expiry.ValueInt64Pointer(),
			ModifiedBy:                plan.ModifiedBy.ValueStringPointer(),
			SesSourceArn:              plan.SesSourceArn.ValueStringPointer(),
			SesSourceEmail:            plan.SesSourceEmail.ValueStringPointer(),
			SlackToken:                plan.SlackToken.ValueStringPointer(),
		}

		out, err := r.client.UpdateSettings(ctx, in)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
			return
		}

		if out == nil {
			resp.Diagnostics.AddError("Update Error", "Received empty Settings.")
			return
		}

		if out.Settings == nil {
			resp.Diagnostics.AddError("Update Error", "Received empty Settings.")
			return
		}

		plan.flatten(out.Settings)

		tflog.Trace(ctx, "updated settings resource")

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SettingsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.DeleteSettingsInput{}

	_, err := r.client.DeleteSettings(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}
}

func (r *SettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (d *SettingsModel) flatten(out *awsteam.Settings) {
	d.Id = types.StringPointerValue(out.Id)
	d.Approval = types.BoolPointerValue(out.Approval)
	d.Comments = types.BoolPointerValue(out.Comments)
	d.CreatedAt = types.StringPointerValue(out.CreatedAt)
	d.Duration = types.Int64PointerValue(out.Duration)
	d.ModifiedBy = types.StringPointerValue(out.ModifiedBy)
	d.Expiry = types.Int64PointerValue(out.Expiry)
	d.SesNotificationsEnabled = types.BoolPointerValue(out.SesNotificationsEnabled)
	d.SesSourceArn = types.StringPointerValue(out.SesSourceArn)
	d.SesSourceEmail = types.StringPointerValue(out.SesSourceEmail)
	d.SlackNotificationsEnabled = types.BoolPointerValue(out.SlackNotificationsEnabled)
	d.SlackToken = types.StringPointerValue(out.SlackToken)
	d.SnsNotificationsEnabled = types.BoolPointerValue(out.SnsNotificationsEnabled)
	d.TeamAdminGroup = types.StringPointerValue(out.TeamAdminGroup)
	d.TeamAuditorGroup = types.StringPointerValue(out.TeamAuditorGroup)
	d.TicketNo = types.BoolPointerValue(out.TicketNo)
	d.UpdatedAt = types.StringPointerValue(out.UpdatedAt)
}
