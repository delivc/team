package test

import (
	"github.com/delivc/team/conf"
	"github.com/delivc/team/storage"
)

// SetupDBConnection setups a new Connection to the Database in a test env
func SetupDBConnection(globalConfig *conf.GlobalConfiguration) (*storage.Connection, error) {
	return storage.Dial(globalConfig)
}
