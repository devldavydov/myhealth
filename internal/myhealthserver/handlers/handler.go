package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/devldavydov/myhealth/internal/common/messages"
	p "github.com/devldavydov/myhealth/internal/myhealthserver/process"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cmdProcessor    *cmdproc.CmdProcessor
	fileStoragePath string
	userID          int64
}

func NewHandler(cmdProcessor *cmdproc.CmdProcessor, fileStoragePath string, userID int64) *Handler {
	return &Handler{
		cmdProcessor:    cmdProcessor,
		fileStoragePath: fileStoragePath,
		userID:          userID}
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

	prc := p.NewCmdProcessImpl(r.fileStoragePath)

	if err := r.cmdProcessor.Process(prc, req.Cmd, r.userID); err != nil {
		c.JSON(http.StatusOK, &p.Response{Error: err.Error()})
	}

	c.JSON(http.StatusOK, prc.GetResponses())
}

func (r *Handler) File(c *gin.Context) {
	fileUUID, _ := c.GetQuery("fileUUID")
	fileName, _ := c.GetQuery("fileName")
	fileMime, _ := c.GetQuery("fileMime")

	c.Header("Content-Type", fileMime)
	c.FileAttachment(filepath.Join(r.fileStoragePath, fileUUID), fileName)
}
