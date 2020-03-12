package models

import (
	"database/sql"
	"time"

	"github.com/delivc/team/storage"
	"github.com/delivc/team/storage/namespace"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Permission exports `Id` and `Match`
type Permission struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
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

// FindPermissionsByName returns permissions
func FindPermissionsByName(tx *storage.Connection, request []string) ([]Permission, error) {
	permissions := []Permission{}
	if len(request) == 0 {
		return permissions, errors.New("No Permissions")
	}
	q := tx.Q()

	q.Where("name IN (?)", request)

	if err := q.All(&permissions); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, PermissionNotFoundError{}
		}
		return permissions, errors.Wrap(err, "error finding permissions")
	}

	logrus.Info(permissions)

	//

	return permissions, nil
}

// FindPermissions returns a list of Permissions if any
func FindPermissions(tx *storage.Connection, pageParams *Pagination, sortParams *SortParams) ([]*Permission, error) {
	permissions := []*Permission{}

	q := tx.Q()

	if sortParams != nil && len(sortParams.Fields) > 0 {
		for _, field := range sortParams.Fields {
			q = q.Order(field.Name + " " + string(field.Dir))
		}
	}

	var err error
	if pageParams != nil {
		err = q.Paginate(int(pageParams.Page), int(pageParams.PerPage)).All(&permissions)
		pageParams.Count = uint64(q.Paginator.TotalEntriesSize)
	} else {
		err = q.All(&permissions)
	}

	return permissions, err
}
