package api

import (
	"encoding/json"
	"net/http"

	"github.com/delivc/team/models"
	"github.com/delivc/team/storage"
	"github.com/go-chi/chi/v4"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
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

	account, err := a.getAccountFromRequest(r)
	if err != nil {
		return err
	}

	user := getUser(ctx)
	if user == nil {
		return badRequestError("Invalid User")
	}
	if user.IsSuperAdmin || account.IsOwner(user.ID) || account.HasPermissionTo(a.db, "account-role-create", user.ID) {
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
	return unauthorizedError("You dont have `account-role-create` Permission, ask your Manager")
}

type updateRoleRequest struct {
	createRoleRequest
}

// RoleUpdate updates a role
// Permission: account-role-update
func (a *API) RoleUpdate(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := &updateRoleRequest{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return badRequestError("Could not read Account Update params: %v", err)
	}

	account, err := a.getAccountFromRequest(r)
	if err != nil {
		return err
	}

	user := getUser(ctx)
	if user == nil {
		return badRequestError("Invalid User")
	}

	// get role from request
	roleID, err := uuid.FromString(chi.URLParam(r, "roleId"))
	if err != nil {
		return badRequestError("Invalid Role ID")
	}
	logrus.Info(roleID)

	// get role from cache or get a new from db
	var role *models.Role
	roleFromCache, exists := a.cache.Get("role-" + roleID.String())
	if exists {
		role = roleFromCache.(*models.Role)
	} else {
		if role, err = models.FindRoleByAccountAndID(a.db, account.ID, roleID); err != nil {
			return internalServerError("Database error finding roles").WithInternalError(err)
		}
	}

	if user.IsSuperAdmin || account.IsOwner(user.ID) || account.HasPermissionTo(a.db, "account-role-update", user.ID) {
		// we have permission, now do the updates :)))

		err = a.db.Transaction(func(conn *storage.Connection) error {
			var terr error
			if params.Name != "" {
				if terr = role.UpdateName(conn, params.Name); terr != nil {
					return internalServerError("Error during name change").WithInternalError(terr)
				}
			}
			if params.Permissions != nil {
				if terr = role.UpdatePermissions(conn, params.Permissions); terr != nil {
					return internalServerError("Error updating permissions").WithInternalError(terr)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		a.cache.SetDefault("role-"+role.ID.String(), role)
		return sendJSON(w, 200, role)
	}

	return unauthorizedError("You dont have `account-role-update` Permission, ask your Manager")
}

// RoleDestroy destroys a role in storage
// Permission: account-role-destroy
func (a *API) RoleDestroy(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	account, err := a.getAccountFromRequest(r)
	if err != nil {
		return err
	}

	user := getUser(ctx)
	if user == nil {
		return badRequestError("Invalid User")
	}

	// get role from request
	roleID, err := uuid.FromString(chi.URLParam(r, "roleId"))
	if err != nil {
		return badRequestError("Invalid Role ID")
	}

	if user.IsSuperAdmin || account.IsOwner(user.ID) || account.HasPermissionTo(a.db, "account-role-destroy", user.ID) {
		err = a.db.Transaction(func(conn *storage.Connection) error {
			return models.DeleteRole(conn, roleID)
		})

		if err != nil {
			return err
		}

		// remove from cache if exists
		a.cache.Delete("role-" + roleID.String())

		return sendJSON(w, http.StatusOK, map[string]interface{}{})
	}

	return unauthorizedError("You dont have `account-role-destroy` Permission, ask your Manager")
}

// detachPermissions
// detachaches all Permissions of given Role
func detachPermissions(tx *storage.Connection, roleID uuid.UUID) error {
	return nil
}
