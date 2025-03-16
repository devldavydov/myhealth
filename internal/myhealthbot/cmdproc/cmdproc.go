//go:generate go run ./gen/gen.go -in commands.yaml -out cmdproc_generated.go

package cmdproc

import (
	"strings"
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
	cmdParts := []string{}
	for _, part := range strings.Split(cmd, ",") {
		cmdParts = append(cmdParts, strings.Trim(part, " "))
	}

	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.String("command", cmd),
			zap.Int64("userID", userID),
		)
		return c.Send(MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "h":
		resp = r.processHelp(userID)
	case "s":
		resp = r.processSport("s", cmdParts[1:], userID)
	default:
		r.logger.Error(
			"unknown command",
			zap.String("command", cmd),
			zap.Int64("userID", userID),
		)
		resp = NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	if r.debugMode {
		if err := c.Send("!!! ОТЛАДОЧНЫЙ РЕЖИМ !!!"); err != nil {
			return err
		}
	}

	for _, rItem := range resp {
		if err := c.Send(rItem.what, rItem.opts...); err != nil {
			return err
		}
	}

	return nil
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

func (r *CmdProcessor) parseTimestamp(sTimestamp string) (time.Time, error) {
	var t time.Time
	var err error

	if sTimestamp == "" {
		t = time.Now().In(r.tz)
	} else {
		t, err = time.Parse("02.01.2006", sTimestamp)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, r.tz), nil
}

func formatTimestamp(ts time.Time) string {
	return ts.Format("02.01.2006")
}
