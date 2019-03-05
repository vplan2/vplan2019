// Package drivers (auth/drivers) contains various general
// drivers for accessing authentication services
//   Authors: Ringo Hoffmann
package drivers

import (
	"strconv"
	"strings"

	"github.com/zekroTJA/vplan2019/internal/logger"

	"gopkg.in/ldap.v2"

	"github.com/zekroTJA/vplan2019/internal/auth"
)

type LDAPOptions struct {
	BaseDN string
	Host   string
	Port   string
	UseSSL bool
}

// DebugAuthProvider is an auth provider, which
// is only purposed to use in debugging and testing
type LDAPAuthProvider struct {
	cfg  map[string]string
	opts *LDAPOptions
	conn *ldap.Conn
}

// Connect _
func (d *LDAPAuthProvider) Connect(options map[string]string) error {
	var err error
	d.cfg = options
	d.opts = new(LDAPOptions)

	// iPort, err := strconv.Atoi(options["port"])
	// if err != nil {
	// 	return err
	// }

	d.opts.BaseDN = options["base"]
	d.opts.Host = options["host"]
	d.opts.Port = options["port"]

	d.opts.UseSSL, err = strconv.ParseBool(options["usessl"])
	if err != nil {
		return err
	}

	d.conn, err = ldap.Dial("tcp", d.opts.Host+":"+d.opts.Port)
	if err != nil {
		return err
	}

	// d.client = &ldap.LDAPClient{
	// 	Base:        options["base"],
	// 	Host:        options["host"],
	// 	Port:        iPort,
	// 	UseSSL:      bSSL,
	// 	UserFilter:  "(uid=%s)",
	// 	GroupFilter: "(memberUid=%s)",
	// }

	// if err = d.client.Connect(); err != nil {
	// 	return err
	// }

	return nil
}

// Close _
func (d *LDAPAuthProvider) Close() {
	d.conn.Close()
}

// GetConfigModel _
func (d *LDAPAuthProvider) GetConfigModel() map[string]string {
	// TODO: LDAP Result Code 200 "Network Error": tls: either ServerName or InsecureSkipVerify must be specified in the tls.Config
	return map[string]string{
		"base":   "dc=example,dc=com",
		"host":   "ldap.example.com",
		"port":   "389",
		"usessl": "false",
	}
}

// Authenticate _
func (d *LDAPAuthProvider) Authenticate(username, group, password string) (*auth.Response, error) {
	// ok, data, err := d.client.Authenticate(username, password)
	// fmt.Println(ok, data, err)
	// if err != nil {
	// 	return nil, err
	// }
	// if !ok {
	// 	return nil, errors.New("unauthorized")
	// }

	// fmt.Println(data)

	dnArr := []string{"cn=" + username}
	if group != "" {
		dnArr = append(dnArr, "ou="+group)
	}
	dnArr = append(dnArr, d.opts.BaseDN)
	dn := strings.Join(dnArr, ",")

	err := d.conn.Bind(dn, password)
	if err != nil {
		logger.Debug("LDAP auth error: ", err)
		return nil, err
	}

	return nil, nil
}
