// Package database contains general utilities and the database
// Driver interface for managing database connections
//   Authors: Ringo Hoffmann
package database

import (
	"errors"
	"time"

	"github.com/gorilla/sessions"
)

const (
	// LoginTypeWebInterface describes a web interface login
	LoginTypeWebInterface LoginType = iota
	// LoginTypeToken describes an API token generation
	LoginTypeToken
)

// ErrConfig describes the type of error returned if the config
// interface could not be parsed to the database scheme specified
var ErrConfig = errors.New("failed parsing config for database")

// Timestamp is the standard timestamp type
// returned from the MySQL database
type Timestamp []uint8

// LoginType is the type of login
type LoginType int8

// VPlan contains the information of
// a VPlan database structure
type VPlan struct {
	ID       int           `json:"id"`
	DateEdit time.Time     `json:"date_edit"`
	DateFor  time.Time     `json:"date_for"`
	Block    string        `json:"block"`
	Header   string        `json:"header"`
	Footer   string        `json:"footer"`
	Entries  []*VPlanEntry `json:"entries"`
}

// VPlanEntry contains the information of
// a VPlan entry database structure
type VPlanEntry struct {
	ID         int    `json:"id"`
	VPlanID    int    `json:"vplan_id"`
	Class      string `json:"class"`
	Time       string `json:"time"`
	Messures   string `json:"messures"`
	Resposible string `json:"responsible"`
}

// TickerEntry contains the informations of
// a Newsticker entry database structure
type TickerEntry struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	Headline string    `json:"headline"`
	Short    string    `json:"short"`
	Story    string    `json:"story"`
}

// Login contains the data of a login
// log database structure
type Login struct {
	Ident     string    `json:"ident"`
	Timestamp time.Time `json:"timestamp"`
	Type      LoginType `json:"type"`
	Useragent string    `json:"useragent"`
	IPAddress string    `json:"ipaddress"`
}

// UserSetting contains the data of a
// user setting database structure
type UserSetting struct {
	Ident  string    `json:"ident"`
	Class  string    `json:"class"`
	Theme  string    `json:"theme"`
	Edited time.Time `json:"edited"`
}

// Driver is the general interface for database drivers
type Driver interface {

	////////////////////////
	// SETUP AND TEARDOWN //
	////////////////////////

	// Connect to the database or open database file
	// with the passed options
	Connect(options map[string]string) error
	// Close database connection or file
	Close()
	// Setup deos stuff like chekcing for database schemes,
	// creating them if necessary, validating data ...
	Setup() error

	//////////////
	// SETTINGS //
	//////////////

	// Get map which defines the key-value config
	// model structure
	GetConfigModel() map[string]string
	// Get session store driver
	GetSessionStoreDriver(mayAge int, secrets ...[]byte) (sessions.Store, error)

	////////////////
	// API TOKENS //
	////////////////

	// GetAPIToken returns the passing ident to the fount token.
	// If no matching token was found, an empty stirng should be
	// returned with an nil error.
	// Only return an error when the database access failed.
	GetAPIToken(token string) (indent string, expire time.Time, err error)
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

	///////////////////
	// LOGIN LOGGING //
	///////////////////

	// InsertLogin inserts logiin infrmation to login table
	InsertLogin(loginType LoginType, ident, useragent, ipaddress string) error

	GetLogins(ident string, afterTimestamp time.Time) ([]*Login, error)

	////////////
	// VPLANS //
	////////////

	// GetVPlans collects VPlans and VPlanEntries filtered by the
	// passed Timestamp (all after time) and Class name.
	GetVPlans(class string, timestamp time.Time) ([]*VPlan, error)

	////////////////
	// NEWSTICKER //
	////////////////

	// GetNewsTicker collects news ticker entries after the
	// specified timestamp
	GetNewsTicker(timestamp time.Time) ([]*TickerEntry, error)

	//////////////////
	// USER SETTNGS //
	//////////////////

	// GetUserSettings returns the personal settings of a user. If there is no setting,
	// an empty struct will be returned with 'false' as second return value. If there
	// was a setting found, the second return value will be 'true'.
	GetUserSettings(ident string) (*UserSetting, bool, error)
	// SetUserSetting sets or inserts the personal user settings of a user.
	// Only changed values in the settings object will be updated in the database.
	// If a value should be reset (default init value, e.g. '' for strings), set the
	// settings value to "reset" or -1.
	SetUserSetting(ident string, updateSetting *UserSetting) error
}

// -------------------------------------------

// ToTime parses the timestamp to a time object
func (t Timestamp) ToTime(format string) (time.Time, error) {
	return time.Parse(format, string(t))
}
