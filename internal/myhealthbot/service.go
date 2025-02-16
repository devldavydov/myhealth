package myhealthbot

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/devldavydov/myhealth/internal/myhealthbot/cmdproc"
	s "github.com/devldavydov/myhealth/internal/storage"
	slite "github.com/devldavydov/myhealth/internal/storage/sqlite"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

const (
	_backupFile = "backup.tar.gz"
)

type Service struct {
	settings *ServiceSettings
	cmdProc  *cmdproc.CmdProcessor
	logger   *zap.Logger
}

func NewService(settings *ServiceSettings, logger *zap.Logger) (*Service, error) {
	stg, err := slite.NewStorageSQLite(settings.DBFilePath, logger)
	if err != nil {
		return nil, err
	}

	srv := &Service{
		settings: settings,
		cmdProc:  cmdproc.NewCmdProcessor(stg, settings.TZ, settings.DebugMode, logger),
		logger:   logger,
	}

	if err := srv.tryRestoreFromBackup(stg); err != nil {
		return nil, err
	}

	return srv, nil
}

func (r *Service) Run(ctx context.Context) error {
	pref := tele.Settings{
		Token:  r.settings.Token,
		Poller: &tele.LongPoller{Timeout: r.settings.PollTimeOut},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return err
	}

	r.setupRouting(b, r.settings.AllowedUserIDs)
	go b.Start()

	<-ctx.Done()
	b.Stop()
	r.cmdProc.Stop()

	return nil
}

func (r *Service) setupRouting(b *tele.Bot, allowedUserIDs []int64) {
	b.Handle("/start", r.onStart)

	allowedGroup := b.Group()
	allowedGroup.Use(middleware.Whitelist(allowedUserIDs...))
	allowedGroup.Handle(tele.OnText, r.onText)
}

func (r *Service) onStart(c tele.Context) error {
	return c.Send(
		fmt.Sprintf(
			"Привет, %s [%d]!\nДобро пожаловать в MyHealthBot!\nОтправь 'h' для помощи",
			c.Sender().Username,
			c.Sender().ID,
		),
	)
}

func (r *Service) onText(c tele.Context) error {
	return r.cmdProc.Process(c, c.Text(), c.Sender().ID)
}

func (r *Service) tryRestoreFromBackup(stg s.Storage) error {
	r.logger.Info("trying to restore from backup")

	f, err := os.Open(_backupFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			r.logger.Info("backup file not found")
			return nil
		}

		r.logger.Error("open backup file error", zap.Error(err))
		return err
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		r.logger.Error("backup gzip read error", zap.Error(err))
		return err
	}
	defer gr.Close()

	var backup s.Backup
	if err := json.NewDecoder(gr).Decode(&backup); err != nil {
		r.logger.Error("backup json read error", zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.StorageRestoreTimeout)
	defer cancel()

	if err := stg.Restore(ctx, &backup); err != nil {
		r.logger.Error("backup restore error", zap.Error(err))
		return err
	}

	r.logger.Info("backup finished")
	return nil
}
