package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/devldavydov/myhealth/internal/common/html"
	m "github.com/devldavydov/myhealth/internal/common/messages"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) foodSetCommand(
	userID int64,
	key string,
	name string,
	brand string,
	cal100 float64,
	prot100 float64,
	fat100 float64,
	carb100 float64,
	comment string,
) []CmdResponse {
	food := &storage.Food{
		Key:     key,
		Name:    name,
		Brand:   brand,
		Cal100:  cal100,
		Prot100: prot100,
		Fat100:  fat100,
		Carb100: carb100,
		Comment: comment,
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetFood(ctx, userID, food); err != nil {
		if errors.Is(err, storage.ErrFoodInvalid) {
			return NewSingleCmdResponse(m.MsgErrInvalidCommand)
		}

		r.logger.Error(
			"food set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) foodSetTemplateCommand(userID int64, key string) []CmdResponse {
	// Get food from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, userID, key)
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(m.MsgErrFoodNotFound)
		}

		r.logger.Error(
			"food set template command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
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

func (r *CmdProcessor) foodFindCommand(userID int64, pattern string) []CmdResponse {
	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	foodLst, err := r.stg.FindFood(ctx, userID, pattern)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(m.MsgErrEmptyResult)
		}

		r.logger.Error(
			"food find command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
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

func (r *CmdProcessor) foodCalcCommand(userID int64, key string, foodWeight float64) []CmdResponse {
	// Get in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, userID, key)
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(m.MsgErrFoodNotFound)
		}

		r.logger.Error(
			"food calc command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
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
			return NewSingleCmdResponse(m.MsgErrEmptyResult)
		}

		r.logger.Error(
			"food list command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
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

func (r *CmdProcessor) foodDelCommand(userID int64, key string) []CmdResponse {
	// Delete from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteFood(ctx, userID, key); err != nil {
		if errors.Is(err, storage.ErrFoodIsUsed) {
			return NewSingleCmdResponse(m.MsgErrFoodIsUsed)
		}

		r.logger.Error(
			"food del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}
