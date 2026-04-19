package process

import (
	"fmt"

	"github.com/devldavydov/myhealth/internal/cmdproc"
	"github.com/google/uuid"
)

var _ cmdproc.ICmdProcess = (*CmdProcessImpl)(nil)

type Response struct {
	IsFile bool `json:"isFile"`
	//
	TextResponse string `json:"textResponse"`
	//
	FileUUID string `json:"fileUUID"`
	FileMime string `json:"fileMime"`
	FileName string `json:"fileName"`
	//
	Error string `json:"error"`
}

type CmdProcessImpl struct {
	responses []Response
}

func NewCmdProcessImpl() *CmdProcessImpl {
	return &CmdProcessImpl{}
}

func (r *CmdProcessImpl) Send(what any, opts ...any) error {
	switch w := what.(type) {
	case string:
		r.responses = append(r.responses, Response{
			TextResponse: w,
		})
	case *FileType:
		fileUUID := string(uuid.New().String())

		r.responses = append(r.responses, Response{
			IsFile:   true,
			FileUUID: fileUUID,
			FileMime: w.Mime,
			FileName: w.Name,
		})
	default:
		panic(fmt.Sprintf("unknown type [%T] in CmdProcess Send: %v", what, what))
	}

	return nil
}

func (r *CmdProcessImpl) GetResponses() []Response {
	return r.responses
}
