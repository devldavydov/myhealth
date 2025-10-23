package food

import (
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Attach(group *gin.RouterGroup, stg storage.Storage, userID int64, logger *zap.Logger) {
	foodHandler := NewFoodHander(stg, userID, logger)

	group.GET("/", foodHandler.ListPage)
	group.GET("/edit", foodHandler.EditPage)

	group.GET("/api/list", foodHandler.GetListAPI)
	group.GET("/api/get/:key", foodHandler.GetFoodAPI)
}
