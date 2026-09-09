package awsteam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/smithy-go/ptr"
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
		return nil, errors.New("id is required to create Approvers")
	}

	approversJson, err := json.Marshal(in.Approvers)

	if err != nil {
		return nil, err
	}

	groupIdsJson, err := json.Marshal(in.GroupIds)

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`mutation CreateApprovers {
		createApprovers(
			input: {
				id: "%s"
				name: "%s"
				type: "%s"
				approvers: %s
				groupIds: %s
				ticketNo: "%s"
				modifiedBy: "%s"
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
	}`, ptr.ToString(in.Id),
		ptr.ToString(in.Name),
		ptr.ToString(in.Type),
		string(approversJson),
		string(groupIdsJson),
		ptr.ToString(in.TicketNo),
		ptr.ToString(in.ModifiedBy),
	)

	raw, err := client.GraphClient.ExecRaw(ctx, q, nil)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, out)

	if err != nil {
		return nil, err
	}

	return out, nil
}
