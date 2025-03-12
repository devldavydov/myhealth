package cmdproc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
)

func (r *CmdProcessor) processUserSettings(baseCmd string, cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	// Sport
	case "set":
		resp = r.userSettingsSetCommand(cmdParts[1:], userID)
	case "get":
		resp = r.userSettingsGetCommand(userID)
	case "st":
		resp = r.userSettingsSetTemplateCommand(userID)
	case "h":
		resp = r.userSettingsHelpCommand(baseCmd)
	default:
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		resp = NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) userSettingsSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	calLimit, err := strconv.ParseFloat(cmdParts[0], 64)
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

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
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) userSettingsSetCommand2(userID int64, calLimit float64) []CmdResponse {
	return nil
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

func (r *CmdProcessor) userSettingsHelpCommand(baseCmd string) []CmdResponse {
	return NewSingleCmdResponse(
		newCmdHelpBuilder(baseCmd, "Управление настройками пользователя").
			addCmd(
				"Установка",
				"set",
				"Лимит калорий [Дробное>0]",
			).
			addCmd(
				"Шаблон команды установки",
				"st",
			).
			addCmd(
				"Получение",
				"get",
			).
			build(),
		optsHTML)
}
