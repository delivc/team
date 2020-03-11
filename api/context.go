package api

import (
	"context"

	"github.com/delivc/identity/models"
	"github.com/delivc/team/conf"
	"github.com/gofrs/uuid"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	userKey       = contextKey("user")
	configKey     = contextKey("config")
	instanceIDKey = contextKey("instance_id")
	instanceKey   = contextKey("instance")
	requestIDKey  = contextKey("request_id")
)

// withUser adds the JWT token to the context.
func withUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// getUser reads the JWT token from the context.
func getUser(ctx context.Context) *models.User {
	obj := ctx.Value(userKey)
	if obj == nil {
		return nil
	}

	return obj.(*models.User)
}

// withRequestID adds the provided request ID to the context.
func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// getRequestID reads the request ID from the context.
func getRequestID(ctx context.Context) string {
	obj := ctx.Value(requestIDKey)
	if obj == nil {
		return ""
	}

	return obj.(string)
}

// withConfig adds the tenant configuration to the context.
func withConfig(ctx context.Context, config *conf.Configuration) context.Context {
	return context.WithValue(ctx, configKey, config)
}

// Reads the configuration from context
func getConfig(ctx context.Context) *conf.Configuration {
	obj := ctx.Value(configKey)
	if obj == nil {
		return nil
	}
	return obj.(*conf.Configuration)
}

// withInstanceID adds the instance id to the context.
func withInstanceID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, instanceIDKey, id)
}

// getInstanceID reads the instance id from the context.
func getInstanceID(ctx context.Context) uuid.UUID {
	obj := ctx.Value(instanceIDKey)
	if obj == nil {
		return uuid.Nil
	}
	return obj.(uuid.UUID)
}
