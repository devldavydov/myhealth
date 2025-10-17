package handlers

import (
	cc "github.com/devldavydov/myhealth/internal/myhealthserver/constants"
	rr "github.com/devldavydov/myhealth/internal/myhealthserver/response"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	rr.
		NewResponse(cc.TotalConstants["Page_Main"], "index.html", nil).
		OK(c)
}

func NotFound(c *gin.Context) {
	rr.
		NewResponse(cc.TotalConstants["Page_NotFound"], "notFound.html", nil).
		OK(c)
}
