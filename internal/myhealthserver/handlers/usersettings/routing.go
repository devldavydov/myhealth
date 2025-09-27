package usersettings

import (
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Attach(group *gin.RouterGroup, stg storage.Storage, userID int64, logger *zap.Logger) {
	usHandler := NewUserSettingsHandler(stg, userID, logger)

	group.GET("/", usHandler.GetAPI)
	group.POST("/", usHandler.SetAPI)
}
