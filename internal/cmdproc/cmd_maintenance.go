package cmdproc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"time"

	m "github.com/devldavydov/myhealth/internal/common/messages"
	tele "gopkg.in/telebot.v4"

	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
)

func (r *CmdProcessor) maintenanceBackupCommand(userID int64) []CmdResponse {
	// Get backup from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout*10)
	defer cancel()

	backup, err := r.stg.Backup(ctx)
	if err != nil {
		r.logger.Error(
			"backup command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	// Generate response.
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if err := json.NewEncoder(zw).Encode(&backup); err != nil {
		r.logger.Error(
			"backup json err",
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	if err := zw.Close(); err != nil {
		r.logger.Error(
			"gzip err",
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(&buf),
		MIME:     "application/x-gzip-compressed",
		FileName: fmt.Sprintf("backup_%s.json.gz", formatTimestamp(time.Now().In(r.tz))),
	})
}
