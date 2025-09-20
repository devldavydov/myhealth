package handlers

import (
	"github.com/devldavydov/myhealth/internal/myhealthserver/handlers/food"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(router *gin.Engine, stg storage.Storage, userID int64, logger *zap.Logger) {
	api := router.Group("/api")

	food.Attach(api.Group("/food"), stg, userID, logger)
}
