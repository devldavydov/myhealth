package handlers

import (
	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine, cmdProceccor *cmdproc.CmdProcessor, fileStoragePath string, userID int64) {
	handler := NewHandler(cmdProceccor, fileStoragePath, userID)

	router.GET("/", handler.Index)
	router.GET("/file", handler.File)
	router.POST("/api", handler.Api)
	router.NoRoute(handler.NotFound)
}
