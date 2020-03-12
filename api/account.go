package api

import (
	"encoding/json"
	"net/http"

	"github.com/delivc/team/models"
	"github.com/delivc/team/storage"
	"github.com/go-chi/chi/v4"
	"github.com/gofrs/uuid"
)

// accountCreateParams
type accountCreateParams struct {
	Name string `json:"name"`
	Aud  string `json:"-"`
}

// AccountCreate trys to create a new Account
// [POST]/accounts {accountCreateParams}
func (a *API) AccountCreate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	instanceID := getInstanceID(ctx)
	params := &accountCreateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	user := getUser(ctx)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("Could not read params: %v", err)
	}

	params.Aud = a.requestAud(ctx, r)

	var account *models.Account

	err = a.db.Transaction(func(conn *storage.Connection) error {
		var terr error
		account, terr = models.NewAccount(instanceID, params.Name, params.Aud)
		if terr != nil {
			return internalServerError("Database error creating account").WithInternalError(err)
		}

		if account.OwnerIDs == nil {
			account.OwnerIDs = make(map[string]interface{})
		}

		account.OwnerIDs["0"] = user.ID

		err = conn.Transaction(func(tx *storage.Connection) error {
			// create account
			if txerr := tx.Create(account); txerr != nil {
				return internalServerError("Database error saving new account").WithInternalError(txerr)
			}
			// create admin role
			permissions, err := models.AllPermissions(tx)
			if err != nil {
				return err
			}

			admin, err := models.NewRole(account.ID, "Admin")
			if err != nil {
				return err
			}
			admin.Permissions = permissions

			if txerr := tx.Create(admin); txerr != nil {
				return internalServerError("Database error saving new admin role").WithInternalError(txerr)
			}

			// attach user to account
			if txerr := models.AttachUserToAccount(tx, user.ID, account.ID, admin.ID); txerr != nil {
				return internalServerError("Database error attaching user to account").WithInternalError(txerr)
			}

			// do we need more default roles?, ehehehe
			// nope: this is related to the Owner of the Account
			account.Roles = []models.Role{*admin}

			return nil
		})

		return err
	})
	if err != nil {
		return err
	}
	// we can cache this on creation
	// "users" will not exists, but it will get updated
	// with the next "update" until then, data is fine.
	a.cache.SetDefault("account-"+account.ID.String(), account)
	return sendJSON(w, http.StatusOK, account)
}

// AccountDelete deletes an Account from the Storage
// The user who is calling "delete" must fullfill:
// a: is Superadmin
// b: is one of the owners of account
// [DELETE]/accounts/{id}
func (a *API) AccountDelete(w http.ResponseWriter, r *http.Request) error {
	accountID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return badRequestError("Invalid Account ID")
	}
	ctx := r.Context()
	user := getUser(ctx)

	if user.IsSuperAdmin {
		// if the user is a super admin
		// we are not doing any more validation
		// just kick'em off directly
		if _, err := models.DeleteAccount(a.db, accountID); err != nil {
			return internalServerError("Database error deleting account").WithInternalError(err)
		}
	}
	// if given user is not super admin
	// we have to make sure that he has proper permission (is_owner)
	// 1: get the account from storage
	account, err := models.FindAccountByID(a.db, accountID)
	if err != nil {
		if models.IsNotFoundError(err) {
			return notFoundError(err.Error())
		}
		return internalServerError("Database error finding account").WithInternalError(err)
	}

	if !account.IsOwner(user.ID) {
		return unauthorizedError("You dont have proper permission")
	}

	if _, err := models.DeleteAccount(a.db, accountID); err != nil {
		return internalServerError("Database error deleting account").WithInternalError(err)
	}

	return sendJSON(w, http.StatusOK, map[string]interface{}{})
}

// AccountsGet returns a list of all related accounts
func (a *API) AccountsGet(w http.ResponseWriter, r *http.Request) error {
	// what is our caching key?
	// what are we receiving?
	// nothing ;s we can just save on a user base ;O
	// [accounts-userid?]
	ctx := r.Context()
	aud := a.requestAud(ctx, r)
	user := getUser(ctx)

	var userID uuid.UUID
	if !user.IsSuperAdmin {
		userID = user.ID
	}
	pageParams, err := paginate(r)
	if err != nil {
		return badRequestError("Bad Pagination Parameters: %v", err)
	}

	sortParams, err := sort(r, map[string]bool{models.CreatedAt: true}, []models.SortField{models.SortField{Name: models.CreatedAt, Dir: models.Descending}})
	if err != nil {
		return badRequestError("Bad Sort Parameters: %v", err)
	}

	accounts, err := models.FindAccounts(a.db, userID, pageParams, sortParams)
	if err != nil {
		return internalServerError("Database error finding accounts").WithInternalError(err)
	}
	addPaginationHeaders(w, r, pageParams)

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"accounts": accounts,
		"aud":      aud,
	})
}

