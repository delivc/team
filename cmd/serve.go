package cmd

import (
	"context"
	"fmt"

	"github.com/delivc/team/api"
	"github.com/delivc/team/conf"
	"github.com/delivc/team/storage"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

func serve(globalConfig *conf.GlobalConfiguration, config *conf.Configuration) {
	db, err := storage.Dial(globalConfig)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	ctx := api.WithInstanceConfig(context.Background(), config, uuid.Nil)
	api := api.New(ctx, globalConfig, db, Version)

	l := fmt.Sprintf("%v:%v", globalConfig.API.Host, globalConfig.API.Port)
	logrus.Infof("Delivc Team API (%s) started on: %s", Version, l)
	api.ListenAndServe(l)
}
