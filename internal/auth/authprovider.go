// Package auth provides general structures and utilities
// for authentication at various user authentication
// services
//   Authors: Ringo Hoffmann
package auth

// Response is an object which contains the Ident,
// which is a user unique string used for session
// and user data association, and an additional
// context depending on the authentication service
type Response struct {
	Ident string
	Ctx   interface{}
}

// Provider is an interface for authentication
// provider drivers
type Provider interface {
	// Connect to the authentication service with
	// the passed options
	Connect(options map[string]string) error
	// Close connection
	Close()

	// Get map which defines the key-value config
	// model structure
	GetConfigModel() map[string]string
	// Get authentication response by passed
	// username and password credentials.
	// If authentication failes, an error must
	// be returned with an empty response object.
	Authenticate(username, group, password string) (*Response, error)
}
