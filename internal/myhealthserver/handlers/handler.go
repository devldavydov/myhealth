package handlers

import (
	"net/http"

	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/devldavydov/myhealth/internal/common/messages"
	p "github.com/devldavydov/myhealth/internal/myhealthserver/process"
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

type Request struct {
	Cmd string `json:"cmd"`
}

func (r *Handler) Api(c *gin.Context) {
	req := Request{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, &p.Response{Error: messages.MsgErrBadRequest})
	}

	prc := p.NewCmdProcessImpl()

	if err := r.cmdProc.Process(prc, req.Cmd, r.userID); err != nil {
		c.JSON(http.StatusOK, &p.Response{Error: err.Error()})
	}

	c.JSON(http.StatusOK, prc.GetResponses())
}

func (r *Handler) File(c *gin.Context) {
	fileUUID, _ := c.GetQuery("fileUUID")
	fileName, _ := c.GetQuery("fileName")
	fileMime, _ := c.GetQuery("fileMime")

	c.String(http.StatusOK, "Download file [name=%s, mime=%s, uuid=%s]", fileName, fileMime, fileUUID)
}
