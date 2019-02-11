// Package webserver contains general utilities
// for exposing and handling a web server.
//   Authors: Ringo Hoffmann
package webserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Config contains the configuration
// for the web server
type Config struct {
	Addr string     `json:"addr"`
	TLS  *ConfigTLS `json:"tls"`
}

// ConfigTLS contains the cert file path
// and the key file path for TLS configuration
type ConfigTLS struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// Server contains the instance of the
// http server and the mux router
type Server struct {
	server *http.Server
	router *mux.Router
	store  sessions.Store
}

var (
	// Errors
	errServerInstanceNil = errors.New("web server instance must be initialized")
	// Vars
	defaultConfig = &Config{":80", nil}
)

// StartBlocking starts the web server and block the current thread
//   server : the initialized instance of an empty Server struct
//   config : instance of Config; if this is nil, defaultConfig will be used
func StartBlocking(server *Server, config *Config, sessionStorage sessions.Store) error {
	if server == nil {
		return errServerInstanceNil
	}

	if config == nil {
		config = defaultConfig
	}

	server.router = mux.NewRouter()
	server.store = sessionStorage

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
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "test")
		fmt.Println(session.Values["a"])
		session.Values["a"] = "b"
		session.Save(r, w)
		w.Write([]byte("hey"))
	}).Methods("GET")
}
