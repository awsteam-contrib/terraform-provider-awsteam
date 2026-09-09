package awsteam

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/smithy-go/ptr"
)

type UpdateSettingsInput struct {
	Approval                  *bool
	Comments                  *bool
	Duration                  *int64
	Expiry                    *int64
	Id                        *string
	SesNotificationsEnabled   *bool
	SnsNotificationsEnabled   *bool
	SlackNotificationsEnabled *bool
	SesSourceEmail            *string
	SesSourceArn              *string
	SlackToken                *string
	TeamAdminGroup            *string
	TeamAuditorGroup          *string
	TicketNo                  *bool
	UseOUCache                *bool
	ModifiedBy                *string
	CreatedAt                 *string
	UpdatedAt                 *string
}

type UpdateSettingsOutput struct {
	Settings *Settings `json:"updateSettings"`
}

const updateSettingsMutation = `mutation UpdateSettings($input: UpdateSettingsInput!) {
	updateSettings(input: $input) {
		id
		duration
		expiry
		comments
		ticketNo
		approval
		modifiedBy
		sesNotificationsEnabled
		snsNotificationsEnabled
		slackNotificationsEnabled
		sesSourceEmail
		sesSourceArn
		slackToken
		teamAdminGroup
		teamAuditorGroup
		createdAt
		updatedAt
	}
}`

const updateSettingsMutationWithOUCache = `mutation UpdateSettings($input: UpdateSettingsInput!) {
	updateSettings(input: $input) {
		id
		duration
		expiry
		comments
		ticketNo
		approval
		modifiedBy
		sesNotificationsEnabled
		snsNotificationsEnabled
		slackNotificationsEnabled
		sesSourceEmail
		sesSourceArn
		slackToken
		teamAdminGroup
		teamAuditorGroup
		useOUCache
		createdAt
		updatedAt
	}
}`

func (client *Client) UpdateSettings(ctx context.Context, in *UpdateSettingsInput) (*UpdateSettingsOutput, error) {
	out := &UpdateSettingsOutput{}

	id := "settings"
	if in.Id != nil {
		id = *in.Id
	}

	q := updateSettingsMutation
	inputVars := map[string]interface{}{
		"id":                        id,
		"duration":                  fmt.Sprintf("%d", ptr.ToInt64(in.Duration)),
		"expiry":                    fmt.Sprintf("%d", ptr.ToInt64(in.Expiry)),
		"comments":                  ptr.ToBool(in.Comments),
		"ticketNo":                  ptr.ToBool(in.TicketNo),
		"approval":                  ptr.ToBool(in.Approval),
		"modifiedBy":                ptr.ToString(in.ModifiedBy),
		"sesNotificationsEnabled":   ptr.ToBool(in.SesNotificationsEnabled),
		"snsNotificationsEnabled":   ptr.ToBool(in.SnsNotificationsEnabled),
		"slackNotificationsEnabled": ptr.ToBool(in.SlackNotificationsEnabled),
		"sesSourceEmail":            ptr.ToString(in.SesSourceEmail),
		"sesSourceArn":              ptr.ToString(in.SesSourceArn),
		"slackToken":                ptr.ToString(in.SlackToken),
		"teamAdminGroup":            ptr.ToString(in.TeamAdminGroup),
		"teamAuditorGroup":          ptr.ToString(in.TeamAuditorGroup),
	}

	if client.SettingsCapabilities.UseOUCacheSupported {
		q = updateSettingsMutationWithOUCache
		inputVars["useOUCache"] = ptr.ToBool(in.UseOUCache)
	}

	vars := map[string]interface{}{"input": inputVars}

	raw, err := client.GraphClient.ExecRaw(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
