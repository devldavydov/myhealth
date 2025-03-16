package cmdproc

import (
	"context"
	"errors"
	"fmt"

	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
)

func (r *CmdProcessor) userSettingsSetCommand(userID int64, calLimit float64) []CmdResponse {
	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetUserSettings(ctx, userID, &storage.UserSettings{
		CalLimit: calLimit,
	}); err != nil {
		if errors.Is(err, storage.ErrUserSettingsInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		r.logger.Error(
			"user settings set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) userSettingsGetCommand(userID int64) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	us, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserSettingsNotFound) {
			return NewSingleCmdResponse(MsgErrUserSettingsNotFound)
		}

		r.logger.Error(
			"user settings get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("<b>Лимит калорий:</b> %.2f", us.CalLimit), optsHTML)
}

func (r *CmdProcessor) userSettingsSetTemplateCommand(userID int64) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	us, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserSettingsNotFound) {
			return NewSingleCmdResponse(MsgErrUserSettingsNotFound)
		}

		r.logger.Error(
			"user settings get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("u,set,%.2f", us.CalLimit))
}
