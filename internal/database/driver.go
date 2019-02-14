// Package database contains general utilities and the database
// Driver interface for managing database connections
//   Authors: Ringo Hoffmann
package database

import (
	"errors"

	"github.com/gorilla/sessions"
)

// ErrConfig describes the type of error returned if the config
// interface could not be parsed to the database scheme specified
var ErrConfig = errors.New("failed parsing config for database")

// Driver is the general interface for database drivers
type Driver interface {
	// Connect to the database or open database file
	// with the passed options
	Connect(options map[string]string) error
	// Close database connection or file
	Close()

	// Get map which defines the key-value config
	// model structure
	GetConfigModel() map[string]string
	// Get session store driver
	GetSessionStoreDriver(secrets ...[]byte) (sessions.Store, error)
}
