// Package webserver contains general utilities
// for exposing and handling a web server.
//   Authors: Ringo Hoffmann
package webserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/zekroTJA/vplan2019/internal/database"

	"github.com/zekroTJA/vplan2019/internal/auth"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	tokenLifetime = 30 * 24 * time.Hour
)

var (
	// ErrUnauthorized is returned on unauthorized access
	ErrUnauthorized = errors.New("unauthorized")
)

// Config contains the configuration
// for the web server
type Config struct {
	Addr     string          `json:"addr"`
	Sessions *ConfigSessions `json:"sessions"`
	TLS      *ConfigTLS      `json:"tls"`
	TVUser   string          `json:"tvuser"`

	StaticFiles string `json:",omitempty"`
}

// ConfigTLS contains the cert file path
// and the key file path for TLS configuration
type ConfigTLS struct {
	UseSSL   bool   `json:"usessl"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// ConfigSessions contains the configuration
// for the session management
type ConfigSessions struct {
	DefaultMaxAge    int    `json:"defaultMaxAge"`
	RememberMaxAge   int    `json:"rememberMaxAge"`
	EncryptionSecret string `json:"encryptionSecret"`
}

// Server contains the instance of the
// http server and the mux router
type Server struct {
	config       *Config
	tokenManager *auth.TokenManager
	reqAuth      *auth.RequestAuthManager
	db           database.Driver
	server       *http.Server
	router       *mux.Router
	store        sessions.Store
	authProvider auth.Provider
	limiter      *RateLimiter
}

var (
	// Errors
	errServerInstanceNil = errors.New("web server instance must be initialized")
	errConfigNil         = errors.New("web server config is nil")
)

// StartBlocking starts the web server and block the current thread
//   server : the initialized instance of an empty Server struct
//   config : instance of Config; if this is nil, defaultConfig will be used
func StartBlocking(server *Server, config *Config, db database.Driver, sessionStorage sessions.Store, authProvider auth.Provider) error {
	if server == nil {
		return errServerInstanceNil
	}

	if config == nil {
		return errConfigNil
	}

	server.db = db
	server.config = config
	server.router = mux.NewRouter()
	server.store = sessionStorage
	server.authProvider = authProvider
	server.tokenManager = auth.NewTokenManager(db, tokenLifetime)
	server.reqAuth = auth.NewRequestAuthManager(db, server.tokenManager, sessionStorage,
		server.handlerAPIUnauthorizedError, server.handlerAPIInternalError)

	server.limiter = NewRateLimiter(&LimiterOpts{10, 10}, server.handlerAPIRateLimitError)

	server.server = &http.Server{
		Addr:    config.Addr,
		Handler: server.router,
	}

	server.initializeHnalders()

	if config.TLS == nil || !config.TLS.UseSSL {
		return server.server.ListenAndServe()
	}
	return server.server.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile)
}

// addHandler is a help function for adding handlers to the router
// and for registering the entry for the rate limiter in one function
//   path         : url path of the route
//   ident        : the endpoint ident string for the rate limiter
//   handler      : the handler function of the request
//   limiterRate  : ammount of tokens regenerated per seconds
//   limiterBurst : initial and total size of the token bucket
//   ...methods   : allowed HTTP methods
func (s *Server) addHandler(path string, ident string, handler func(w http.ResponseWriter, r *http.Request), limiterRate float64, limiterBurst int, methods ...string) {
	s.router.HandleFunc(path, handler).Methods(methods...)
	s.limiter.Register(ident, limiterRate, limiterBurst)
}

// initializeHandlers contains all setup functions for router
// endpoints and their handlers
func (s *Server) initializeHnalders() {

	// NOTICE:
	// Each route over another overwrites the access area of
	// of the path of the route below.

	// ---------------------------
	// API

	// POST /api/authenticate/:USERNAME
	s.addHandler("/api/authenticate/{username}", "authenticate",
		s.handlerAPIAuthenticate, 0.2, 3, "POST")

	// POST /api/logout
	s.addHandler("/api/logout", "logout", s.handlerAPILogout, 1, 3, "POST")

	// GET /api/logins
	s.addHandler("/api/logins", "getLogins", s.handlerAPIGetLogins, 1, 3, "GET")

	// GET /api/vplan
	s.addHandler("/api/vplan", "getVPlan", s.handlerAPIGetVPlan, 0.2, 3, "GET")

	// GET /api/newsticker
	s.addHandler("/api/newsticker", "getNewsTicker", s.handlerAPIGetNewsTicker, 0.2, 3, "GET")

	// GET /api/settings
	s.addHandler("/api/settings", "getUserSettings", s.handleAPIGetUserSettings, 1, 3, "GET")

	// POST /api/settings
	s.addHandler("/api/settings", "setUserSettings", s.handleAPISetUserSettings, 0.5, 5, "POST")

	// POST /api/test
	s.addHandler("/api/test", "test", s.handlerAPITest, 1, 1, "POST")

	// ---------------------------
	// FRONTEND

	s.router.HandleFunc(`/login`, s.handlerFELogin)

	// GET /:FILENAME
	// Matches the root path ('/'), all pathes not passing a file name
	// at the end (e.g.: '/vplan') and all paths ending with '.html'
	// or '.xhtml' (e.g.: 'vplan/index.html').
	s.router.HandleFunc(`/{file:$|[\w+\/]+|[\w\/]+.x?html$}`, s.handlerFEMainRoot)

	// ---------------------------
	// STATIC FRONTEND FILES

	// All other pathes will be interpreted as static file accesses.
	// If there is no coresponding file to the path requested, an
	// error 404 will be returned.
	s.router.Handle("/{file:.+}", http.FileServer(http.Dir(s.config.StaticFiles)))
}

// jsonResponse sends a response containing the response code
// and the data content transformed as JSON data
//   w    : response rwiter
//   code : HTTP status code
//   data : content data
func jsonResponse(w http.ResponseWriter, code int, data interface{}) error {
	var bData []byte
	var err error

	w.Header().Add("Content-Type", "application/json")

	if data != nil {
		bData, err = json.MarshalIndent(data, "", "  ")
	}

	if err != nil {
		return err
	}

	w.WriteHeader(code)
	_, err = w.Write(bData)

	return err
}

// parseJSONBody parses a JSON body to a struct pointer
//   body : Read closer of the request body
//   v    : pointer to the data object
func (s *Server) parseJSONBody(body io.ReadCloser, v interface{}) error {
	dec := json.NewDecoder(body)
	return dec.Decode(v)
}

// checkAuth checks if the current request is authorized by
// a valid session token or by a passed API access token
func (s *Server) checkAuth(w http.ResponseWriter, r *http.Request, authToken string) (string, error) {
	session, err := s.store.Get(r, "main")
	if err != nil {
		return "", err
	}

	if session.IsNew {
		return "", ErrUnauthorized
	}

	_ident, ok := session.Values["ident"]
	if !ok || _ident == "" {
		return "", ErrUnauthorized
	}

	ident, ok := _ident.(string)
	if !ok {
		return "", errors.New("failed parsing ident to string")
	}

	return ident, nil
}
