package awsteam

import (
	"context"
	"errors"
	"fmt"
)

type CreateApproversInput struct {
	Id         *string   `json:"id"`
	Name       *string   `json:"name"`
	Type       *string   `json:"type"`
	Approvers  []*string `json:"approvers"`
	GroupIds   []*string `json:"groupIds"`
	TicketNo   *string   `json:"ticketNo"`
	ModifiedBy *string   `json:"modifiedBy"`
}

type CreateApproversOutput struct {
	Approvers *Approvers `json:"createApprovers"`
}

func (client *Client) CreateApprovers(ctx context.Context, in *CreateApproversInput) (*CreateApproversOutput, error) {
	out := &CreateApproversOutput{}

	if in.Id == nil {
		return nil, errors.New("Id is required to create Approvers.")
	}

	approversJSON, err := graphqlJSON(in.Approvers)
	if err != nil {
		return nil, err
	}

	groupIDsJSON, err := graphqlJSON(in.GroupIds)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`mutation CreateApprovers {
		createApprovers(
			input: {
				id: %s
				name: %s
				type: %s
				approvers: %s
				groupIds: %s
				ticketNo: %s
				modifiedBy: %s
			}
		)  {
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
	}`, graphqlStringPtr(in.Id),
		graphqlStringPtr(in.Name),
		graphqlStringPtr(in.Type),
		approversJSON,
		groupIDsJSON,
		graphqlStringPtr(in.TicketNo),
		graphqlStringPtr(in.ModifiedBy),
	)

	if err := client.executeGraphQL(ctx, query, nil, out); err != nil {
		return nil, err
	}

	return out, nil
}
