package provider

import (
	"context"
	"fmt"

	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/names"
	"github.com/awsteam-contrib/terraform-provider-awsteam/internal/sdk/awsteam"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &SettingsDataSource{}

func NewSettingsDataSource() datasource.DataSource {
	return &SettingsDataSource{}
}

type SettingsDataSource struct {
	client *awsteam.Client
}

func (d *SettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (d *SettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a data source for AWS TEAM Settings",

		Attributes: map[string]schema.Attribute{
			"approval": schema.BoolAttribute{
				MarkdownDescription: "If disabled, approval will not be required for all elevated access requests. If enabled, approval requirement is managed in eligibility policy configuration.",
				Computed:            true,
			},
			"comments": schema.BoolAttribute{
				MarkdownDescription: "Determines if comment field is mandatory for all elevated access requests.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Settings identifier.",
				Computed:            true,
			},
			"duration": schema.Int64Attribute{
				MarkdownDescription: "Default maximum request duration in hours.",
				Computed:            true,
			},
			"expiry": schema.Int64Attribute{
				MarkdownDescription: "Number of time in hours before an unapproved TEAM request expires.",
				Computed:            true,
			},
			"ses_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable sending notifications via Amazon SES.",
				Computed:            true,
			},
			"ses_source_arn": schema.StringAttribute{
				MarkdownDescription: "ARN of a verified SES identity in another AWS account. Must be configured to authorize sending mail from the TEAM account.",
				Computed:            true,
			},
			"ses_source_email": schema.StringAttribute{
				MarkdownDescription: "Email address to send notifications from. Must be verified in SES.",
				Computed:            true,
			},
			"sns_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Send notifications via Amazon SNS. Once enabled, create a subscription to the SNS topic (TeamNotifications-main) in the TEAM account.",
				Computed:            true,
			},
			"slack_notifications_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable to send notifications directly to users in Slack via a Slack bot app.",
				Computed:            true,
			},
			"slack_token": schema.StringAttribute{
				MarkdownDescription: "Slack OAuth token associated with the installed app.",
				Computed:            true,
			},
			"team_admin_group": schema.StringAttribute{
				MarkdownDescription: "Group of users responsible for managing TEAM administrative configurations",
				Computed:            true,
			},
			"team_auditor_group": schema.StringAttribute{
				MarkdownDescription: "Group of users allowed to audit TEAM elevated access requests",
				Computed:            true,
			},
			"ticket_no": schema.BoolAttribute{
				MarkdownDescription: "Determines if ticket number field is mandatory for elevated access requests",
				Computed:            true,
			},
			names.AttrModifiedBy: ModifiedByAttribute(),
			names.AttrCreatedAt:  CreatedAtAttribute(),
			names.AttrUpdatedAt:  UpdatedAtAttribute(),
		},
	}
}

func (d *SettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SettingsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	in := &awsteam.GetSettingsInput{}

	out, err := d.client.GetSettings(ctx, in)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}

	if out.Settings == nil {
		resp.Diagnostics.AddError("Client Error", "Settings does not exist")
		return
	}

	data.flatten(out.Settings)
	tflog.Trace(ctx, "read settings resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
