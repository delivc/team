package cmd

import (
	"database/sql"
	"net/url"

	"github.com/delivc/team/conf"
	"github.com/delivc/team/models"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrateCmd = cobra.Command{
	Use:  "migrate",
	Long: "Migrate database strucutures. This will create new tables and add missing columns and indexes.",
	Run:  migrate,
}

func migrate(cmd *cobra.Command, args []string) {
	writer := logrus.StandardLogger().Writer()
	globalConfig, err := conf.LoadGlobal(configFile)
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %+v", err)
	}
	if globalConfig.DB.Driver == "" && globalConfig.DB.URL != "" {
		u, err := url.Parse(globalConfig.DB.URL)
		if err != nil {
			logrus.Fatalf("%+v", errors.Wrap(err, "parsing db connection url"))
		}
		globalConfig.DB.Driver = u.Scheme
	}
	pop.Debug = true

	logrus.Info((globalConfig.DB.URL))

	deets := &pop.ConnectionDetails{
		Dialect: globalConfig.DB.Driver,
		URL:     globalConfig.DB.URL,
	}

	if globalConfig.DB.Namespace != "" {
		deets.Options = map[string]string{
			"Namespace": globalConfig.DB.Namespace + "_",
		}
	}

	db, err := pop.NewConnection(deets)
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "opening db connection"))
	}
	defer db.Close()

	if err := db.Open(); err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "checking database connection"))
	}

	logrus.Infof("Reading migrations from %s", globalConfig.DB.MigrationsPath)
	mig, err := pop.NewFileMigrator(globalConfig.DB.MigrationsPath, db)
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "creating db migrator"))
	}
	logrus.Infof("before status")
	err = mig.Status(writer)
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "migration status"))
	}
	// turn off schema dump
	mig.SchemaPath = ""

	err = mig.Up()
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "running db migrations"))
	}

	logrus.Infof("after status")
	err = mig.Status(writer)
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "migration status"))
	}

	err = db.Transaction(func(tx *pop.Connection) error {
		return nil
	})
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "seeding status"))
	}

	permissions := []string{
		"spaces-create",
		"spaces-edit",
		"spaces-delete",
		"spaces-read-apikeys",
		"spaces-create-apikeys",
		"spaces-destroy-apikeys",
		"spaces-create-models",
		"spaces-edit-models",
		"spaces-destroy-models",
		"spaces-create-content",
		"spaces-edit-content",
		"spaces-destroy-content",
		"spaces-create-assets",
		"spaces-destroy-assets",
		"account-edit",
		"account-destroy",
		"account-users-invite",
		"account-users-remove",
	}

	err = db.Transaction(func(tx *pop.Connection) error {
		for _, permission := range permissions {
			p, err := models.NewPermission(permission)
			if err != nil {
				logrus.Fatalf("%+v", errors.Wrap(err, "creating permission"))
				return err
			}
			obj := &models.Permission{}
			if err := tx.Q().Where("name = ?", p.Name).First(obj); err != nil {
				if errors.Cause(err) != sql.ErrNoRows {
					return err
				}
			}
			if obj.Name == "" {
				tx.Create(p)
			}
		}

		return nil
	})

	// Migrate default Permissions
	// can:
	// manage billing
}
