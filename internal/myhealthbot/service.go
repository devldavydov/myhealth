package myhealthbot

import (
	"context"
	"fmt"

	"github.com/devldavydov/myhealth/internal/myhealthbot/cmdproc"
	storage "github.com/devldavydov/myhealth/internal/storage/sqlite"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

type Service struct {
	settings *ServiceSettings
	cmdProc  *cmdproc.CmdProcessor
}

func NewService(settings *ServiceSettings, logger *zap.Logger) (*Service, error) {
	stg, err := storage.NewStorageSQLite(settings.DBFilePath, logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		settings: settings,
		cmdProc:  cmdproc.NewCmdProcessor(stg, settings.TZ, settings.DebugMode, logger),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
	pref := tele.Settings{
		Token:  s.settings.Token,
		Poller: &tele.LongPoller{Timeout: s.settings.PollTimeOut},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return err
	}

	s.setupRouting(b, s.settings.AllowedUserIDs)
	go b.Start()

	<-ctx.Done()
	b.Stop()
	s.cmdProc.Stop()

	return nil
}

func (s *Service) setupRouting(b *tele.Bot, allowedUserIDs []int64) {
	b.Handle("/start", s.onStart)

	allowedGroup := b.Group()
	allowedGroup.Use(middleware.Whitelist(allowedUserIDs...))
	allowedGroup.Handle(tele.OnText, s.onText)
}

func (s *Service) onStart(c tele.Context) error {
	return c.Send(
		fmt.Sprintf(
			"Привет, %s [%d]!\nДобро пожаловать в MyHealthBot!\nОтправь 'h' для помощи",
			c.Sender().Username,
			c.Sender().ID,
		),
	)
}

func (s *Service) onText(c tele.Context) error {
	return s.cmdProc.Process(c, c.Text(), c.Sender().ID)
}
