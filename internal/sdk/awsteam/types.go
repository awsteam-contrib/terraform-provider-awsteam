package awsteam

type Account struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type Approvers struct {
	Id         *string   `json:"id"`        // "Account ID" or "OU id"
	Name       *string   `json:"name"`      // Account Name or OU Name
	Type       *string   `json:"type"`      // "Account" or "OU"
	Approvers  []*string `json:"approvers"` // List of Group names that can approve
	GroupIds   []*string `json:"groupIds"`  // list of group ids that can approve
	TicketNo   *string   `json:"ticketNo"`
	ModifiedBy *string   `json:"modifiedBy"`
	CreatedAt  *string   `json:"createdAt"`
	UpdatedAt  *string   `json:"updatedAt"`
}

type Eligibility struct {
	Id               *string                  `json:"id"` // User or Group Id
	Name             *string                  `json:"name"`
	Type             *string                  `json:"type"` // "User" or "Group"
	Accounts         []*EligibilityAccount    `json:"accounts"`
	OUs              []*EligibilityOU         `json:"ous"`
	Permissions      []*EligibilityPermission `json:"permissions"`
	TicketNo         *string                  `json:"ticketNo"`
	ApprovalRequired *bool                    `json:"approvalRequired"`
	Duration         *int64                   `json:"duration,string"`
	ModifiedBy       *string                  `json:"modifiedBy"`
	CreatedAt        *string                  `json:"createdAt"`
	UpdatedAt        *string                  `json:"updatedAt"`
}

type EligibilityAccount struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type EligibilityOU struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type EligibilityPermission struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type OU struct {
	Id       *string `json:"id"`
	Arn      *string `json:"arn"`
	Name     *string `json:"name"`
	Children []OU    `json:"children"`
}

type Permission struct {
	Name     *string
	Arn      *string
	Duration *string
}

type Settings struct {
	Approval                  *bool   `json:"approval"`
	Comments                  *bool   `json:"comments"`
	Duration                  *int64  `json:"duration,string"`
	Expiry                    *int64  `json:"expiry,string"`
	Id                        *string `json:"id"`
	SesNotificationsEnabled   *bool   `json:"sesNotificationsEnabled"`
	SnsNotificationsEnabled   *bool   `json:"snsNotificationsEnabled"`
	SlackNotificationsEnabled *bool   `json:"slackNotificationsEnabled"`
	SesSourceEmail            *string `json:"sesSourceEmail"`
	SesSourceArn              *string `json:"sesSourceArn"`
	SlackToken                *string `json:"slackToken"`
	TeamAdminGroup            *string `json:"teamAdminGroup"`
	TeamAuditorGroup          *string `json:"teamAuditorGroup"`
	TicketNo                  *bool   `json:"ticketNo"`
	UseOUCache                *bool   `json:"useOUCache"`
	ModifiedBy                *string `json:"modifiedBy"`
	CreatedAt                 *string `json:"createdAt"`
	UpdatedAt                 *string `json:"updatedAt"`
}
