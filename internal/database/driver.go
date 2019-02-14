// Package database contains general utilities and the database
// Driver interface for managing database connections
//   Authors: Ringo Hoffmann
package database

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/zekroTJA/vplan2019/internal/config"
)

// ErrConfig describes the type of error returned if the config
// interface could not be parsed to the database scheme specified
var ErrConfig = errors.New("failed parsing config for database")

// Driver is the general interface for database drivers
type Driver interface {
	// Connect to the database or open database file
	// with the passed options
	Connect(options config.Model) error
	// Close database connection or file
	Close()

	// Get map which defines the key-value config
	// model structure
	GetConfigModel() config.Model
	// Get session store driver
	GetSessionStoreDriver(secrets ...[]byte) (sessions.Store, error)
}
