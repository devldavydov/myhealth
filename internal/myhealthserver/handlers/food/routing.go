package food

import (
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Attach(group *gin.RouterGroup, stg storage.Storage, userID int64, logger *zap.Logger) {
	foodHandler := NewFoodHander(stg, userID, logger)

	group.GET("/", foodHandler.ListAPI)
	group.GET("/:key", foodHandler.GetAPI)
	group.DELETE("/:key", foodHandler.DeleteAPI)
	group.POST("/set", foodHandler.SetAPI)
}
