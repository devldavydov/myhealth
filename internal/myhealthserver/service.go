//go:generate go run ./gen/gen.go -in service.yaml

package myhealthserver

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/devldavydov/myhealth/internal/myhealthserver/handlers"
	p "github.com/devldavydov/myhealth/internal/myhealthserver/process"
	"github.com/devldavydov/myhealth/internal/storage"
	slite "github.com/devldavydov/myhealth/internal/storage/sqlite"
	"go.uber.org/zap"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Service struct {
	settings     *ServerSettings
	cmdProceccor *cmdproc.CmdProcessor
	logger       *zap.Logger
	stg          storage.Storage
	wg           sync.WaitGroup
}

//go:embed templates/* static/*
var embedFS embed.FS

func NewService(settings *ServerSettings, logger *zap.Logger) (*Service, error) {
	stg, err := slite.NewStorageSQLite(settings.DBFilePath, logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		settings: settings,
		cmdProceccor: cmdproc.NewCmdProcessor(
			stg,
			p.NewTypeAdapter(),
			settings.TZ,
			settings.DebugMode,
			logger),
		stg:    stg,
		logger: logger}, nil
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

	handlers.Init(router, r.cmdProceccor, r.settings.FileStoragePath, r.settings.UserID)

	r.wg.Add(1)
	go r.filesCleanJob(ctx)

	// Start server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", r.settings.RunAddress.Hostname(), r.settings.RunAddress.Port()),
		Handler: router,
	}

	errChan := make(chan error)
	go func(ch chan error) {
		ch <- httpServer.ListenAndServeTLS(r.settings.TLSCertFile, r.settings.TLSKeyFile)
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

func (r *Service) filesCleanJob(ctx context.Context) {
	defer r.wg.Done()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	cleanOldFiles := func() {
		now := time.Now()
		threshold := 24 * time.Hour

		files, err := os.ReadDir(r.settings.FileStoragePath)
		if err != nil {
			r.logger.Error(
				"failed to read file storage path",
				zap.String("fileStoragePath", r.settings.FileStoragePath),
				zap.Error(err))
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			info, err := file.Info()
			if err != nil {
				continue
			}

			if now.Sub(info.ModTime()) > threshold {
				fullPath := filepath.Join(r.settings.FileStoragePath, file.Name())
				err := os.Remove(fullPath)
				if err != nil {
					r.logger.Info(
						"failed to delete file",
						zap.String("path", fullPath),
						zap.Error(err))
				} else {
					r.logger.Info(
						"file deleted",
						zap.String("path", fullPath),
					)
				}
			}
		}
	}

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("files clean job context canceled")
			return
		case <-ticker.C:
			cleanOldFiles()
		}
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
