package api

import (
	"net/http"

	"github.com/delivc/team/models"
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
		roleFromCache, exists := a.cache.Get("role-"+roleID.String())
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

	rolesFromCache, exists := a.cache.Get("roles-"+accountID.String())
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
