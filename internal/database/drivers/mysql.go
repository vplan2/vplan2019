// Package drivers (database/drivers) contains database
// driver structs for accessing various database types
//   Authors: Ringo Hoffmann
package drivers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gorilla/sessions"
	"github.com/zekroTJA/mysqlstore"
)

// MySql contains database functions
// for MySql database
type MySql struct {
	cfg map[string]string
	dsn string
	db  *sql.DB
}

// Connect opens a MySql3 database file or creates
// it if it does not exist depending on the passed options
func (s *MySql) Connect(options map[string]string) error {
	var err error

	s.cfg = options
	s.dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		options["user"], options["password"], options["host"], options["database"])

	s.db, err = sql.Open("mysql", s.dsn)

	return err
}

// Close closes the MySql3 database file
func (s *MySql) Close() {
	s.db.Close()
}

// Setup creates tables if they do not exist yet
func (s *MySql) Setup() error {
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS `apitoken` (" +
		"`ident` text NOT NULL DEFAULT ''," +
		"`token` text NOT NULL DEFAULT ''," +
		"`expire` text NOT NULL DEFAULT '' );")
	if err != nil {
		return err
	}

	return nil
}

// GetAPIToken returns the matching indent and expire time to a found token.
// If the token could not be matched, this returns an empty string without
// and error. Errors are only returned if the database request failes.
func (s *MySql) GetAPIToken(token string) (string, time.Time, error) {
	var ident string
	var expire string

	row := s.db.QueryRow("SELECT ident, expire FROM apitoken WHERE token = ?", token)
	err := row.Scan(&ident, &expire)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return "", time.Time{}, err
	}

	if ident == "" {
		return ident, time.Time{}, nil
	}

	tExpire, err := time.Parse("2006-01-02 15:04:05.999999-07:00", expire)

	return ident, tExpire, err
}

// GetUserAPIToken gets the users API token with the time, the token expires,
// if existent. Else, the returned string will be empty. If the query failes,
// an error will be returned
func (s *MySql) GetUserAPIToken(ident string) (string, time.Time, error) {
	var token string
	var expire string

	row := s.db.QueryRow("SELECT token, expire FROM apitoken WHERE ident = ?", ident)
	err := row.Scan(&token, &expire)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return "", time.Time{}, err
	}

	tExpire, err := time.Parse("2006-01-02 15:04:05.999999-07:00", expire)

	return token, tExpire, err
}

// SetUserAPIToken sets the API token an the expire time of it for a user
func (s *MySql) SetUserAPIToken(ident, token string, expire time.Time) error {
	res, err := s.db.Exec("UPDATE apitoken SET token = ?, expire = ? WHERE ident = ?",
		token, expire, ident)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); err != nil {
		return err
	} else if ar < 1 {
		_, err = s.db.Exec("INSERT INTO apitoken (ident, token, expire) VALUES (?, ?, ?)",
			ident, token, expire)
	}
	return err
}

// DeleteUserAPIToken removes a users token from the database
func (s *MySql) DeleteUserAPIToken(ident string) error {
	_, err := s.db.Exec("DELETE FROM apitoken WHERE ident = ?", ident)
	return err
}

// GetConfigModel returns a map with preset config
// keys and values
func (s *MySql) GetConfigModel() map[string]string {
	return map[string]string{
		"host":     "localhost",
		"user":     "vplan2",
		"password": "",
		"database": "vplan2",
	}
}

// GetSessionStoreDriver returns a new instance of the session
// store driver, which should be used for saving encrypted session data
func (s *MySql) GetSessionStoreDriver(maxAge int, secrets ...[]byte) (sessions.Store, error) {
	return mysqlstore.NewMySQLStoreFromConnection(s.db, "apisessions", "/", maxAge, secrets...)
}
