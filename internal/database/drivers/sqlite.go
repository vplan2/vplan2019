// Package drivers contains database driver structs for
// accessing various database types
//   Authors: Ringo Hoffmann
package drivers

import (
	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/michaeljs1990/sqlitestore"
	"github.com/zekroTJA/vplan2019/internal/config"
)

// SQLite contains database functions
// for SQLite database
type SQLite struct {
	cfg config.Model
	db  *sql.DB
}

// Connect opens a sqlite3 database file or creates
// it if it does not exist depending on the passed options
func (s *SQLite) Connect(options config.Model) error {
	var err error

	s.cfg = options
	dsn := "file:" + s.cfg["file"]
	s.db, err = sql.Open("sqlite3", dsn)

	return err
}

// Close closes the sqlite3 database file
func (s *SQLite) Close() {
	s.db.Close()
}

// GetConfigModel returns a map with preset config
// keys and values
func (s *SQLite) GetConfigModel() config.Model {
	return map[string]string{
		"file": "main.db.sqlite3",
	}
}

// GetSessionStoreDriver returns a new instance of the session
// store driver, which should be used for saving encrypted session data
func (s *SQLite) GetSessionStoreDriver(maxAge int, secrets ...[]byte) (sessions.Store, error) {
	return sqlitestore.NewSqliteStore(s.cfg["file"], "sessions", "/", maxAge, secrets...)
}
