// Package webserver contains general utilities
// for exposing and handling a web server.
//   Authors: Ringo Hoffmann
package webserver

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
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
func StartBlocking(server *Server, config *Config) error {
	if server == nil {
		return errServerInstanceNil
	}

	if config == nil {
		config = defaultConfig
	}

	server.server = &http.Server{
		Addr: config.Addr,
	}

	server.router = mux.NewRouter()

	server.initializeHnalders()

	if config.TLS == nil {
		return server.server.ListenAndServe()
	}
	return server.server.ListenAndServeTLS(config.TLS.CertFile, config.TLS.KeyFile)
}

// initializeHandlers contains all setup functions for router
// endpoints and their handlers
func (s *Server) initializeHnalders() {

}
