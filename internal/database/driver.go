// Package database contains general utilities and the database
// Driver interface for managing database connections
//   Authors: Ringo Hoffmann
package database

import (
	"errors"
	"time"

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
	// Setup deos stuff like chekcing for database schemes,
	// creating them if necessary, validating data ...
	Setup() error

	// Get map which defines the key-value config
	// model structure
	GetConfigModel() map[string]string
	// Get session store driver
	GetSessionStoreDriver(secrets ...[]byte) (sessions.Store, error)

	// GetUserAPIToken returns the user API token, if existent
	// If there is no token existent, this should only return an
	// empty string and only an error if the db reuqest failes
	GetUserAPIToken(ident string) (token string, expire time.Time, err error)
	// SetUserAPIToken sets a new API token for the specified user
	SetUserAPIToken(ident, token string, expire time.Time) error
	// DeleteUserAPIToken removes a users API token from database.
	// This should only return an error if the db request failes.
	// If the user has no token to delete, you should not return
	// an error.
	DeleteUserAPIToken(ident string) error
}
