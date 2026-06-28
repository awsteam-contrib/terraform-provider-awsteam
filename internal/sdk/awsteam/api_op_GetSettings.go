package awsteam

import (
	"context"
	"fmt"
)

type GetSettingsInput struct {
	Id *string
}

type GetSettingsOutput struct {
	Settings *Settings `json:"getSettings"`
}

func (client *Client) GetSettings(ctx context.Context, in *GetSettingsInput) (*GetSettingsOutput, error) {
	out := &GetSettingsOutput{}
	id := defaultString(in.Id, "settings")

	query := fmt.Sprintf(`query GetSettings {
		getSettings(id: %s) {
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
	}
	`, graphqlString(id))

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