// AccountGet gets an Account from the Database if exists
// or returns nil
// single requests are cached!
func (a *API) AccountGet(w http.ResponseWriter, r *http.Request) error {
	var accountID uuid.UUID
	var account *models.Account
	var err error

	ctx := r.Context()
	user := getUser(ctx)

	accountID, err = uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return badRequestError("Invalid Account ID")
	}

	fromCache, exists := a.cache.Get("account-" + accountID.String())
	if exists {
		var ok bool
		account, ok = fromCache.(*models.Account)
		if ok {
			// do we have permission to view this entity
			if user.IsSuperAdmin || account.IsOwner(user.ID) || account.IsMember(user.ID) {
				// we are just reading, this is a default permission
				// so simple is this.
				return sendJSON(w, http.StatusOK, account)
			}
			// usually we would throw an 401
			// but this would give attackers an idea about the existence of this
			// account

			return notFoundError("Account not found")
		}
	}

	account, err = models.FindAccountByID(a.db, accountID)
	if err != nil {
		if models.IsNotFoundError(err) {
			return notFoundError(err.Error())
		}
		return internalServerError("Database error finding account").WithInternalError(err)
	}

	// cache it
	a.cache.SetDefault("account-"+account.ID.String(), account)

	if user.IsSuperAdmin || account.IsOwner(user.ID) || account.IsMember(user.ID) {
		return sendJSON(w, http.StatusOK, account)
	}
	return notFoundError("Account not found")
}

type accountUpdateParams struct {
	Name           string         `json:"name"`
	BillingName    string         `json:"billing_name"`
	BillingEmail   string         `json:"billing_email"`
	BillingDetails string         `json:"billing_details"`
	MetaData       models.JSONMap `json:"account_meta_data"`
}

// AccountsUpdate updates given account if proper permission
func (a *API) AccountsUpdate(w http.ResponseWriter, r *http.Request) error {
	var accountID uuid.UUID
	var account *models.Account
	var err error
	accountID, err = uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return badRequestError("Invalid Account ID")
	}

	ctx := r.Context()
	params := &accountUpdateParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err = jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("Could not read Account Update params: %v", err)
	}

	user := getUser(ctx)
	if user == nil {
		return internalServerError("Error finding user object")
	}

	// get the account,check permissions
	// check account cache
	var ok bool
	fromCache, exists := a.cache.Get("account-" + accountID.String())
	if exists {
		if account, ok = fromCache.(*models.Account); !ok {
			account, err = models.FindAccountByID(a.db, accountID)
			if err != nil {
				if models.IsNotFoundError(err) {
					return notFoundError(err.Error())
				}
				return internalServerError("Database error finding account").WithInternalError(err)
			}
		}
	} else {
		account, err = models.FindAccountByID(a.db, accountID)
		if err != nil {
			if models.IsNotFoundError(err) {
				return notFoundError(err.Error())
			}
			return internalServerError("Database error finding account").WithInternalError(err)
		}
	}
	err = a.db.Transaction(func(tx *storage.Connection) error {
		var terr error
		// get permissions eg. hasPermission
		if account.HasPermissionTo(tx, "account-edit") || account.IsOwner(user.ID) || user.IsSuperAdmin {
			if params.Name != "" {
				if terr = account.UpdateName(tx, params.Name); terr != nil {
					return internalServerError("Error during name change").WithInternalError(terr)
				}
			}
			if params.BillingName != "" {
				if terr = account.UpdateBillingName(tx, params.BillingName); terr != nil {
					return internalServerError("Error during billing name change").WithInternalError(terr)
				}
			}
			if len(params.BillingEmail) > 254 || params.BillingEmail != "" {
				if !emailRegex.MatchString(params.BillingEmail) {
					return unprocessableEntityError("Email is Invalid")
				}
				// no further validation happens here
				if terr = account.UpdateBillingEmail(tx, params.BillingEmail); terr != nil {
					return internalServerError("Error during billing email change").WithInternalError(terr)
				}
			}
			if params.BillingDetails != "" {
				if terr = account.UpdateBillingDetails(tx, params.BillingDetails); terr != nil {
					return internalServerError("Error during billing details change").WithInternalError(terr)
				}
			}
			if params.MetaData != nil {
				if !user.IsSuperAdmin {
					return unauthorizedError("Updating account_meta_data requires admin privileges")
				}
				if terr = account.UpdateAccountMetaData(tx, params.MetaData); terr != nil {
					return internalServerError("Error updating user").WithInternalError(terr)
				}
			}

			// TODO: audit!
			// we should remove the audit from identity and team
			// and create a service by its own
			return nil
		}
		return unauthorizedError("You dont have `account-edit` Permission, ask your Manager")
	})

	if err != nil {
		return err
	}

	// cache it
	a.cache.SetDefault("account-"+account.ID.String(), account)

	return sendJSON(w, http.StatusOK, account)
}
