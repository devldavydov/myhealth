package handlers

import (
	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, cmdProc *cmdproc.CmdProcessor, userID int64) {
	handler := NewHandler(cmdProc, userID)

	router.GET("/", handler.Index)
	router.GET("/file", handler.File)
	router.POST("/api", handler.Api)
	router.NoRoute(handler.NotFound)
}
