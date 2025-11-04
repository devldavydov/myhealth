package weight

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"

	"github.com/devldavydov/myhealth/internal/common/messages"
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

	wList, err := r.stg.GetWeightList(ctx, r.userID, storage.Timestamp(from), storage.Timestamp(to))
	if err != nil && !errors.Is(err, storage.ErrEmptyResult) {
		r.logger.Error(
			"weight list api DB error",
			zap.Error(err),
		)

		c.JSON(http.StatusOK, rr.NewErrorAPIResponse(messages.MsgErrInternal))
		return
	}

	if c.Query("order") == "desc" {
		slices.Reverse(wList)
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
