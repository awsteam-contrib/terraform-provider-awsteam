package awsteam

import (
	"context"
	"fmt"
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
	ModifiedBy                *string
}

type CreateSettingsOutput struct {
	Settings *Settings `json:"createSettings"`
}

func (client *Client) CreateSettings(ctx context.Context, in *CreateSettingsInput) (*CreateSettingsOutput, error) {
	out := &CreateSettingsOutput{}
	id := defaultString(in.Id, "settings")

	query := fmt.Sprintf(`mutation CreateSettings {
		createSettings(
			input: {
				id: %s
				duration: %s
				expiry: %s
				comments: %s
				ticketNo: %s
				approval: %s
				modifiedBy: %s
				sesNotificationsEnabled: %s
				snsNotificationsEnabled: %s
				slackNotificationsEnabled: %s
				sesSourceEmail: %s
				sesSourceArn: %s
				slackToken: %s
				teamAdminGroup: %s
				teamAuditorGroup: %s
			}
		) {
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
	`, graphqlString(id),
		graphqlQuotedInt64Ptr(in.Duration),
		graphqlQuotedInt64Ptr(in.Expiry),
		graphqlBoolPtr(in.Comments),
		graphqlBoolPtr(in.TicketNo),
		graphqlBoolPtr(in.Approval),
		graphqlStringPtr(in.ModifiedBy),
		graphqlBoolPtr(in.SesNotificationsEnabled),
		graphqlBoolPtr(in.SnsNotificationsEnabled),
		graphqlBoolPtr(in.SlackNotificationsEnabled),
		graphqlStringPtr(in.SesSourceEmail),
		graphqlStringPtr(in.SesSourceArn),
		graphqlStringPtr(in.SlackToken),
		graphqlStringPtr(in.TeamAdminGroup),
		graphqlStringPtr(in.TeamAuditorGroup),
	)

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
