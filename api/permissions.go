package api

import (
	"net/http"

	"github.com/delivc/team/models"
)

// PermissionsGet returns a list of all Permissions
func (a *API) PermissionsGet(w http.ResponseWriter, r *http.Request) error {
	// this is available for all who are signed in
	// no security checks

	pageParams, err := paginate(r)
	if err != nil {
		return badRequestError("Bad Pagination Parameters: %v", err)
	}

	sortParams, err := sort(r, map[string]bool{models.CreatedAt: true}, []models.SortField{models.SortField{Name: models.CreatedAt, Dir: models.Descending}})
	if err != nil {
		return badRequestError("Bad Sort Parameters: %v", err)
	}

	permissions, err := models.FindPermissions(a.db, pageParams, sortParams)
	if err != nil {
		return internalServerError("Database error finding permissions").WithInternalError(err)
	}
	addPaginationHeaders(w, r, pageParams)
	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"permissions": permissions,
	})
}
