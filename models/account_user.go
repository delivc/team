package models

import (
	"time"

	"github.com/delivc/team/storage"
	"github.com/delivc/team/storage/namespace"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// AccountUser account type
// maintains the relationship between a User and the Account
type AccountUser struct {
	ID          uuid.UUID  `json:"-" db:"id"`
	AccountID   uuid.UUID  `json:"account_id" db:"account_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID      uuid.UUID  `json:"role_id" db:"role_id"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	InvitedAt   *time.Time `json:"invited_at,omitempty" db:"invited_at"`
	InvitedBy   uuid.UUID  `json:"invited_by,omitempty" db:"invited_by"`
}

// TableName returns the given tablename of the model
func (AccountUser) TableName() string {
	tableName := "accounts_users"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// AttachUserToAccount attaches a user to given account
func AttachUserToAccount(tx *storage.Connection, userID uuid.UUID, accountID uuid.UUID, roleID uuid.UUID) error {
	relation := AccountUser{
		AccountID: accountID,
		UserID:    userID,
		RoleID:    roleID,
	}
	if err := tx.Create(&relation); err != nil {
		return errors.Wrap(err, "Error generating attaching user to account")
	}
	return nil
}
