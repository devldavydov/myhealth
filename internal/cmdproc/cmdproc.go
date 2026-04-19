//go:generate go run ./gen/gen.go -in commands.yaml -out cmdproc_generated.go

package cmdproc

import (
	"bytes"
	"time"

	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
)

type ICmdProcess interface {
	Send(what any, opts ...any) error
}

type ITypeAdapter interface {
	File(buf *bytes.Buffer, mime string, fileName string) any
	OptsHTML() any
}

type CmdProcessor struct {
	stg         storage.Storage
	typeAdapter ITypeAdapter
	tz          *time.Location
	logger      *zap.Logger
	debugMode   bool
}

func NewCmdProcessor(
	stg storage.Storage,
	typeAdapter ITypeAdapter,
	tz *time.Location,
	debugMode bool,
	logger *zap.Logger,
) *CmdProcessor {
	return &CmdProcessor{stg: stg, typeAdapter: typeAdapter, tz: tz, debugMode: debugMode, logger: logger}
}

func (r *CmdProcessor) Stop() {
	if err := r.stg.Close(); err != nil {
		r.logger.Error("storage close error", zap.Error(err))
	}
}

func (r *CmdProcessor) Process(c ICmdProcess, cmd string, userID int64) error {
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
