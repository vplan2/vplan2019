// Package drivers (auth/drivers) contains various general
// drivers for accessing authentication services
//   Authors: Ringo Hoffmann
package drivers

import (
	"strconv"

	"github.com/jtblin/go-ldap-client"
	"github.com/zekroTJA/vplan2019/internal/auth"
)

// DebugAuthProvider is an auth provider, which
// is only purposed to use in debugging and testing
type LDAPAuthProvider struct {
	cfg    map[string]string
	client *ldap.LDAPClient
}

// Connect _
func (d *LDAPAuthProvider) Connect(options map[string]string) error {
	d.cfg = options

	iPort, err := strconv.Atoi(options["port"])
	if err != nil {
		return err
	}

	bSSL, err := strconv.ParseBool(options["usessl"])
	if err != nil {
		return err
	}

	d.client = &ldap.LDAPClient{
		Base:   options["base"],
		Host:   options["host"],
		Port:   iPort,
		UseSSL: bSSL,
	}

	if err = d.client.Connect(); err != nil {
		return err
	}
	defer d.client.Close()

	return nil
}

// Close _
func (d *LDAPAuthProvider) Close() {}

// GetConfigModel _
func (d *LDAPAuthProvider) GetConfigModel() map[string]string {
	return map[string]string{
		"base":   "dc=example,dc=com",
		"host":   "ldap.example.com",
		"port":   "389",
		"usessl": "false",
	}
}

// Authenticate _
func (d *LDAPAuthProvider) Authenticate(username, password string) (*auth.Response, error) {
	return nil, nil
}
