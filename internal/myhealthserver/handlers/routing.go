package handlers

import (
	"github.com/devldavydov/myhealth/internal/myhealthserver/handlers/food"
	"github.com/devldavydov/myhealth/internal/myhealthserver/handlers/settings"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(router *gin.Engine, stg storage.Storage, userID int64, logger *zap.Logger) {
	router.GET("/", Index)
	router.NoRoute(NotFound)

	food.Attach(router.Group("/food"), stg, userID, logger)
	settings.Attach(router.Group("/settings"), stg, userID, logger)
}
