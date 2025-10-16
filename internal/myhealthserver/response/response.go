package response

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/devldavydov/myhealth/internal/myhealthserver/constants"
	"github.com/gin-gonic/gin"
)

type Response struct {
	meta     gin.H
	tmplName string
}

func (r *Response) WithCustomScript(scripts ...string) *Response {
	var scriptHTML strings.Builder
	for _, s := range scripts {
		scriptHTML.WriteString(fmt.Sprintf(`<script src="%s"></script>`, s))
	}
	r.meta["CustomScripts"] = template.HTML(scriptHTML.String())
	return r
}

func (r *Response) WithCustomStyle(style string) *Response {
	r.meta["CustomStyle"] = template.HTML(style)
	return r
}

func (r *Response) OK(c *gin.Context) {
	c.HTML(http.StatusOK, r.tmplName, r.meta)
}

func NewResponse(title, tmplName string, meta gin.H) *Response {
	totalMeta := meta
	if totalMeta == nil {
		totalMeta = make(gin.H)
	}
	totalMeta["Constants"] = constants.TotalConstants
	totalMeta["Title"] = title

	return &Response{meta: totalMeta, tmplName: tmplName}
}
