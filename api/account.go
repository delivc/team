package api

import (
	"encoding/json"
	"net/http"

	"github.com/delivc/team/models"
	"github.com/delivc/team/storage"
	"github.com/go-chi/chi/v4"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
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

	pop.Debug = true

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
func (a *API) AccountGet(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "Error generating unique id")
	}
	account, err := models.FindAccountByID(a.db, id)
	if err != nil {
		if models.IsNotFoundError(err) {
			return notFoundError(err.Error())
		}
		return internalServerError("Database error finding user").WithInternalError(err)
	}
	return sendJSON(w, http.StatusOK, account)
}
