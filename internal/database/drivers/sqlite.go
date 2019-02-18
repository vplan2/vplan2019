// Package drivers (database/drivers) contains database
// driver structs for accessing various database types
//   Authors: Ringo Hoffmann
package drivers

import (
	"database/sql"
	"time"

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

// Setup creates tables if they do not exist yet
func (s *SQLite) Setup() error {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS `api.token` (" +
		"`ident` text NOT NULL DEFAULT ''," +
		"`token` text NOT NULL DEFAULT ''," +
		"`expire` int NOT NULL DEFAULT 0 );")
	if err != nil {
		return err
	}

	return nil
}

// GetUserAPIToken gets the users API token with the time, the token expires,
// if existent. Else, the returned string will be empty. If the query failes,
// an error will be returned
func (s *SQLite) GetUserAPIToken(ident string) (string, time.Time, error) {
	var token string
	var expire time.Time

	row := s.db.QueryRow("SELECT token, expire FROM api.token WHERE ident = ?", ident)
	err := row.Scan(&token, &expire)

	return token, expire, err
}

// SetUserAPIToken sets the API token an the expire time of it for a user
func (s *SQLite) SetUserAPIToken(ident, token string, expire time.Time) error {
	res, err := s.db.Exec("UPDATE api.token SET token = ?, expire = ? WHERE ident = ?",
		token, expire, ident)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); err != nil {
		return err
	} else if ar < 1 {
		_, err = s.db.Exec("INSERT INTO api.token (ident, token, expire) VALUES (?, ?, ?)",
			ident, token, expire)
	}
	return err
}

// DeleteUserAPIToken removes a users token from the database
func (s *SQLite) DeleteUserAPIToken(ident string) error {
	_, err := s.db.Exec("DELETE FROM api.token WHERE ident = ?", ident)
	return err
}

// GetConfigModel returns a map with preset config
// keys and values
func (s *SQLite) GetConfigModel() map[string]string {
	return map[string]string{
		"file": "main.db.sqlite3",
	}
}

// GetSessionStoreDriver returns a new instance of the session
// store driver, which should be used for saving encrypted session data
func (s *SQLite) GetSessionStoreDriver(maxAge int, secrets ...[]byte) (sessions.Store, error) {
	return sqlitestore.NewSqliteStore(s.cfg["file"], "api.sessions", "/", maxAge, secrets...)
}
