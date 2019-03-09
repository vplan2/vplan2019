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
	"github.com/zekroTJA/vplan2019/internal/database"
	"github.com/zekroTJA/vplan2019/pkg/multierror"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

// MySQL contains database functions
// for MySQL database
type MySQL struct {
	cfg   map[string]string
	dsn   string
	db    *sql.DB
	stmts *prepStatements
}

type prepStatements struct {
	selectAPITokenByToken *sql.Stmt
	selectAPITokenByIdent *sql.Stmt
	updateAPIToken        *sql.Stmt
	insertAPIToken        *sql.Stmt
	deleteAPIToken        *sql.Stmt

	insertLogin *sql.Stmt
	getLogins   *sql.Stmt

	selectVPlans              *sql.Stmt
	selectVPlanEntries        *sql.Stmt
	selectVPlanEntriesByClass *sql.Stmt

	getUserSettings    *sql.Stmt
	setUserSettings    *sql.Stmt
	insertUserSettings *sql.Stmt
}

////////////////////////
// SETUP AND TEARDOWN //
////////////////////////

// Connect opens a MySql3 database file or creates
// it if it does not exist depending on the passed options
func (s *MySQL) Connect(options map[string]string) error {
	var err error

	s.cfg = options
	s.dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		options["user"], options["password"], options["host"], options["database"])

	s.db, err = sql.Open("mysql", s.dsn)

	return err
}

func (s *MySQL) prepareStatement(multiError *multierror.MultiError, query string) *sql.Stmt {
	stmt, err := s.db.Prepare(query)
	multiError.Append(err)
	return stmt
}

func (s *MySQL) setupPrepStatements() error {
	s.stmts = new(prepStatements)
	m := multierror.NewMultiError(nil)

	s.stmts.selectAPITokenByToken = s.prepareStatement(m, "SELECT ident, expire FROM apitoken WHERE token = ?")
	s.stmts.selectAPITokenByIdent = s.prepareStatement(m, "SELECT token, expire FROM apitoken WHERE ident = ?")
	s.stmts.updateAPIToken = s.prepareStatement(m, "UPDATE apitoken SET token = ?, expire = ? WHERE ident = ?")
	s.stmts.insertAPIToken = s.prepareStatement(m, "INSERT INTO apitoken (ident, token, expire) VALUES (?, ?, ?)")
	s.stmts.deleteAPIToken = s.prepareStatement(m, "DELETE FROM apitoken WHERE ident = ?")

	s.stmts.insertLogin = s.prepareStatement(m,
		"INSERT INTO logins (ident, type, useragent, ipaddress) VALUES (?, ?, ?, ?)")
	s.stmts.getLogins = s.prepareStatement(m,
		"SELECT ident, timestamp, type, useragent, ipaddress FROM logins WHERE "+
			"ident = ? AND timestamp >= ?")

	s.stmts.selectVPlans = s.prepareStatement(m,
		"SELECT id, date_edit, date_for, block, header, footer FROM vplan WHERE "+
			"date_for >= ? AND deleted = 0 "+
			"ORDER BY date_for ASC")
	s.stmts.selectVPlanEntries = s.prepareStatement(m,
		"SELECT id, vplan_id, class, time, messures, responsible FROM vplan_details WHERE vplan_id = ? AND deleted = 0")
	s.stmts.selectVPlanEntriesByClass = s.prepareStatement(m,
		"SELECT id, vplan_id, class, time, messures, responsible FROM vplan_details WHERE vplan_id = ? AND class = ? AND deleted = 0")

	s.stmts.getUserSettings = s.prepareStatement(m,
		"SELECT ident, class, theme, edited FROM usersettings WHERE ident = ?")
	s.stmts.setUserSettings = s.prepareStatement(m,
		"UPDATE usersettings SET class = ?, theme = ? WHERE ident = ?")
	s.stmts.insertUserSettings = s.prepareStatement(m,
		"INSERT INTO usersettings (ident, class, theme) VALUES (?, ?, ?)")

	return m.Concat()
}

// Close closes the MySql3 database file
func (s *MySQL) Close() {
	s.db.Close()
}

