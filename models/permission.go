package models

import (
	"time"

	"github.com/delivc/team/storage"
	"github.com/delivc/team/storage/namespace"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Permission exports `Id` and `Match`
type Permission struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Roles     []Role    `json:"-" many_to_many:"roles_permissions"`
}

// Permissions is a list of Permissions
type Permissions []Permission

// TableName returns the given tablename of the model
func (Permission) TableName() string {
	tableName := "permissions"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// NewPermission creates a new Permission
func NewPermission(name string) (*Permission, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating unique id")
	}

	permission := &Permission{
		ID:   id,
		Name: name,
	}
	return permission, nil
}

// AllPermissions returns all Permissions from the Database
func AllPermissions(tx *storage.Connection) (permissions []Permission, err error) {
	err = tx.All(&permissions)
	if err != nil {
		return nil, errors.Wrap(err, "error finding permissions")
	}
	return permissions, nil
}
