package awsteam

import (
	"context"
)

type GetAccountsInput struct{}

type GetAccountsOutput struct {
	Accounts []*Account `json:"getAccounts"`
}

func (client *Client) GetAccounts(ctx context.Context, in *GetAccountsInput) (*GetAccountsOutput, error) {
	out := &GetAccountsOutput{}

	query := `query GetAccounts {
		getAccounts {
			name
			id
		}
	}`

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