// Setup creates tables if they do not exist yet
func (s *MySQL) Setup() error {
	m := multierror.NewMultiError(nil)

	// TABLE `vplan`
	_, err := s.db.Exec("CREATE TABLE IF NOT EXISTS `vplan` (" +
		"`id` int(11) NOT NULL," +
		"`date_edit` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
		"`date_for` datetime NOT NULL," +
		"`block` char(1) NOT NULL," +
		"`header` text NOT NULL," +
		"`footer` text NOT NULL," +
		"`published` tinyint(1) NOT NULL," +
		"`deleted` int(1) NOT NULL DEFAULT '0' );")
	m.Append(err)

	// TABLE `vplan_details`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `vplan_details` (" +
		"`id` int(11) NOT NULL," +
		"`vplan_id` int(11) NOT NULL," +
		"`class` varchar(45) NOT NULL," +
		"`time` varchar(45) NOT NULL," +
		"`messures` varchar(255) NOT NULL," +
		"`responsible` varchar(255) NOT NULL," +
		"`reason` int(1) NOT NULL DEFAULT '1'," +
		"`geteilt` int(1) NOT NULL," +
		"`notiz` varchar(45) NOT NULL," +
		"`deleted` int(1) NOT NULL DEFAULT '0'," +
		"`selected` int(1) NOT NULL DEFAULT '0' );")
	m.Append(err)

	// TABLE `apitoken`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `apitoken` (" +
		"`id` int PRIMARY KEY AUTO_INCREMENT," +
		"`ident` text NOT NULL," +
		"`token` text NOT NULL," +
		"`expire` timestamp NOT NULL );")
	m.Append(err)

	// TABLE `usersettings`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `usersettings` (" +
		"`id` int PRIMARY KEY AUTO_INCREMENT," +
		"`ident` text NOT NULL," +
		"`class` text NOT NULL," +
		"`theme` text NOT NULL," +
		"`edited` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP);")
	m.Append(err)

	// TABLE `logins`
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS `logins` (" +
		"`id` int PRIMARY KEY AUTO_INCREMENT," +
		"`ident` text NOT NULL," +
		"`timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"`type` int NOT NULL DEFAULT 0," +
		"`useragent` text NOT NULL," +
		"`ipaddress` text NOT NULL );")
	m.Append(err)

	if err := m.Concat(); err != nil {
		return err
	}

	return s.setupPrepStatements()
}

//////////////
// SETTINGS //
//////////////

// GetConfigModel returns a map with preset config
// keys and values
func (s *MySQL) GetConfigModel() map[string]string {
	return map[string]string{
		"host":     "localhost",
		"user":     "vplan2",
		"password": "",
		"database": "vplan2",
	}
}

// GetSessionStoreDriver returns a new instance of the session
// store driver, which should be used for saving encrypted session data
func (s *MySQL) GetSessionStoreDriver(maxAge int, secrets ...[]byte) (sessions.Store, error) {
	return mysqlstore.NewMySQLStoreFromConnection(s.db, "apisessions", "/", maxAge, secrets...)
}

////////////////
// API TOKENS //
////////////////

// GetAPIToken returns the matching indent and expire time to a found token.
// If the token could not be matched, this returns an empty string without
// and error. Errors are only returned if the database request failes.
func (s *MySQL) GetAPIToken(token string) (string, time.Time, error) {
	var ident string
	var expire database.Timestamp

	row := s.stmts.selectAPITokenByToken.QueryRow(token)
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

	tExpire, err := expire.ToTime(timeFormat)

	return ident, tExpire, err
}

// GetUserAPIToken gets the users API token with the time, the token expires,
// if existent. Else, the returned string will be empty. If the query failes,
// an error will be returned
func (s *MySQL) GetUserAPIToken(ident string) (string, time.Time, error) {
	var token string
	var expire database.Timestamp

	row := s.stmts.selectAPITokenByIdent.QueryRow(ident)
	err := row.Scan(&token, &expire)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return "", time.Time{}, err
	}

	tExpire, err := expire.ToTime(timeFormat)

	return token, tExpire, err
}

// SetUserAPIToken sets the API token an the expire time of it for a user
func (s *MySQL) SetUserAPIToken(ident, token string, expire time.Time) error {
	res, err := s.stmts.updateAPIToken.Exec(token, expire, ident)
	if err != nil {
		return err
	}
	if ar, err := res.RowsAffected(); err != nil {
		return err
	} else if ar < 1 {
		_, err = s.stmts.insertAPIToken.Exec(ident, token, expire)
	}
	return err
}

// DeleteUserAPIToken removes a users token from the database
func (s *MySQL) DeleteUserAPIToken(ident string) error {
	_, err := s.stmts.deleteAPIToken.Exec(ident)
	return err
}

///////////////////
// LOGIN LOGGING //
///////////////////

// InsertLogin inserts logiin infrmation to login table
func (s *MySQL) InsertLogin(loginType database.LoginType, ident, useragent, ipaddress string) error {
	_, err := s.stmts.insertLogin.Exec(ident, loginType, useragent, ipaddress)
	return err
}

