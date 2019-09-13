package auth

import (
	"errors"
	"net/http"

	"github.com/vplan2/vplan2019/internal/logger"

	"github.com/gorilla/sessions"

	"github.com/vplan2/vplan2019/internal/database"
)

const (
	// MainSessionName describes the name of the cookie
	// of login sessions used for general authorization
	MainSessionName = "session_main"
)

var (
	// ErrUnauthorized defines an error when an authorization
	// failed because of wrong credentials
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

// Authorize checks the request for an 'Authorization' header passing
// an API token to authenticate against the server. If there was no token
// passed or if the token was invalid, the request will be checked for a
// valid session cookie containing a token which will be used to authenticate
// against the server. If the authorization succeeds, the Ident will be
// returned. Else, an error will be passed with an empty Indent string.
func (ram *RequestAuthManager) Authorize(r *http.Request) (string, error) {
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

// Check authorizes the request by executing the 'Authorize' function returning
// the users Ident. If the authorization failes, an empty string will be returned
// and an error message will be respondet to the client.
func (ram *RequestAuthManager) Check(w http.ResponseWriter, r *http.Request) string {
	ident, err := ram.Authorize(r)

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
