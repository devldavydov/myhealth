//go:generate go run ./gen/gen.go -in commands.yaml -out cmdproc_generated.go

package cmdproc

import (
	"time"

	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

type CmdProcessor struct {
	stg       storage.Storage
	tz        *time.Location
	logger    *zap.Logger
	debugMode bool
}

func NewCmdProcessor(stg storage.Storage, tz *time.Location, debugMode bool, logger *zap.Logger) *CmdProcessor {
	return &CmdProcessor{stg: stg, tz: tz, debugMode: debugMode, logger: logger}
}

func (r *CmdProcessor) Stop() {
	if err := r.stg.Close(); err != nil {
		r.logger.Error("storage close error", zap.Error(err))
	}
}

func (r *CmdProcessor) Process(c tele.Context, cmd string, userID int64) error {
	return r.process(c, cmd, userID)
}

type CmdResponse struct {
	what any
	opts []any
}

func NewCmdResponse(what any, opts ...any) CmdResponse {
	return CmdResponse{what: what, opts: opts}
}

func NewSingleCmdResponse(what any, opts ...any) []CmdResponse {
	return []CmdResponse{
		{what: what, opts: opts},
	}
}
