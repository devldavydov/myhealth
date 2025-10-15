package myhealthserver

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/devldavydov/myhealth/internal/myhealthserver/handlers"
	"github.com/devldavydov/myhealth/internal/storage"
	slite "github.com/devldavydov/myhealth/internal/storage/sqlite"
	"go.uber.org/zap"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Service struct {
	settings *ServerSettings
	logger   *zap.Logger
	stg      storage.Storage
}

//go:embed templates/* static/*
var embedFS embed.FS

func NewService(settings *ServerSettings, logger *zap.Logger) (*Service, error) {
	stg, err := slite.NewStorageSQLite(settings.DBFilePath, logger)
	if err != nil {
		return nil, err
	}

	return &Service{settings: settings, stg: stg, logger: logger}, nil
}

func (r *Service) Run(ctx context.Context) error {
	// Init HTTP
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	files, err := loadTemplates("templates")
	if err != nil {
		return err
	}
	router.LoadHTMLFS(http.FS(embedFS), files...)

	staticFS, err := fs.Sub(embedFS, "static")
	if err != nil {
		return err
	}
	router.StaticFS("/static", http.FS(staticFS))

	handlers.Init(router, r.stg, r.settings.UserID, r.logger)

	// Start server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", r.settings.RunAddress.Hostname(), r.settings.RunAddress.Port()),
		Handler: router,
	}

	errChan := make(chan error)
	go func(ch chan error) {
		ch <- httpServer.ListenAndServe()
	}(errChan)

	select {
	case err := <-errChan:
		return fmt.Errorf("service exited with err: %w", err)
	case <-ctx.Done():
		r.logger.Info("Service context canceled")

		ctx, cancel := context.WithTimeout(context.Background(), r.settings.ShutdownTimeout)
		defer cancel()

		err := httpServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("service shutdown err: %w", err)
		}

		r.logger.Info("Service finished")
		return nil
	}
}

func loadTemplates(root string) (files []string, err error) {
	err = fs.WalkDir(embedFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fi, err := d.Info()
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if path != root {
				loadTemplates(path)
			}
		} else {
			files = append(files, path)
		}

		return err
	})

	return files, err
}
