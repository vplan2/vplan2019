package auth

import (
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/zekroTJA/vplan2019/internal/database"
)

type RequestAuthManager struct {
	tokenManager    *TokenManager
	sessionStore    sessions.Store
	disallowHandler func(w http.ResponseWriter, r *http.Request)
	errorHandler    func(w http.ResponseWriter, r *http.Request, err error)
}

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

func (ram *RequestAuthManager) Check(w http.ResponseWriter, r *http.Request) string {
	token := r.Header.Get("Authorization")

	if token != "" {
		ident, err := ram.tokenManager.Check(token)
		if err != nil {
			ram.errorHandler(w, r, err)
			return ""
		}

		if ident != "" {
			return ident
		}
	}

	session, err := ram.sessionStore.Get(r, "main")
	if err != nil {
		ram.errorHandler(w, r, err)
		return ""
	}

	if session.IsNew {
		ram.disallowHandler(w, r)
		return ""
	}

	ident, ok := session.Values["ident"].(string)

	if !ok || ident == "" {
		ram.disallowHandler(w, r)
		return ""
	}

	return ident
}
