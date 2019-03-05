// Package drivers (auth/drivers) contains various general
// drivers for accessing authentication services
//   Authors: Ringo Hoffmann
package drivers

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/zekroTJA/vplan2019/internal/auth"
)

// DebugAuthProvider is an auth provider, which
// is only purposed to use in debugging and testing
type DebugAuthProvider struct {
	cfg   map[string]string
	creds map[string]string
}

// Connect _
func (d *DebugAuthProvider) Connect(options map[string]string) error {
	d.cfg = options
	d.creds = map[string]string{
		"mustermax": "password",
	}
	return nil
}

// Close _
func (d *DebugAuthProvider) Close() {}

// GetConfigModel _
func (d *DebugAuthProvider) GetConfigModel() map[string]string {
	return make(map[string]string)
}

// Authenticate _
func (d *DebugAuthProvider) Authenticate(username, group, password string) (*auth.Response, error) {
	if pw, ok := d.creds[username]; ok && pw == password {
		ident := fmt.Sprintf("%x", sha256.Sum256([]byte(username+password)))
		return &auth.Response{
			Ident: ident,
		}, nil
	}
	return nil, errors.New("unauthorized")
}
