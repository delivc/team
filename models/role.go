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
	RoleID       uuid.UUID `db:"role_id"`
	PermissionID uuid.UUID `db:"permission_id"`
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

// UpdateName updates the name of the role
func (r *Role) UpdateName(tx *storage.Connection, newName string) error {
	if newName == "" {
		return errors.New("Error: invalid name")
	}
	r.Name = newName
	return tx.UpdateOnly(r, "name", "updated_at")
}

// UpdatePermissions syncs given permissions of role
func (r *Role) UpdatePermissions(tx *storage.Connection, perms []string) error {
	if perms == nil {
		return errors.New("Error: invalid Permissions")
	}
	var err error
	var permissions []Permission

	if permissions, err = FindPermissionsByName(tx, perms); err != nil {
		return errors.Wrap(err, "Error finding Permissions")
	}

	if err = detachAllPermissions(tx, r.ID); err != nil {
		return err
	}

	for _, permission := range permissions {
		if err = attachPermission(tx, r.ID, permission.ID); err != nil {
			return err
		}
	}
	r.Permissions = permissions

	return nil
}

func attachPermission(tx *storage.Connection, roleID uuid.UUID, permissionID uuid.UUID) error {
	id, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "Error generating unique id")
	}
	p := RolePermission{
		ID:           id,
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return tx.Create(&p)
}

func detachAllPermissions(tx *storage.Connection, roleID uuid.UUID) error {
	tableName := RolePermission{}.TableName()

	if err := tx.RawQuery("DELETE FROM "+tableName+" WHERE role_id = ?", roleID).Exec(); err != nil {
		return err
	}
	return nil
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

func findRole(tx *storage.Connection, query string, args ...interface{}) (*Role, error) {
	obj := &Role{}
	if err := tx.Q().Eager().Where(query, args...).First(obj); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, RoleNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding role")
	}
	return obj, nil
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

// FindRoleByAccountAndID returns roles by account and id
func FindRoleByAccountAndID(tx *storage.Connection, accountID uuid.UUID, roleID uuid.UUID) (*Role, error) {
	return findRole(tx, "account_id = ? and id = ?", accountID, roleID)
}

// DeleteRole destroys a role in storage
func DeleteRole(tx *storage.Connection, roleID uuid.UUID) error {
	return tx.Destroy(&Role{ID: roleID})
}
