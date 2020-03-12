package api

import (
	"encoding/json"
	"net/http"

	"github.com/delivc/team/models"
	"github.com/delivc/team/storage"
	"github.com/go-chi/chi/v4"
	"github.com/gofrs/uuid"
)

/**
 * Role attaches, detaches, CRUD
 * Roles with Permissions for an Account
 */

// RoleGet returns a list of roles or a single role
func (a *API) RoleGet(w http.ResponseWriter, r *http.Request) error {
	var err error
	var accountID uuid.UUID

	accountID, err = uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return badRequestError("Invalid Account ID")
	}

	// universal controller
	// accepts ID or nil
	roleIDString := chi.URLParam(r, "roleId")
	if roleIDString != "" {
		var roleID uuid.UUID
		var role *models.Role
		// create uuid from roleIDString
		roleID, err = uuid.FromString(roleIDString)
		if err != nil {
			return badRequestError("Invalid Role ID")
		}
		// check cache before query
		roleFromCache, exists := a.cache.Get("role-" + roleID.String())
		if exists {
			return sendJSON(w, http.StatusOK, roleFromCache)
		}

		role, err = models.FindRoleByAccountAndID(a.db, accountID, roleID)
		if err != nil {
			return internalServerError("Database error finding roles").WithInternalError(err)
		}
		a.cache.SetDefault("role-"+role.ID.String(), role)
		return sendJSON(w, http.StatusOK, role)
	}

	rolesFromCache, exists := a.cache.Get("roles-" + accountID.String())
	if exists {
		return sendJSON(w, http.StatusOK, map[string]interface{}{
			"roles": rolesFromCache,
		})
	}

	// check cache before query database
	var roles []*models.Role
	roles, err = models.FindRolesByAccount(a.db, accountID)
	if err != nil {
		return internalServerError("Database error finding roles").WithInternalError(err)
	}

	a.cache.SetDefault("roles-"+accountID.String(), roles)

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"roles": roles,
	})
}

type createRoleRequest struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// RoleCreate create a new role with permissions if given
func (a *API) RoleCreate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	// get account id from url
	// read params
	// bind params to new role object
	// why not bind directly to role object?
	// user could specifiy an ID from an other account
	// why not override id from request? yep...

	account, err := a.getAccountFromRequest(r)
	if err != nil {
		return err
	}

	user := getUser(ctx)
	if user == nil {
		return badRequestError("Invalid User")
	}
	if user.IsSuperAdmin || account.IsOwner(user.ID) || account.HasPermissionTo(a.db, "account-role-create", user.ID) {
		// we have a user
		// we have permission
		// we know the account where we create the new role
		params := &createRoleRequest{}
		jsonDecoder := json.NewDecoder(r.Body)
		err = jsonDecoder.Decode(params)
		if err != nil {
			return badRequestError("Could not read Role Update params: %v", err)
		}
		var role *models.Role
		err = a.db.Transaction(func(conn *storage.Connection) error {
			var terr error
			var permissions []models.Permission
			role, terr = models.NewRole(account.ID, params.Name)
			if role, terr = models.NewRole(account.ID, params.Name); terr != nil {
				return internalServerError("Database error creating role").WithInternalError(terr)
			}
			if params.Permissions != nil {
				if permissions, terr = models.FindPermissionsByName(conn, params.Permissions); terr != nil {
					return internalServerError("Database error creating role").WithInternalError(terr)
				}
				role.Permissions = permissions
			}

			if terr := conn.Create(role); terr != nil {
				return internalServerError("Database error saving new role").WithInternalError(terr)
			}
			return nil
		})

		if err != nil {
			return err
		}

		a.cache.SetDefault("role-"+role.ID.String(), role)

		return sendJSON(w, 200, role)

	}
	// get account from cache or db
	return unauthorizedError("You dont have `account-role-create` Permission, ask your Manager")
}

// Update a Role

// Destroy a Role

// detachPermissions
// detachaches all Permissions of given Role
func detachPermissions(tx *storage.Connection, roleID uuid.UUID) error {
	return nil
}
