package settings

import (
	cc "github.com/devldavydov/myhealth/internal/myhealthserver/constants"
	rr "github.com/devldavydov/myhealth/internal/myhealthserver/response"
	"github.com/devldavydov/myhealth/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SettingsHandler struct {
	stg    storage.Storage
	userID int64
	logger *zap.Logger
}

func NewSettingsHandler(stg storage.Storage, userID int64, logger *zap.Logger) *SettingsHandler {
	return &SettingsHandler{stg: stg, userID: userID, logger: logger}
}

func (r *SettingsHandler) CalcCalPage(c *gin.Context) {
	rr.
		NewResponse(cc.TotalConstants["Page_Settings_CalcCal"],
			"layout.html",
			nil).
		WithScripts("/static/myhealth/js/settings/calccal.js").
		OK(c)
}
