package food

import (
	cc "github.com/devldavydov/myhealth/internal/myhealthserver/constants"
	rr "github.com/devldavydov/myhealth/internal/myhealthserver/response"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FoodHandler struct {
	stg    storage.Storage
	userID int64
	logger *zap.Logger
}

func NewFoodHander(stg storage.Storage, userID int64, logger *zap.Logger) *FoodHandler {
	return &FoodHandler{stg: stg, userID: userID, logger: logger}
}

func (r *FoodHandler) ListPage(c *gin.Context) {
	rr.
		NewPageResponse(
			cc.TotalConstants["Page_Food_FoodList"],
			"/static/myhealth/js/food/list.js").
		OK(c)
}

func (r *FoodHandler) EditPage(c *gin.Context) {
	rr.
		NewPageResponse(
			cc.TotalConstants["Page_Food_FoodEdit"],
			"/static/myhealth/js/food/edit.js").
		OK(c)
}

func (r *FoodHandler) CreatePage(c *gin.Context) {
	rr.
		NewPageResponse(
			cc.TotalConstants["Page_Food_FoodCreate"],
			"/static/myhealth/js/food/create.js").
		OK(c)
}
