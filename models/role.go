package models

import (
	"database/sql"
	"time"

	"github.com/delivc/team/storage"
	"github.com/delivc/team/storage/namespace"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Role reflects a given Role within an Account
type Role struct {
	AccountID   uuid.UUID   `json:"-" db:"account_id"`
	ID          uuid.UUID   `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time   `json:"updatedAt" db:"updated_at"`
	Permissions Permissions `json:"permissions,omitempty" many_to_many:"roles_permissions"`
}

// RolePermission relationship between roles and permissions
type RolePermission struct {
	ID           uuid.UUID `json:"id" db:"id"`
	RoleID       uuid.UUID
	PermissionID uuid.UUID
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// TableName returns the given tablename of the model
func (RolePermission) TableName() string {
	tableName := "roles_permissions"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// TableName returns the given tablename of the model
func (Role) TableName() string {
	tableName := "roles"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// NewRole creates a new Role
func NewRole(accountID uuid.UUID, name string) (*Role, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating unique id")
	}

	role := &Role{
		AccountID: accountID,
		ID:        id,
		Name:      name,
	}

	return role, nil
}

// AttachRole sets the Role of the user within the Account Space
// The relationship between the User and Account MUST exists when calling this
func AttachRole(tx *storage.Connection, accountID uuid.UUID, userID uuid.UUID, roleID uuid.UUID) error {
	return nil
}

func findRoles(tx *storage.Connection, query string, args ...interface{}) ([]*Role, error) {
	obj := []*Role{}
	if err := tx.Q().Eager().Where(query, args...).All(&obj); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, RoleNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding roles")
	}
	return obj, nil
}

// FindRolesByAccount returns a list of roles by account or error
func FindRolesByAccount(tx *storage.Connection, id uuid.UUID) ([]*Role, error) {
	return findRoles(tx, "account_id = ?", id)
}
