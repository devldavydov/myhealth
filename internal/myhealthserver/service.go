package myhealthserver

import (
	"context"

	"github.com/devldavydov/myhealth/internal/storage"
	slite "github.com/devldavydov/myhealth/internal/storage/sqlite"
	"go.uber.org/zap"
)

type Service struct {
	settings *ServerSettings
	logger   *zap.Logger
	stg      storage.Storage
}

func NewService(settings *ServerSettings, logger *zap.Logger) (*Service, error) {
	stg, err := slite.NewStorageSQLite(settings.DBFilePath, logger)
	if err != nil {
		return nil, err
	}

	return &Service{settings: settings, stg: stg, logger: logger}, nil
}

func (r *Service) Run(ctx context.Context) error {
	return nil
}
