package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/delivc/team/conf"
	"github.com/delivc/team/models"
	"github.com/go-chi/chi/v4"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func addRequestID(globalConfig *conf.GlobalConfiguration) middlewareHandler {
	return func(w http.ResponseWriter, r *http.Request) (context.Context, error) {
		id := ""
		if globalConfig.API.RequestIDHeader != "" {
			id = r.Header.Get(globalConfig.API.RequestIDHeader)
		}
		if id == "" {
			uid, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			id = uid.String()
		}

		ctx := r.Context()
		ctx = withRequestID(ctx, id)
		return ctx, nil
	}
}

func sendJSON(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(b)
	return err
}

func (a *API) getAccountFromRequest(r *http.Request) (*models.Account, error) {
	var accountID uuid.UUID
	var err error
	var account *models.Account

	accountID, err = uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		return account, badRequestError("Invalid Account ID")
	}

	fromCache, exists := a.cache.Get("account-" + accountID.String())
	if exists {
		var ok bool
		account, ok := fromCache.(*models.Account)
		if !ok {
			account, err = models.FindAccountByID(a.db, accountID)
			if err != nil {
				if models.IsNotFoundError(err) {
					return account, notFoundError(err.Error())
				}
				return account, internalServerError("Database error finding account").WithInternalError(err)
			}
		}
	} else {
		account, err = models.FindAccountByID(a.db, accountID)
		if err != nil {
			if models.IsNotFoundError(err) {
				return account, notFoundError(err.Error())
			}
			return account, internalServerError("Database error finding account").WithInternalError(err)
		}
	}

	return account, nil
}
