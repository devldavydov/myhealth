package process

import (
	"bytes"

	"github.com/devldavydov/myhealth/internal/cmdproc"
)

type FileType struct {
	Buffer *bytes.Buffer
	Mime   string
	Name   string
}

var _ cmdproc.ITypeAdapter = (*TypeAdapter)(nil)

type TypeAdapter struct{}

func NewTypeAdapter() *TypeAdapter {
	return &TypeAdapter{}
}

func (t *TypeAdapter) File(buf *bytes.Buffer, mime string, fileName string) any {
	return &FileType{
		Buffer: buf,
		Mime:   mime,
		Name:   fileName,
	}
}

func (t *TypeAdapter) OptsHTML() any {
	return ""
}
