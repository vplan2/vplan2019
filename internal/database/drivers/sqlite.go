// Package drivers contains database driver structs for
// accessing various database types
//   Authors: Ringo Hoffmann
package drivers

import (
	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/michaeljs1990/sqlitestore"
)

// SQLite contains database functions
// for SQLite database
type SQLite struct {
	cfg map[string]string
	db  *sql.DB
}

// Connect opens a sqlite3 database file or creates
// it if it does not exist depending on the passed options
func (s *SQLite) Connect(options map[string]string) error {
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

func (s *SQLite) GetConfigModel() map[string]string {
	return map[string]string{
		"file": "",
	}
}

func (s *SQLite) GetSessionStoreDriver(maxAge int, secrets ...[]byte) (sessions.Store, error) {
	return sqlitestore.NewSqliteStore(s.cfg["file"], "sessions", "/", maxAge, secrets...)
}
