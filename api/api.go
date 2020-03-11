package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/delivc/team/conf"
	"github.com/delivc/team/storage"
	"github.com/go-chi/chi/v4"
	"github.com/gofrs/uuid"
	"github.com/rs/cors"
	"github.com/sebest/xff"
	"github.com/sirupsen/logrus"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

const (
	audHeaderName  = "X-JWT-AUD"
	defaultVersion = "unknown version"
)

// API is the main REST API
type API struct {
	handler http.Handler
	db      *storage.Connection
	config  *conf.GlobalConfiguration
	version string
}

// New creates a new API Instance
func New(ctx context.Context, globalConfig *conf.GlobalConfiguration, db *storage.Connection, version string) *API {
	api := &API{config: globalConfig, db: db, version: version}

	xffmw, _ := xff.Default()
	logger := newStructuredLogger(logrus.StandardLogger())

	r := newRouter()
	r.UseBypass(xffmw.Handler)
	r.Use(addRequestID(globalConfig))
	r.Use(recoverer)

	r.Get("/health", api.HealthCheck)

	r.Route("/", func(r *router) {
		r.UseBypass(logger)
		r.Use(api.requireAuthentication)

		r.Get("/accounts", api.AccountsGet)
		r.Post("/accounts", api.AccountCreate)
		r.Delete("/accounts/{id}", api.AccountDelete)
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://app.delivc.com", "http://app.delivc.com:8081"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	api.handler = corsHandler.Handler(chi.ServerBaseContext(ctx, r))

	return api
}

// ListenAndServe starts the REST API
func (a *API) ListenAndServe(hostAndPort string) {
	log := logrus.WithField("component", "api")
	server := &http.Server{
		Addr:    hostAndPort,
		Handler: a.handler,
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		waitForTermination(log, done)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		server.Shutdown(ctx)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Fatal("http server listen failed")
	}
}

// WaitForShutdown blocks until the system signals termination or done has a value
func waitForTermination(log logrus.FieldLogger, done <-chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-signals:
		log.Infof("Triggering shutdown from signal %s", sig)
	case <-done:
		log.Infof("Shutting down...")
	}
}

// WithInstanceConfig creates a new ctx with a config
func WithInstanceConfig(ctx context.Context, config *conf.Configuration, instanceID uuid.UUID) context.Context {
	ctx = withConfig(ctx, config)
	ctx = withInstanceID(ctx, instanceID)
	return ctx
}

// HealthCheck returns a "is a live"
func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, map[string]string{
		"version":     a.version,
		"name":        "Team",
		"description": "Team is a management Service for Delivc Teams",
	})
}

func (a *API) requestAud(ctx context.Context, r *http.Request) string {
	// First check for an audience in the header
	if aud := r.Header.Get(audHeaderName); aud != "" {
		return aud
	}

	return "app.delivc.com"
}

func (a *API) getConfig(ctx context.Context) *conf.Configuration {
	obj := ctx.Value(configKey)
	if obj == nil {
		return nil
	}
	config := obj.(*conf.Configuration)

	return config
}
