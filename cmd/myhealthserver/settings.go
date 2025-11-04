package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	srv "github.com/devldavydov/myhealth/internal/myhealthserver"
)

const (
	_defaultRunAddress      = "http://127.0.0.1:8080"
	_defaultShutdownTimeout = 15 * time.Second
	_defaultLogLevel        = "INFO"
	_defaultDBFilePath      = ""
	_defaultUserID          = -1
	_defaultTLSCertFile     = ""
	_defaultTLSKeyFile      = ""
)

type Config struct {
	RunAddress      string
	ShutdownTimeout time.Duration
	DBFilePath      string
	LogLevel        string
	UserID          int64
	TLSCertFile     string
	TLSKeyFile      string
}

func LoadConfig(flagSet flag.FlagSet, flags []string) (*Config, error) {
	config := &Config{}

	flagSet.StringVar(&config.RunAddress, "a", _defaultRunAddress, "Server run address")
	flagSet.StringVar(&config.DBFilePath, "d", _defaultDBFilePath, "DB file path")
	flagSet.StringVar(&config.LogLevel, "l", _defaultLogLevel, "Log level")
	flagSet.Int64Var(&config.UserID, "u", _defaultUserID, "User ID")
	flagSet.DurationVar(&config.ShutdownTimeout, "t", _defaultShutdownTimeout, "Server shutdown timeout")
	flagSet.StringVar(&config.TLSCertFile, "c", _defaultTLSCertFile, "TLS cert file")
	flagSet.StringVar(&config.TLSKeyFile, "k", _defaultTLSKeyFile, "TLS key file")

	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flagSet.PrintDefaults()
	}

	err := flagSet.Parse(flags)
	if err != nil {
		return nil, err
	}

	if config.DBFilePath == _defaultDBFilePath {
		return nil, fmt.Errorf("invalid DB file path")
	}

	if config.UserID == _defaultUserID {
		return nil, fmt.Errorf("invalid user ID")
	}

	if config.TLSCertFile == _defaultTLSCertFile {
		return nil, fmt.Errorf("invalid TLS cert file")
	}

	if config.TLSKeyFile == _defaultTLSKeyFile {
		return nil, fmt.Errorf("invalid TLS key file")
	}

	return config, nil
}

func ServiceSettingsAdapt(config *Config) (*srv.ServerSettings, error) {
	return srv.NewServerSettings(
		config.RunAddress,
		config.DBFilePath,
		config.ShutdownTimeout,
		config.UserID,
		config.TLSCertFile,
		config.TLSKeyFile)
}
