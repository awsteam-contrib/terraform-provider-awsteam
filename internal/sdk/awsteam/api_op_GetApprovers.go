package awsteam

import (
	"context"
	"errors"
	"fmt"
)

type GetApproversInput struct {
	Id *string
}

type GetApproversOutput struct {
	Approvers *Approvers `json:"getApprovers"`
}

func (client *Client) GetApprovers(ctx context.Context, in *GetApproversInput) (*GetApproversOutput, error) {
	out := &GetApproversOutput{}

	if in.Id == nil {
		return nil, errors.New("Id is required to get Approvers.")
	}

	query := fmt.Sprintf(`query GetApprovers {
		getApprovers(id: %s) {
			id
			name
			type
			approvers
			groupIds
			ticketNo
			modifiedBy
			createdAt
			updatedAt
		}
	}
	`, graphqlStringPtr(in.Id))

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
