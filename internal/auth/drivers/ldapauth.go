// Package drivers (auth/drivers) contains various general
// drivers for accessing authentication services
//   Authors: Ringo Hoffmann
package drivers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zekroTJA/vplan2019/internal/logger"

	"gopkg.in/ldap.v2"

	"github.com/zekroTJA/vplan2019/internal/auth"
)

type LDAPOptions struct {
	BaseDN    string
	Host      string
	Port      string
	Attrbutes []string
	UseSSL    bool
	CertFile  string
	KeyFile   string
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

	d.opts.BaseDN = options["base"]
	d.opts.Host = options["host"]
	d.opts.Port = options["port"]
	d.opts.CertFile = options["certfile"]
	d.opts.KeyFile = options["keyfile"]
	d.opts.Attrbutes = strings.Split(options["attributes"], ",")

	for i, a := range d.opts.Attrbutes {
		d.opts.Attrbutes[i] = strings.Trim(a, " \t")
	}

	d.opts.UseSSL, err = strconv.ParseBool(options["usessl"])
	if err != nil {
		return err
	}

	if d.opts.UseSSL {
		cert, err := tls.LoadX509KeyPair(d.opts.CertFile, d.opts.KeyFile)
		if err != nil {
			return err
		}

		d.conn, err = ldap.DialTLS("tcp", d.opts.Host+":"+d.opts.Port, &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
	} else {
		d.conn, err = ldap.Dial("tcp", d.opts.Host+":"+d.opts.Port)
	}
	if err != nil {
		return err
	}

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
		"base":       "dc=example,dc=com",
		"host":       "ldap.example.com",
		"port":       "389",
		"attributes": "cn, ou, giveName",
		"usessl":     "false",
		"certfile":   "",
		"keyfile":    "",
	}
}

// Authenticate _
func (d *LDAPAuthProvider) Authenticate(username, group, password string) (*auth.Response, error) {
	dnArr := []string{"cn=" + username}
	if group != "" {
		dnArr = append(dnArr, "ou="+group)
	}
	dnArr = append(dnArr, d.opts.BaseDN)
	dn := strings.Join(dnArr, ",")

	err := d.conn.Bind(dn, password)
	if ldap.IsErrorWithCode(err, ldap.ErrorNetwork) {
		if err = d.Connect(d.cfg); err != nil {
			return nil, err
		}
	} else if err != nil {
		logger.Debug("LDAP auth error: %s", err.Error())
		return nil, err
	}

	searchRequest := ldap.NewSearchRequest(
		d.opts.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", username),
		d.opts.Attrbutes,
		nil,
	)

	res, err := d.conn.Search(searchRequest)
	if err != nil {
		logger.Debug("LDAP search error: %s", err.Error())
		return nil, err
	}

	if len(res.Entries) == 0 {
		logger.Debug("LDAP search error: search results are empty")
		return nil, errors.New("search results emptys")
	}

	userRes := res.Entries[0]

	ctx := make(map[string][]string)
	for _, a := range userRes.Attributes {
		ctx[a.Name] = a.Values
	}

	reqRes := &auth.Response{
		Ident: userRes.DN,
		Ctx:   ctx,
	}

	return reqRes, nil
}
