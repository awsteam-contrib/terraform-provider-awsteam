package awsteam

import (
	"context"
	"fmt"
)

type DeleteSettingsInput struct {
	Id *string
}

type DeleteSettingsOutput struct {
	Settings *Settings `json:"deleteSettings"`
}

func (client *Client) DeleteSettings(ctx context.Context, in *DeleteSettingsInput) (*DeleteSettingsOutput, error) {
	out := &DeleteSettingsOutput{}
	id := defaultString(in.Id, "settings")

	query := fmt.Sprintf(`mutation DeleteSettings {
		deleteSettings(input: { id: %s }) {
			id
		}
	}
	`, graphqlString(id))

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
