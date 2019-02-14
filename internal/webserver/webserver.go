// Package webserver contains general utilities
// for exposing and handling a web server.
//   Authors: Ringo Hoffmann
package webserver

import (
	"encoding/json"
	"errors"
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
	server       *http.Server
	router       *mux.Router
	store        sessions.Store
	authProvider auth.Provider
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

	server.router = mux.NewRouter()
	server.store = sessionStorage
	server.authProvider = authProvider

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

// initializeHandlers contains all setup functions for router
// endpoints and their handlers
func (s *Server) initializeHnalders() {
	// Frontent
	s.router.HandleFunc("/", s.handlerMainPage).Methods("GET")

	// API
	s.router.HandleFunc("/api/authenticate/{username}", s.handlerAPIAuthenticate).Methods("POST")

	// Serve static files from './web/static'
	s.router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("./web/static/"))))
}

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
