package myhealthbot

import (
	"bytes"

	"github.com/devldavydov/myhealth/internal/cmdproc"
	tele "gopkg.in/telebot.v4"
)

var _ cmdproc.ITypeAdapter = (*BotTypeAdapter)(nil)

type BotTypeAdapter struct{}

func NewBotTypeAdapter() *BotTypeAdapter {
	return &BotTypeAdapter{}
}

func (b *BotTypeAdapter) File(buf *bytes.Buffer, mime string, fileName string) any {
	return &tele.Document{
		File:     tele.FromReader(buf),
		MIME:     mime,
		FileName: fileName,
	}
}

func (b *BotTypeAdapter) OptsHTML() any {
	return &tele.SendOptions{ParseMode: tele.ModeHTML}
}
