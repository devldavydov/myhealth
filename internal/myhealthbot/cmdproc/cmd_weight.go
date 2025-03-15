package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) weightSetCommand(userID int64, ts time.Time, weight float64) []CmdResponse {
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetWeight(ctx,
		userID,
		&storage.Weight{
			Timestamp: storage.NewTimestamp(ts),
			Value:     weight,
		},
	); err != nil {
		if errors.Is(err, storage.ErrWeightInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		r.logger.Error(
			"weight set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) weightDelCommand(userID int64, ts time.Time) []CmdResponse {
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteWeight(ctx,
		userID,
		storage.NewTimestamp(ts),
	); err != nil {
		r.logger.Error(
			"weight del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) weightListCommand(userID int64, tsFrom, tsTo time.Time) []CmdResponse {
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	lst, err := r.stg.GetWeightList(ctx,
		userID,
		storage.NewTimestamp(tsFrom),
		storage.NewTimestamp(tsTo),
	)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"weight list command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	// Report table
	tsFromStr, tsToStr := formatTimestamp(tsFrom), formatTimestamp(tsTo)

	htmlBuilder := html.NewBuilder("Таблица веса")
	accordion := html.NewAccordion("accordionWeight")

	// Table
	tbl := html.NewTable([]string{"Дата", "Вес"})

	xlabels := make([]string, 0, len(lst))
	data := make([]float64, 0, len(lst))
	for _, w := range lst {
		tbl.AddRow(
			html.NewTr(nil).
				AddTd(html.NewTd(html.NewS(formatTimestamp(w.Timestamp.ToTime(r.tz))), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.1f", w.Value)), nil)),
		)
		xlabels = append(xlabels, formatTimestamp(w.Timestamp.ToTime(r.tz)))
		data = append(data, w.Value)
	}

	accordion.AddItem(
		html.HewAccordionItem(
			"tbl",
			fmt.Sprintf("Таблица веса за %s - %s", tsFromStr, tsToStr),
			tbl))

	// Chart
	chart := html.NewCanvas("chart")
	accordion.AddItem(
		html.HewAccordionItem(
			"graph",
			fmt.Sprintf("График веса за %s - %s", tsFromStr, tsToStr),
			chart))

	chartSnip, err := GetChartSnippet(&ChartData{
		ElemID:  "chart",
		XLabels: xlabels,
		Type:    "line",
		Datasets: []ChartDataset{
			{
				Data:  data,
				Label: "Вес",
				Color: ChartColorBlue,
			},
		},
	})
	if err != nil {
		r.logger.Error(
			"weight list command chart error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			accordion,
		),
		html.NewScript(_jsBootstrapURL),
		html.NewScript(_jsChartURL),
		html.NewS(chartSnip),
	)

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("weight_%s_%s.html", tsFromStr, tsToStr),
	})
}
