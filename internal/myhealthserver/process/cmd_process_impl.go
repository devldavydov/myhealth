package process

import (
	"fmt"
	"os"
	"path/filepath"

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
	responses       []Response
	fileStoragePath string
}

func NewCmdProcessImpl(fileStoragePath string) *CmdProcessImpl {
	return &CmdProcessImpl{fileStoragePath: fileStoragePath}
}

func (r *CmdProcessImpl) Send(what any, opts ...any) error {
	switch w := what.(type) {
	case string:
		r.responses = append(r.responses, Response{
			TextResponse: w,
		})
	case *FileType:
		fileUUID := string(uuid.New().String())

		if err := os.WriteFile(
			filepath.Join(r.fileStoragePath, fileUUID),
			w.Buffer.Bytes(),
			0644); err != nil {
			r.responses = append(r.responses, Response{Error: err.Error()})
		} else {
			r.responses = append(r.responses, Response{
				IsFile:   true,
				FileUUID: fileUUID,
				FileMime: w.Mime,
				FileName: w.Name,
			})
		}
	default:
		panic(fmt.Sprintf("unknown type [%T] in CmdProcess Send: %v", what, what))
	}

	return nil
}

func (r *CmdProcessImpl) GetResponses() []Response {
	return r.responses
}
