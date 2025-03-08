package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) processFood(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "set":
		resp = r.foodSetCommand(cmdParts[1:], userID)
	case "st":
		resp = r.foodSetTemplateCommand(cmdParts[1:], userID)
	case "find":
		resp = r.foodFindCommand(cmdParts[1:], userID)
	case "calc":
		resp = r.foodCalcCommand(cmdParts[1:], userID)
	case "list":
		resp = r.foodListCommand(userID)
	case "del":
		resp = r.foodDelCommand(cmdParts[1:], userID)
	case "h":
		resp = r.foodHelpCommand()
	default:
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		resp = NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	return resp
}

func (r *CmdProcessor) foodSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 8 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse fields
	food := &storage.Food{
		Key:     cmdParts[0],
		Name:    cmdParts[1],
		Brand:   cmdParts[2],
		Comment: cmdParts[7],
	}

	cal100, err := strconv.ParseFloat(cmdParts[3], 64)
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}
	food.Cal100 = cal100

	prot100, err := strconv.ParseFloat(cmdParts[4], 64)
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}
	food.Prot100 = prot100

	fat100, err := strconv.ParseFloat(cmdParts[5], 64)
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}
	food.Fat100 = fat100

	carb100, err := strconv.ParseFloat(cmdParts[6], 64)
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}
	food.Carb100 = carb100

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetFood(ctx, userID, food); err != nil {
		if errors.Is(err, storage.ErrFoodInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		r.logger.Error(
			"food set command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) foodSetTemplateCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Get food from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, userID, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(MsgErrFoodNotFound)
		}

		r.logger.Error(
			"food set template command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	foodSetTemplate := fmt.Sprintf(
		"f,set,%s,%s,%s,%.2f,%.2f,%.2f,%.2f,%s",
		food.Key,
		food.Name,
		food.Brand,
		food.Cal100,
		food.Prot100,
		food.Fat100,
		food.Carb100,
		food.Comment,
	)
	return NewSingleCmdResponse(foodSetTemplate, optsHTML)
}

func (r *CmdProcessor) foodFindCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	foodLst, err := r.stg.FindFood(ctx, userID, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"food find command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	var sb strings.Builder

	for i, food := range foodLst {
		sb.WriteString(fmt.Sprintf("<b>Ключ:</b> %s\n", food.Key))
		sb.WriteString(fmt.Sprintf("<b>Наименование:</b> %s\n", food.Name))
		sb.WriteString(fmt.Sprintf("<b>Бренд:</b> %s\n", food.Brand))
		sb.WriteString(fmt.Sprintf("<b>ККал100:</b> %.2f\n", food.Cal100))
		sb.WriteString(fmt.Sprintf("<b>Бел100:</b> %.2f\n", food.Prot100))
		sb.WriteString(fmt.Sprintf("<b>Жир100:</b> %.2f\n", food.Fat100))
		sb.WriteString(fmt.Sprintf("<b>Угл100:</b> %.2f\n", food.Carb100))
		sb.WriteString(fmt.Sprintf("<b>Комментарий:</b> %s\n", food.Comment))

		if i != len(foodLst)-1 {
			sb.WriteString("\n")
		}
	}

	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func (r *CmdProcessor) foodCalcCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	foodWeight, err := strconv.ParseFloat(cmdParts[1], 64)
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, userID, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(MsgErrFoodNotFound)
		}

		r.logger.Error(
			"food calc command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>Наименование:</b> %s\n", food.Name))
	sb.WriteString(fmt.Sprintf("<b>Бренд:</b> %s\n", food.Brand))
	sb.WriteString(fmt.Sprintf("<b>Вес:</b> %.1f\n", foodWeight))
	sb.WriteString(fmt.Sprintf("<b>ККал:</b> %.2f\n", foodWeight/100*food.Cal100))
	sb.WriteString(fmt.Sprintf("<b>Бел:</b> %.2f\n", foodWeight/100*food.Prot100))
	sb.WriteString(fmt.Sprintf("<b>Жир:</b> %.2f\n", foodWeight/100*food.Fat100))
	sb.WriteString(fmt.Sprintf("<b>Угл:</b> %.2f\n", foodWeight/100*food.Carb100))

	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func (r *CmdProcessor) foodListCommand(userID int64) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	foodList, err := r.stg.GetFoodList(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"food list command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	// Build html
	htmlBuilder := html.NewBuilder("Список продуктов")

	// Table
	tbl := html.NewTable([]string{
		"Ключ", "Наименование", "Бренд", "ККал в 100г.", "Белки в 100г.",
		"Жиры в 100г.", "Углеводы в 100г.", "Комментарий",
	})

	for _, item := range foodList {
		tr := html.NewTr(nil)
		tr.
			AddTd(html.NewTd(html.NewS(item.Key), nil)).
			AddTd(html.NewTd(html.NewS(item.Name), nil)).
			AddTd(html.NewTd(html.NewS(item.Brand), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Cal100)), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Prot100)), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Fat100)), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", item.Carb100)), nil)).
			AddTd(html.NewTd(html.NewS(item.Comment), nil))
		tbl.AddRow(tr)
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				"Список продуктов и энергетической ценности",
				5,
				html.Attrs{"align": "center"},
			),
			tbl))

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: "food.html",
	})
}

func (r *CmdProcessor) foodDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteFood(ctx, userID, cmdParts[0]); err != nil {
		if errors.Is(err, storage.ErrFoodIsUsed) {
			return NewSingleCmdResponse(MsgErrFoodIsUsed)
		}

		r.logger.Error(
			"food del command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) foodHelpCommand() []CmdResponse {
	return NewSingleCmdResponse(
		newCmdHelpBuilder("Управление едой").
			addCmd(
				"Установка",
				"f,set",
				"Ключ [Строка>0]",
				"Наименование [Строка>0]",
				"Бренд [Строка>=0]",
				"ККал 100г [Дробное>=0]",
				"Б 100г [Дробное>=0]",
				"Ж 100г [Дробное>=0]",
				"У 100г [Дробное>=0]",
				"Комментарий [Строка>=0]",
			).
			addCmd(
				"Шаблон команды установки",
				"f,st",
				"Ключ [Строка>0]",
			).
			addCmd(
				"Поиск",
				"f,find",
				"Подстрока [Строка>=0]",
			).
			addCmd(
				"Расчет КБЖУ",
				"f,calc",
				"Ключ [Строка>0]",
				"Вес [Дробное>=0]",
			).
			addCmd(
				"Список",
				"f,list",
			).
			addCmd(
				"Удаление",
				"f,del",
				"Ключ [Строка>0]",
			).
			build(),
		optsHTML)
}
