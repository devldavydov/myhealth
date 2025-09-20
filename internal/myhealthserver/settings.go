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
}

func NewServerSettings(
	runAddress string,
	dbFilePath string,
	shutdownTimeout time.Duration,
	userID int64,
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
	}, nil
}
