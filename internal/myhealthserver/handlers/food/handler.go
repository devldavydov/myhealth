package food

import (
	"context"
	"errors"
	"net/http"

	"github.com/devldavydov/myhealth/internal/common/messages"
	"github.com/devldavydov/myhealth/internal/myhealthserver/model"
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

func (r *FoodHandler) ListAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	foodList, err := r.stg.GetFoodList(ctx, r.userID)
	if err != nil && !errors.Is(err, storage.ErrEmptyResult) {
		r.logger.Error(
			"food list api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
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

	c.JSON(http.StatusOK, model.NewDataResponse(data))
}

func (r *FoodHandler) GetAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	f, err := r.stg.GetFood(ctx, r.userID, c.Param("key"))
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrFoodNotFound))
			return
		}

		r.logger.Error(
			"food get api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, model.NewDataResponse(FoodItem{
		Key:     f.Key,
		Name:    f.Name,
		Brand:   f.Brand,
		Cal100:  f.Cal100,
		Prot100: f.Prot100,
		Fat100:  f.Fat100,
		Carb100: f.Carb100,
		Comment: f.Comment,
	}))
}
func (r *FoodHandler) DeleteAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteFood(ctx, r.userID, c.Param("key")); err != nil {
		if errors.Is(err, storage.ErrFoodIsUsed) {
			c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrFoodIsUsed))
			return
		}

		r.logger.Error(
			"food del command DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, model.NewOKResponse())
}

func (r *FoodHandler) SetAPI(c *gin.Context) {
	req := FoodItem{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetFood(ctx, r.userID, &storage.Food{
		Key:     req.Key,
		Name:    req.Name,
		Brand:   req.Brand,
		Cal100:  req.Cal100,
		Prot100: req.Prot100,
		Fat100:  req.Fat100,
		Carb100: req.Carb100,
		Comment: req.Comment,
	}); err != nil {
		if errors.Is(err, storage.ErrFoodInvalid) {
			c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrFoodInvalid))
			return
		}

		r.logger.Error(
			"food set api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, model.NewErrorResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, model.NewOKResponse())
}
