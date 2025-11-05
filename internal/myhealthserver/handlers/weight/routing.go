package weight

import (
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Attach(group *gin.RouterGroup, stg storage.Storage, userID int64, logger *zap.Logger) {
	weightHandler := NewWeightHander(stg, userID, logger)

	group.GET("/", weightHandler.ListPage)
	group.GET("/edit", weightHandler.EditPage)
	group.GET("/create", weightHandler.CreatePage)

	group.GET("/api/list", weightHandler.GetListAPI)
	group.GET("/api/:key", weightHandler.GetWeightAPI)
	group.POST("/api/set", weightHandler.SetWeightAPI)
	group.DELETE("/api/:key", weightHandler.DeleteWeightAPI)
}
