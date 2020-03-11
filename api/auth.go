package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/delivc/identity/models"
	jwt "github.com/dgrijalva/jwt-go"
)

// auth.go
// validates an users auth state
// and caches it, until the token expires
var (
	cache = authCache{}
)

type authCacheItem struct {
	User      *models.User
	ExpiresAt time.Time `json:"exp"`
}
type authCache struct {
	Items map[string]*authCacheItem
	mutex sync.Mutex
}

func init() {
	// creates a new global authCache
	cache = authCache{
		Items: map[string]*authCacheItem{},
	}

	// cache clearing
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// run cache clearer
				for tk := range cache.Items {
					if cache.Items[tk].ExpiresAt.Before(time.Now()) {
						cache.mutex.Lock()
						delete(cache.Items, tk)
						cache.mutex.Unlock()
					}

				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// requireAuthentication is a middleware to check if the user who made the
// request is authenticated with our Identity Service
func (a *API) requireAuthentication(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	token, err := a.extractBearerToken(w, r)
	if err != nil {
		// user is maybe not authentified
		return nil, err
	}
	return a.validateToken(token, r, w)
}

func (a *API) validateToken(bearer string, r *http.Request, w http.ResponseWriter) (context.Context, error) {
	var user models.User
	ctx := r.Context()
	// check if the token is already cached
	cached, ok := cache.Items[bearer]
	if ok {
		if cached.ExpiresAt.After(time.Now()) {
			return withUser(ctx, cache.Items[bearer].User), nil
		}
	}

	// token is not cached, heck!
	// we have to ask the Identity Service if our user is valid
	client := &http.Client{}
	// get endpoint from config
	request, _ := http.NewRequest("GET", a.config.IdentityEndpoint+"/user", nil)
	request.Header.Add("Authorization", "Bearer "+bearer)
	resp, err := client.Do(request)
	if err != nil {
		return nil, unauthorizedError("Invalid token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, unauthorizedError("Invalid token: token is expired")
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, unauthorizedError("Invalid token: %v", err)
	}
	claims := jwt.StandardClaims{}
	p := jwt.Parser{
		ValidMethods:         []string{jwt.SigningMethodHS256.Name},
		SkipClaimsValidation: true,
	}
	_, _, err = p.ParseUnverified(bearer, &claims)
	if err != nil {
		return nil, unauthorizedError("Invalid token: Token is does not match Schema")
	}

	cached = &authCacheItem{
		User:      &user,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}

	cache.mutex.Lock()
	cache.Items[bearer] = cached
	cache.mutex.Unlock()

	return withUser(ctx, &user), nil
}

func (a *API) extractBearerToken(w http.ResponseWriter, r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", unauthorizedError("This endpoint requires a Bearer token")
	}

	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", unauthorizedError("This endpoint requires a Bearer token")
	}

	return matches[1], nil
}
