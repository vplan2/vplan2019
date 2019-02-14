package main

import (
	"flag"
	"os"

	"github.com/ghodss/yaml"

	"github.com/zekroTJA/vplan2019/internal/config"
	"github.com/zekroTJA/vplan2019/internal/logger"
	"github.com/zekroTJA/vplan2019/internal/webserver"

	authDrivers "github.com/zekroTJA/vplan2019/internal/auth/drivers"
	dbDrivers "github.com/zekroTJA/vplan2019/internal/database/drivers"
)

var (
	flagConfig = flag.String("c", "config.yml", "location of the config file")
)

func main() {
	flag.Parse()

	database := new(dbDrivers.SQLite)
	authProvider := new(authDrivers.DebugAuthProvider)

	//////////////////
	// LOGGER SETUP //
	//////////////////

	logger.Setup(`%{color}▶  %{level:.4s} %{id:03x}%{color:reset} %{message}`, 5)

	////////////////////
	// CONFIG PARSING //
	////////////////////

	// Define unmarshal function for config decoding
	unmarshalFunc := func(data []byte, v interface{}) error {
		return yaml.Unmarshal(data, v)
	}
	// Define marshal function for config encoding
	marshalFunc := func(v interface{}, prefix, indent string) ([]byte, error) {
		return yaml.Marshal(v)
	}
	// try to laod existing config
	cfg, err := config.Open(*flagConfig, unmarshalFunc)
	// If it was a file not found error, try to create a new config file
	if os.IsNotExist(err) {
		providerModels := &config.ProviderModels{
			Database:      database.GetConfigModel(),
			Authorization: authProvider.GetConfigModel(),
		}

		err = config.Create(*flagConfig, nil, providerModels, "", "  ", marshalFunc)
		if err != nil {
			logger.Fatal("Failed creating config: ", err)
		}

		logger.Info("Created new empty config at '%s'. Enter your preferenced values and restart.", *flagConfig)
		return
	}
	// If config laoding set an other error, throw and exit
	if err != nil {
		logger.Fatal("Failed parsing config: ", err)
	}
	// Set log level from config
	logger.SetLogLevel(cfg.Logging.Level)

	////////////////////
	// DATABASE SETUP //
	////////////////////

	err = database.Connect(cfg.Providers.Database)
	if err != nil {
		logger.Fatal("Failed creating database connection: ", err)
	}

	//////////////////////
	// WEB SERVER SETUP //
	//////////////////////

	// output web server starting info and warn if web server was
	// started in non TLS mode
	logger.Info("Starting web server on %s...", cfg.WebServer.Addr)
	if cfg.WebServer.TLS == nil {
		logger.Warning("ATTENTION: THE WEB SERVER IS NOT RUNNING IN TLS MODE")
	}
	// Set session storage module
	store, err := database.GetSessionStoreDriver(
		cfg.WebServer.Sessions.DefaultMaxAge, []byte(cfg.WebServer.Sessions.EncryptionSecret))
	if err != nil {
		logger.Fatal("Failed getting session driver: ", err)
	}
	// Create server instance
	server := new(webserver.Server)
	// Initiate and run the web server, which blocks the main thread.
	// If it fails, thow error and exit.
	logger.Fatal("Failed opening webserver: ",
		webserver.StartBlocking(server, cfg.WebServer, store))
}
