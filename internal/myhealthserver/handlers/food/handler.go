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

func (r *FoodHandler) GetFoodAPI(c *gin.Context) {
	// Get food from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, r.userID, c.Param("key"))
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrFoodNotFound))
			return
		}

		r.logger.Error(
			"food get api DB error",
			zap.Int64("userID", r.userID),
			zap.Error(err),
		)

		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, FoodItem{
		Key:     food.Key,
		Name:    food.Name,
		Brand:   food.Brand,
		Cal100:  food.Cal100,
		Prot100: food.Prot100,
		Fat100:  food.Fat100,
		Carb100: food.Carb100,
		Comment: food.Comment,
	})
}
