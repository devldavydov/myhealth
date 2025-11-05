package weight

import (
	cc "github.com/devldavydov/myhealth/internal/myhealthserver/constants"
	rr "github.com/devldavydov/myhealth/internal/myhealthserver/response"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WeightHandler struct {
	stg    storage.Storage
	userID int64
	logger *zap.Logger
}

func NewWeightHander(stg storage.Storage, userID int64, logger *zap.Logger) *WeightHandler {
	return &WeightHandler{stg: stg, userID: userID, logger: logger}
}

func (r *WeightHandler) ListPage(c *gin.Context) {
	rr.
		NewPageResponse(
			cc.TotalConstants["Page_Weight_WeightList"],
			"/static/myhealth/js/weight/list.js").
		OK(c)
}

func (r *WeightHandler) EditPage(c *gin.Context) {
	rr.
		NewPageResponse(
			cc.TotalConstants["Page_Weight_WeightEdit"],
			"/static/myhealth/js/weight/edit.js").
		OK(c)
}

func (r *WeightHandler) CreatePage(c *gin.Context) {
	rr.
		NewPageResponse(
			cc.TotalConstants["Page_Weight_WeightCreate"],
			"/static/myhealth/js/weight/create.js").
		OK(c)
}
