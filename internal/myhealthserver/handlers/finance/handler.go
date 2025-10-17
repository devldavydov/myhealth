package finance

import (
	cc "github.com/devldavydov/myhealth/internal/myhealthserver/constants"
	rr "github.com/devldavydov/myhealth/internal/myhealthserver/response"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FinanceHandler struct {
	stg    storage.Storage
	userID int64
	logger *zap.Logger
}

func NewFinanceHandler(stg storage.Storage, userID int64, logger *zap.Logger) *FinanceHandler {
	return &FinanceHandler{stg: stg, userID: userID, logger: logger}
}

func (r *FinanceHandler) BoncCalcPage(c *gin.Context) {
	rr.
		NewResponse(cc.TotalConstants["Page_Finance_BondCalc"],
			"bondCalc.html",
			nil).
		WithCustomScript("/static/myhealth/js/finance/bondcalc.js").
		OK(c)
}
