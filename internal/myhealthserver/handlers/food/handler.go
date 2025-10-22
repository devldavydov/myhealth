package food

import (
	"context"
	"errors"
	"net/http"

	"github.com/devldavydov/myhealth/internal/common/messages"
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
		NewResponse(cc.TotalConstants["Page_Food_FoodList"], "foodList.html", nil).
		WithCustomScript("/static/myhealth/js/food/list.js").
		OK(c)
}

type FoodItem struct {
	Key     string  `json:"key"`
	Name    string  `json:"name"`
	Brand   string  `json:"brand"`
	Cal100  float64 `json:"cal100"`
	Prot100 float64 `json:"prot100"`
	Fat100  float64 `json:"fat100"`
	Carb100 float64 `json:"carb100"`
	Comment string  `json:"comment"`
}

func (r *FoodHandler) GetListAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	foodList, err := r.stg.GetFoodList(ctx, r.userID)
	if err != nil && !errors.Is(err, storage.ErrEmptyResult) {
		r.logger.Error(
			"food list api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrInternal))
		return
	}

	data := make([]FoodItem, 0, len(foodList))
	for _, f := range foodList {
		data = append(data, FoodItem{
			Key:     f.Key,
			Name:    f.Name,
			Brand:   f.Brand,
			Cal100:  f.Cal100,
			Prot100: f.Prot100,
			Fat100:  f.Fat100,
			Carb100: f.Carb100,
			Comment: f.Comment,
		})
	}

	c.JSON(http.StatusOK, rr.NewDataAPIResponse(data))
}
