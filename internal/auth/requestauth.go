package auth

import (
	"errors"
	"net/http"

	"github.com/zekroTJA/vplan2019/internal/logger"

	"github.com/gorilla/sessions"

	"github.com/zekroTJA/vplan2019/internal/database"
)

const (
	// MainSessionName describes the name of the cookie
	// of login sessions used for general authorization
	MainSessionName = "session_main"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// RequestAuthManager contains functionalities to check if
// a request is authorized by checking multiple methods
type RequestAuthManager struct {
	tokenManager    *TokenManager
	sessionStore    sessions.Store
	disallowHandler func(w http.ResponseWriter, r *http.Request)
	errorHandler    func(w http.ResponseWriter, r *http.Request, err error)
}

// NewRequestAuthManager creates a new instance of RequestAuthManager
//   db                     : database driver instance
//   tokenManager           : TokenManager instance which should be used
//   sessionStorageInstance : SessionStorage instance which should be used
//   disallowHandler        : handler which will be executed if a request was unauthorized
//   errorHandler           : handler which will be called if an unexpected error occures
func NewRequestAuthManager(db database.Driver, tokenManager *TokenManager, sessionStorageInstance sessions.Store,
	disallowHandler func(w http.ResponseWriter, r *http.Request), errorHandler func(w http.ResponseWriter, r *http.Request, err error)) *RequestAuthManager {

	if disallowHandler == nil {
		disallowHandler = func(w http.ResponseWriter, r *http.Request) {}
	}

	if errorHandler == nil {
		errorHandler = func(w http.ResponseWriter, r *http.Request, err error) {}
	}

	return &RequestAuthManager{
		tokenManager:    tokenManager,
		sessionStore:    sessionStorageInstance,
		disallowHandler: disallowHandler,
		errorHandler:    errorHandler,
	}
}

func (ram *RequestAuthManager) Authenticate(r *http.Request) (string, error) {
	//-------------------------------------
	// CHECKING AUTHORIZATION HEADER TOKEN

	token := r.Header.Get("Authorization")

	if token != "" {
		ident, err := ram.tokenManager.Check(token)
		if err != nil {
			return "", err
		}

		if ident != "" {
			return ident, nil
		}
	}

	//----------------------
	// CHECK SESSION COOKIE

	session, err := ram.sessionStore.Get(r, MainSessionName)
	if err != nil {
		return "", err
	}

	if session.IsNew {
		return "", ErrUnauthorized
	}

	ident, ok := session.Values["ident"].(string)

	if !ok || ident == "" {
		return "", ErrUnauthorized
	}

	return ident, nil
}

// Check returns the ident of the user, if the request was authorized.
// First, this function will look for the 'Authorization' header containing a
// token value, which will be checked against the database entries.
// If there was no matching token found, the function will check for a valid
// session cookie to authorize.
// If both methods are failing, an empty string will be returned without an error
// and the disallowHandler will be called.
// If an unexpected error occures at any point, the function will also return an
// empty string and will call the errorHandler passing the error object.
func (ram *RequestAuthManager) Check(w http.ResponseWriter, r *http.Request) string {
	ident, err := ram.Authenticate(r)

	if err == ErrUnauthorized {
		logger.Debug("session login error: %s", err.Error())
		ram.disallowHandler(w, r)
		return ""
	}

	if err != nil {
		ram.errorHandler(w, r, err)
		return ""
	}

	return ident
}
