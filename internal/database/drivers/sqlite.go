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

func (s *SQLite) Connect(options map[string]string) error {
	var err error

	s.cfg = options
	dsn := "file:" + s.cfg["file"]
	s.db, err = sql.Open("sqlite3", dsn)

	return err
}

func (s *SQLite) Close() {

}

func (s *SQLite) GetConfigModel() map[string]string {
	return map[string]string{
		"file": "",
	}
}

func (s *SQLite) GetSessionStoreDriver(maxAge int, secrets ...[]byte) (sessions.Store, error) {
	return sqlitestore.NewSqliteStore(s.cfg["file"], "sessions", "/", maxAge, secrets...)
}
