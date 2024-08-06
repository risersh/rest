package invitations

import "time"

type InvitationType string

var (
	InvitationTypeCamera   InvitationType = "camera"
	InvitationTypeLocation InvitationType = "location"
	InvitationTypeTenant   InvitationType = "organization"
)

type InvitationRole string

var (
	InvitationRoleOwner   InvitationRole = "owner"
	InvitationRoleManager InvitationRole = "manager"
	InvitationRoleViewer  InvitationRole = "viewer"
)

type Invitation struct {
	ID      string    `json:"id"`
	Email   string    `json:"email"`
	Type    string    `json:"type"`
	Context string    `json:"context"`
	Role    string    `json:"role"`
	Created time.Time `json:"created"`
	Status  string    `json:"status"`
}
