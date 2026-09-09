package awsteam

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/smithy-go/ptr"
)

type CreateSettingsInput struct {
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
}

type CreateSettingsOutput struct {
	Settings *Settings `json:"createSettings"`
}

const createSettingsMutation = `mutation CreateSettings($input: CreateSettingsInput!) {
	createSettings(input: $input) {
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

const createSettingsMutationWithOUCache = `mutation CreateSettings($input: CreateSettingsInput!) {
	createSettings(input: $input) {
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

func (client *Client) CreateSettings(ctx context.Context, in *CreateSettingsInput) (*CreateSettingsOutput, error) {
	out := &CreateSettingsOutput{}

	id := "settings"
	if in.Id != nil {
		id = *in.Id
	}

	q := createSettingsMutation
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
		q = createSettingsMutationWithOUCache
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