// GetLogins returns a list of entries from the login log filtered by Ident of a user and
// after the time passed
func (s *MySQL) GetLogins(ident string, afterTimestamp time.Time) ([]*database.Login, error) {
	rows, err := s.stmts.getLogins.Query(ident, afterTimestamp)
	if err != nil {
		return nil, err
	}

	mErr := multierror.NewMultiError(nil)
	logins := make([]*database.Login, 0)
	for rows.Next() {
		var timeStamp database.Timestamp
		login := new(database.Login)

		err = rows.Scan(&login.Ident, &timeStamp, &login.Type, &login.Useragent, &login.IPAddress)
		mErr.Append(err)
		if err != nil {
			continue
		}

		login.Timestamp, err = timeStamp.ToTime(timeFormat)
		mErr.Append(err)
		if err != nil {
			continue
		}

		logins = append(logins, login)
	}

	return logins, mErr.Concat()
}

////////////
// VPLANS //
////////////

// GetVPlans collects VPlans wich for-dates are after the passed timestamp.
// Also, a class can be specified for filtering the VPlanEntries.
func (s *MySQL) GetVPlans(class string, timestamp time.Time) ([]*database.VPlan, error) {
	rows, err := s.stmts.selectVPlans.Query(timestamp)
	if err != nil {
		return nil, err
	}

	vplans := make([]*database.VPlan, 0)
	mErr := multierror.NewMultiError(nil)

	var dateEdit, dateFor database.Timestamp
	for rows.Next() {
		vplan := new(database.VPlan)
		vplan.Entries = make([]*database.VPlanEntry, 0)
		err = rows.Scan(&vplan.ID, &dateEdit, &dateFor, &vplan.Block, &vplan.Header, &vplan.Footer)
		mErr.Append(err)
		if err == nil {
			vplan.DateEdit, _ = dateEdit.ToTime(timeFormat)
			vplan.DateFor, _ = dateFor.ToTime(timeFormat)
			vplans = append(vplans, vplan)
		}
	}

	for _, v := range vplans {
		if class != "" {
			rows, err = s.stmts.selectVPlanEntriesByClass.Query(v.ID, class)
		} else {
			rows, err = s.stmts.selectVPlanEntries.Query(v.ID)
		}
		mErr.Append(err)
		if err != nil {
			continue
		}

		for rows.Next() {
			entry := new(database.VPlanEntry)
			err = rows.Scan(&entry.ID, &entry.VPlanID, &entry.Class, &entry.Time, &entry.Messures, &entry.Resposible)
			mErr.Append(err)
			if err == nil {
				v.Entries = append(v.Entries, entry)
			}
		}
	}

	return vplans, mErr.Concat()
}

//////////////////
// USER SETTNGS //
//////////////////

// GetUserSettings returns the personal settings of a user. If there is no setting,
// an empty struct will be returned with 'false' as second return value. If there
// was a setting found, the second return value will be 'true'.
func (s *MySQL) GetUserSettings(ident string) (*database.UserSetting, bool, error) {
	settings := new(database.UserSetting)
	var edited database.Timestamp
	var err error

	err = s.stmts.getUserSettings.QueryRow(ident).Scan(
		&settings.Ident, &settings.Class, &settings.Theme, &edited)
	if err == sql.ErrNoRows {
		return settings, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	settings.Edited, err = edited.ToTime(timeFormat)
	if err != nil {
		return nil, false, err
	}

	return settings, true, nil
}

// SetUserSetting sets or inserts the personal user settings of a user.
// Only changed values in the settings object will be updated in the database.
// If a value should be reset (default init value, e.g. '' for strings), set the
// settings value to "reset" or -1.
func (s *MySQL) SetUserSetting(ident string, updateSetting *database.UserSetting) error {
	settings, exists, err := s.GetUserSettings(ident)
	if err != nil {
		return err
	}

	if !exists {
		_, err = s.stmts.insertUserSettings.Exec(ident, updateSetting.Class, updateSetting.Theme)
		return err
	}

	if updateSetting.Class == "reset" {
		settings.Theme = ""
	} else if updateSetting.Class != "" {
		settings.Class = updateSetting.Class
	}

	if updateSetting.Theme == "reset" {
		settings.Theme = ""
	} else if updateSetting.Theme != "" {
		settings.Theme = updateSetting.Theme
	}

	_, err = s.stmts.setUserSettings.Exec(settings.Class, settings.Theme, ident)

	return err
}
