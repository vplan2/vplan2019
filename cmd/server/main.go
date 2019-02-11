package main

import (
	"flag"
	"os"

	"github.com/zekroTJA/vplan2019/internal/webserver"

	"github.com/ghodss/yaml"

	"github.com/zekroTJA/vplan2019/internal/config"
	"github.com/zekroTJA/vplan2019/internal/logger"
)

var (
	flagConfig = flag.String("c", "config.yml", "location of the config file")
)

func main() {
	flag.Parse()

	unmarshalFunc := func(data []byte, v interface{}) error {
		return yaml.Unmarshal(data, v)
	}

	marshalFunc := func(v interface{}, prefix, indent string) ([]byte, error) {
		return yaml.Marshal(v)
	}

	cfg, err := config.Open(*flagConfig, unmarshalFunc)

	if os.IsNotExist(err) {
		err = config.Create(*flagConfig, nil, "", "  ", marshalFunc)
		if err != nil {
			logger.Fatal("Failed creating config: ", err)
		}
		logger.Infof("Created new empty config at '%s'. Enter your preferenced values and restart.", *flagConfig)
		return
	}

	if err != nil {
		logger.Fatal("Failed parsing config: ", err)
	}

	logger.Setup(`%{color}▶  %{level:.4s} %{id:03x}%{color:reset} %{message}`, cfg.Logging.Level)

	logger.Infof("Starting web server on %s...", cfg.WebServer.Addr)
	if cfg.WebServer.TLS == nil {
		logger.Warning("ATTENTION: THE WEB SERVER IS NOT RUNNING IN TLS MODE")
	}

	server := new(webserver.Server)
	logger.Fatal("Failed opening webserver: ", webserver.StartBlocking(server, cfg.WebServer))
}
