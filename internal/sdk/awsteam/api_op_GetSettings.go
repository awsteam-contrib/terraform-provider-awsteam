package awsteam

import (
	"context"
	"encoding/json"
)

type GetSettingsInput struct {
	Id *string
}

type GetSettingsOutput struct {
	Settings *Settings `json:"getSettings"`
}

const getSettingsQuery = `query GetSettings($id: ID!) {
	getSettings(id: $id) {
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

const getSettingsQueryWithOUCache = `query GetSettings($id: ID!) {
	getSettings(id: $id) {
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

func (client *Client) GetSettings(ctx context.Context, in *GetSettingsInput) (*GetSettingsOutput, error) {
	out := &GetSettingsOutput{}

	id := "settings"
	if in.Id != nil {
		id = *in.Id
	}

	q := getSettingsQuery
	if client.SettingsCapabilities.UseOUCacheSupported {
		q = getSettingsQueryWithOUCache
	}

	vars := map[string]interface{}{"id": id}

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
