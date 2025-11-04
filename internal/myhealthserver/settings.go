package myhealthserver

import (
	"net/url"
	"time"
)

type ServerSettings struct {
	RunAddress      *url.URL
	DBFilePath      string
	ShutdownTimeout time.Duration
	UserID          int64
	TLSCertFile     string
	TLSKeyFile      string
}

func NewServerSettings(
	runAddress string,
	dbFilePath string,
	shutdownTimeout time.Duration,
	userID int64,
	tlsCertFile string,
	tlsKeyFile string,
) (*ServerSettings, error) {

	urlRunAddress, err := url.ParseRequestURI(runAddress)
	if err != nil {
		return nil, err
	}

	return &ServerSettings{
		RunAddress:      urlRunAddress,
		DBFilePath:      dbFilePath,
		ShutdownTimeout: shutdownTimeout,
		UserID:          userID,
		TLSCertFile:     tlsCertFile,
		TLSKeyFile:      tlsKeyFile,
	}, nil
}
