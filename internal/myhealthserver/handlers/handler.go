package handlers

import (
	"net/http"

	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cmdProc *cmdproc.CmdProcessor
	userID  int64
}

func NewHandler(cmdProc *cmdproc.CmdProcessor, userID int64) *Handler {
	return &Handler{cmdProc: cmdProc, userID: userID}
}

func (r *Handler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (r *Handler) NotFound(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (r *Handler) Api(c *gin.Context) {

}
