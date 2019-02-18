// Package webserver contains general utilities
// for exposing and handling a web server.
//   Authors: Ringo Hoffmann
package webserver

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/zekroTJA/vplan2019/internal/auth"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Config contains the configuration
// for the web server
type Config struct {
	Addr     string          `json:"addr"`
	Sessions *ConfigSessions `json:"sessions"`
	TLS      *ConfigTLS      `json:"tls"`
}

// ConfigTLS contains the cert file path
// and the key file path for TLS configuration
type ConfigTLS struct {
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
func StartBlocking(server *Server, config *Config, sessionStorage sessions.Store, authProvider auth.Provider) error {
	if server == nil {
		return errServerInstanceNil
	}

	if config == nil {
		return errConfigNil
	}

	server.config = config
	server.router = mux.NewRouter()
	server.store = sessionStorage
	server.authProvider = authProvider

	server.limiter = NewRateLimiter(&LimiterOpts{10, 10}, server.handlerAPIRateLimitError)

	server.server = &http.Server{
		Addr:    config.Addr,
		Handler: server.router,
	}

	server.initializeHnalders()

	if config.TLS == nil {
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
func (s *Server) addHandler(path string, ident string, handler func(w http.ResponseWriter, r *http.Request), limiterRate, limiterBurst int, methods ...string) {
	s.router.HandleFunc(path, handler).Methods(methods...)
	s.limiter.Register(ident, limiterRate, limiterBurst)
}

// initializeHandlers contains all setup functions for router
// endpoints and their handlers
func (s *Server) initializeHnalders() {

	// FRONTEND
	s.router.HandleFunc("/", s.handlerMainPage).Methods("GET")

	// API
	s.addHandler("/api/authenticate/{username}", "authenticate",
		s.handlerAPIAuthenticate, 1, 3, "POST")
	s.addHandler("/api/test", "test", s.handlerAPITest, 1, 1, "POST")

	// Serve static files from './web/static'
	s.router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("./web/static/"))))
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
