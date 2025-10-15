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

type CalcCalPageState struct {
	Gender string
	Weight float64
	Height float64
	Age    int
	Result []CalcCalResultItem
}

type CalcCalResultItem struct {
	Name  string
	Value float64
}

func (r *SettingsHandler) GetCalcCal(c *gin.Context) {
	defaultState := CalcCalPageState{Gender: "m"}

	rr.
		NewResponse(cc.TotalConstants["Page_Settings_CalcCal"],
			"settingsCalcCal.html",
			gin.H{"State": defaultState}).
		OK(c)
}

type CalcCalRequest struct {
	Gender string  `form:"gender"`
	Weight float64 `form:"weight"`
	Height float64 `form:"height"`
	Age    int     `form:"age"`
}

func (r *SettingsHandler) PostCalcCal(c *gin.Context) {
	var req CalcCalRequest
	if err := c.ShouldBind(&req); err != nil {
		// TODO: error on this page
		panic(err)
	}

	// TODO: calc

	newState := CalcCalPageState{
		Gender: req.Gender,
		Weight: req.Weight,
		Height: req.Height,
		Age:    req.Age,
		Result: []CalcCalResultItem{
			{Name: "Foo", Value: 1.1},
			{Name: "Bar", Value: 2.2},
			{Name: "Fuzz", Value: 3.3},
			{Name: "Buzz", Value: 4.4},
		},
	}

	rr.
		NewResponse(cc.TotalConstants["Page_Settings_CalcCal"],
			"settingsCalcCal.html",
			gin.H{"State": newState}).
		OK(c)
}
