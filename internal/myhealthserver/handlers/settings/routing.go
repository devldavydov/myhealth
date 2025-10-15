package settings

import (
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Attach(group *gin.RouterGroup, stg storage.Storage, userID int64, logger *zap.Logger) {
	sHandler := NewSettingsHandler(stg, userID, logger)

	group.GET("/calccal", sHandler.GetCalcCal)
	group.POST("/calccal", sHandler.PostCalcCal)
}
