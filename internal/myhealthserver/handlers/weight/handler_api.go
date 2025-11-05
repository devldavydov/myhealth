package weight

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	rr "github.com/devldavydov/myhealth/internal/myhealthserver/response"

	"github.com/devldavydov/myhealth/internal/common/messages"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WeightItem struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

func (r *WeightHandler) GetListAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), storage.StorageOperationTimeout)
	defer cancel()

	from, err := strconv.ParseInt(c.Query("from"), 10, 64)
	if err != nil {
		r.logger.Error(
			"weight list 'from' param error",
			zap.String("from", c.Query("from")),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrBadRequest))
		return
	}

	to, err := strconv.ParseInt(c.Query("to"), 10, 64)
	if err != nil {
		r.logger.Error(
			"weight list 'to' param error",
			zap.String("from", c.Query("to")),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrBadRequest))
		return
	}

	wList, err := r.stg.GetWeightList(
		ctx,
		r.userID,
		storage.Timestamp(from),
		storage.Timestamp(to),
		c.Query("order") == "desc",
	)
	if err != nil && !errors.Is(err, storage.ErrEmptyResult) {
		r.logger.Error(
			"weight list api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrInternal))
		return
	}

	data := make([]WeightItem, 0, len(wList))
	for _, w := range wList {
		data = append(data, WeightItem{
			Timestamp: int64(w.Timestamp),
			Value:     w.Value,
		})
	}

	c.JSON(http.StatusOK, rr.NewDataAPIResponse(data))
}

func (r *WeightHandler) GetWeightAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	ts, err := strconv.ParseInt(c.Param("key"), 10, 64)
	if err != nil {
		r.logger.Error(
			"weight 'key' param error",
			zap.String("key", c.Param("key")),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrBadRequest))
		return
	}

	w, err := r.stg.GetWeight(ctx, r.userID, storage.Timestamp(ts))
	if err != nil {
		if errors.Is(err, storage.ErrWeightNotFound) {
			c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrWeightNotFound))
			return
		}

		r.logger.Error(
			"weight get api DB error",
			zap.Int64("userID", r.userID),
			zap.Error(err),
		)

		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, WeightItem{
		Timestamp: int64(w.Timestamp),
		Value:     w.Value,
	})
}

func (r *WeightHandler) SetWeightAPI(c *gin.Context) {
	reqW := &WeightItem{}
	if err := c.Bind(reqW); err != nil {
		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrBadRequest))
		return
	}

	w := &storage.Weight{
		Timestamp: storage.Timestamp(reqW.Timestamp),
		Value:     reqW.Value,
	}

	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetWeight(ctx, r.userID, w); err != nil {
		if errors.Is(err, storage.ErrWeightInvalid) {
			c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrBadRequest))
			return
		}

		r.logger.Error(
			"weight set DB error",
			zap.Int64("userID", r.userID),
			zap.Error(err),
		)

		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, rr.NewOKAPIResponse())
}

func (r *WeightHandler) DeleteWeightAPI(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	ts, err := strconv.ParseInt(c.Param("key"), 10, 64)
	if err != nil {
		r.logger.Error(
			"weight 'key' param error",
			zap.String("key", c.Param("key")),
			zap.Error(err),
		)
		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrBadRequest))
		return
	}

	if err := r.stg.DeleteWeight(ctx, r.userID, storage.Timestamp(ts)); err != nil {
		r.logger.Error(
			"weight del api DB error",
			zap.Int64("userID", r.userID),
			zap.Error(err),
		)

		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrInternal))
		return
	}

	c.JSON(http.StatusOK, rr.NewOKAPIResponse())
}
