// Package config contains the main utilities for
// parsing and generating the main servers config
// file
//   Authors: Ringo Hoffmann
package config

import (
	"io/ioutil"
	"os"

	"github.com/zekroTJA/vplan2019/internal/webserver"
)

// Config contains the general Server config
type Config struct {
	Logging   *Logging          `json:"logging"`
	WebServer *webserver.Config `json:"webServer"`
}

// Logging contains the configuration for logging
type Logging struct {
	Level int `json:"level"`
}

// UnmarshalFunc is a function which can be used to parse
// raw config data to an object instance
type UnmarshalFunc func(data []byte, v interface{}) error

// MarshalIndentFunc is a function which can be used to parse
// an objects content to a config file byte output
type MarshalIndentFunc func(v interface{}, prefix, indent string) ([]byte, error)

// Open reads from the passed file and parses the
// raw data to an Config object by the passed parsing
// function.
//   file          : the file name and path of the file to read from
//   unmarshalFunc : the parsing function which should be used for
//                   parsing the raw config data
func Open(file string, unmarshalFunc UnmarshalFunc) (*Config, error) {
	bData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = unmarshalFunc(bData, config)

	return config, err
}

// Create writes the passed config object content (or an empty config, if passed config is nil)
// to a new or existing file by using the passed MarshalIndentFunc.
//   file        : the fil name and path of the file to write to
//   config      : config object to write (if nil, an empty config will be used)
//   prefix      : prefix which will be directly passed to the MarhsalIndentFunc
//   indent      : indent which will be directly passed to the MarhsalIndentFunc
//   marshalFunc : function which will be used to parse the config object content
//                 to the formatted data which will be written to the created file
func Create(file string, config *Config, prefix, indent string, marshalFunc MarshalIndentFunc) error {
	if config == nil {
		config = &Config{
			Logging: new(Logging),
			WebServer: &webserver.Config{
				TLS: new(webserver.ConfigTLS),
			},
		}
	}

	fHandle, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fHandle.Close()

	bData, err := marshalFunc(config, prefix, indent)
	if err != nil {
		return err
	}

	_, err = fHandle.Write(bData)
	return err
}
