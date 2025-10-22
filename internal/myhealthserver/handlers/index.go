package handlers

import (
	cc "github.com/devldavydov/myhealth/internal/myhealthserver/constants"
	rr "github.com/devldavydov/myhealth/internal/myhealthserver/response"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	rr.
		NewResponse(cc.TotalConstants["Page_Main"], "layout.html", nil).
		WithScripts("/static/myhealth/js/index.js").
		OK(c)
}

func NotFound(c *gin.Context) {
	rr.
		NewResponse(cc.TotalConstants["Page_NotFound"], "layout.html", nil).
		WithScripts("/static/myhealth/js/notFound.js").
		OK(c)
}
