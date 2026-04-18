package handlers

import (
	"net/http"
	"time"

	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/devldavydov/myhealth/internal/common/messages"
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

type ApiRequest struct {
	Cmd string `json:"cmd"`
}

type ApiResponse struct {
	IsFile       bool   `json:"isFile"`
	TextResponse string `json:"textResponse"`
	Error        string `json:"error"`
}

func (r *Handler) Api(c *gin.Context) {
	req := ApiRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, &ApiResponse{Error: messages.MsgErrBadRequest})
	}

	loc, _ := time.LoadLocation("Europe/Moscow")
	c.JSON(http.StatusOK, &ApiResponse{TextResponse: time.Now().In(loc).Format("2006-01-02 15:04:05")})
}
