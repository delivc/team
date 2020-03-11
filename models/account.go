package models

import (
	"database/sql"
	"time"

	"github.com/delivc/identity/models"
	"github.com/delivc/team/storage"
	"github.com/delivc/team/storage/namespace"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Account represents a Team or a Company within the Delivc Org
type Account struct {
	InstanceID uuid.UUID `json:"-" db:"instance_id"`
	ID         uuid.UUID `json:"id" db:"id"`

	Aud             string  `json:"aud" db:"aud"`
	Name            string  `json:"name" db:"name"`
	BillingName     string  `json:"billing_name,omitempty" db:"billing_name"`
	BillingEmail    string  `json:"billing_email,omitempty" db:"billing_email"`
	BillingDetails  string  `json:"billing_details,omitempty" db:"billing_details"`
	BillingPeriod   string  `json:"billing_period,omitempty" db:"billing_period"`
	PaymentMethodID string  `json:"payment_method_id,omitempty" db:"payment_method_id"`
	OwnerIDs        JSONMap `json:"owner_ids" db:"raw_owner_ids"`

	AccountMetaData JSONMap `json:"account_metadata,omitempty" db:"raw_account_meta_data"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`

	Roles []Role        `json:"roles,omitempty" has_many:"roles"`
	Users []models.User `json:"users,omitempty" many_to_many:"accounts_users"`
}

// TableName returns the given tablename of the model
func (Account) TableName() string {
	tableName := "accounts"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "_" + tableName
	}

	return tableName
}

// IsOwner checks if given uid is owner of account
func (a *Account) IsOwner(userID uuid.UUID) bool {
	for _, value := range a.OwnerIDs {
		if value == userID.String() {
			return true
		}
	}
	return false

}

// NewAccount initializes a new account from name
// does not create!!! an account in the database
func NewAccount(instanceID uuid.UUID, name, aud string) (*Account, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating unique id")
	}

	account := &Account{
		InstanceID: instanceID,
		ID:         id,
		Aud:        aud,
		Name:       name,
	}
	return account, nil
}

// DeleteAccount removes an account from storage
func DeleteAccount(tx *storage.Connection, accountID uuid.UUID) (bool, error) {
	if err := tx.Destroy(&Account{ID: accountID}); err != nil {
		return false, err
	}
	return true, nil
}

func findAccount(tx *storage.Connection, query string, args ...interface{}) (*Account, error) {
	obj := &Account{}
	if err := tx.Q().Where(query, args...).First(obj); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, AccountNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding account")
	}
	return obj, nil
}

// FindAccountByID finds a account matching the provided ID.
func FindAccountByID(tx *storage.Connection, id uuid.UUID) (*Account, error) {
	return findAccount(tx, "id = ?", id)
}

// FindAccounts searches for Accounts in the given "Audience"
func FindAccounts(tx *storage.Connection, userID uuid.UUID, pageParams *Pagination, sortParams *SortParams) ([]*Account, error) {
	accounts := []*Account{}
	var err error

	pop.Debug = true
	q := tx.Q()
	if userID.String() != "00000000-0000-0000-0000-000000000000" {
		// UserID is not nil, so we have to query for the relations from
		// account_user
		q.RawQuery(`
		SELECT
			accounts.id as id,
			accounts.name as name,
			accounts.billing_name as billing_name,
			accounts.billing_email as billing_email,
			accounts.billing_details as billing_details,
			accounts.billing_period as billing_period,
			accounts.payment_method_id as payment_method_id,
			accounts.raw_owner_ids as raw_owner_ids,
			accounts.raw_account_meta_data as raw_account_meta_data,
			accounts.created_at as created_at,
			accounts.updated_at as updated_at
		FROM
			accounts_users as accounts_users
			JOIN accounts ON accounts.id = accounts_users.account_id
		WHERE
			accounts_users.user_id = ?`, userID)

		err = q.Eager("Roles").All(&accounts)
		return accounts, err
	}

	if sortParams != nil && len(sortParams.Fields) > 0 {
		for _, field := range sortParams.Fields {
			q = q.Order(field.Name + " " + string(field.Dir))
		}
	}

	if pageParams != nil {
		err = q.Paginate(int(pageParams.Page), int(pageParams.PerPage)).Eager("Roles").All(&accounts)
		pageParams.Count = uint64(q.Paginator.TotalEntriesSize)
	} else {
		err = q.Eager("Roles").All(&accounts)
	}
	return accounts, err
}
