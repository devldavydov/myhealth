package usersettings

import (
	"context"
	"errors"
	"net/http"

	"github.com/devldavydov/myhealth/internal/common/messages"
	"github.com/devldavydov/myhealth/internal/myhealthserver/model"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserSettingsHandler struct {
	stg    storage.Storage
	userID int64
	logger *zap.Logger
}

func NewUserSettingsHandler(stg storage.Storage, userID int64, logger *zap.Logger) *UserSettingsHandler {
	return &UserSettingsHandler{stg: stg, userID: userID, logger: logger}
}

type UserSettings struct {
	CalLimit float64 `json:"calLimit"`
}

func (r *UserSettingsHandler) GetAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	us, err := r.stg.GetUserSettings(ctx, r.userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserSettingsNotFound) {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(messages.MsgErrUserSettingsNotFound))
			return
		}

		r.logger.Error(
			"user settings get api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, model.NewDataResponse(UserSettings{
		CalLimit: us.CalLimit,
	}))
}

func (r *UserSettingsHandler) SetAPI(c *gin.Context) {}
