package awsteam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/smithy-go/ptr"
)

type UpdateApproversInput struct {
	Id         *string   `json:"id"`
	Name       *string   `json:"name"`
	Type       *string   `json:"type"`
	Approvers  []*string `json:"approvers"`
	GroupIds   []*string `json:"groupIds"`
	TicketNo   *string   `json:"ticketNo"`
	ModifiedBy *string   `json:"modifiedBy"`
}

type UpdateApproversOutput struct {
	Approvers *Approvers `json:"updateApprovers"`
}

func (client *Client) UpdateApprovers(ctx context.Context, in *UpdateApproversInput) (*UpdateApproversOutput, error) {
	out := &UpdateApproversOutput{}

	if in.Id == nil {
		return nil, errors.New("id is required to update Approvers")
	}

	approversJson, err := json.Marshal(in.Approvers)

	if err != nil {
		return nil, err
	}

	groupIdsJson, err := json.Marshal(in.GroupIds)

	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`mutation UpdateApprovers {
		updateApprovers(
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
			updatedAt
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
