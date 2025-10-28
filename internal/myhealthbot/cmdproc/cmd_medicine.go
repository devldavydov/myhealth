package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/devldavydov/myhealth/internal/common/html"
	m "github.com/devldavydov/myhealth/internal/common/messages"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) medSetCommand(userID int64, key, name, unit, comment string) []CmdResponse {
	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetMedicine(ctx, userID, &storage.Medicine{
		Key:     key,
		Name:    name,
		Unit:    unit,
		Comment: comment,
	}); err != nil {
		if errors.Is(err, storage.ErrMedicineInvalid) {
			return NewSingleCmdResponse(m.MsgErrInvalidCommand)
		}

		r.logger.Error(
			"medicine set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) medSetTemplateCommand(userID int64, key string) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	med, err := r.stg.GetMedicine(ctx, userID, key)
	if err != nil {
		if errors.Is(err, storage.ErrMedicineNotFound) {
			return NewSingleCmdResponse(m.MsgErrMedicineNotFound)
		}

		r.logger.Error(
			"medicine get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("m,set,%s,%s,%s,%s", med.Key, med.Name, med.Unit, med.Comment))
}

func (r *CmdProcessor) medDelCommand(userID int64, key string) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteMedicine(ctx, userID, key); err != nil {
		if errors.Is(err, storage.ErrMedicineIsUsed) {
			return NewSingleCmdResponse(m.MsgErrMedicineIsUsed)
		}

		r.logger.Error(
			"medicine del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) medListCommand(userID int64) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	sportList, err := r.stg.GetMedicineList(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(m.MsgErrEmptyResult)
		}

		r.logger.Error(
			"medicine list command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	// Build html
	htmlBuilder := html.NewBuilder("Список медицины")

	// Table
	tbl := html.NewTable([]string{"Ключ", "Наименование", "Единица измерения", "Комментарий"})

	for _, item := range sportList {
		tr := html.NewTr(nil)
		tr.
			AddTd(html.NewTd(html.NewS(item.Key), nil)).
			AddTd(html.NewTd(html.NewS(item.Name), nil)).
			AddTd(html.NewTd(html.NewS(item.Unit), nil)).
			AddTd(html.NewTd(html.NewS(item.Comment), nil))
		tbl.AddRow(tr)
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				"Список медицины",
				5,
				html.Attrs{"align": "center"},
			),
			tbl))

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: "medicine.html",
	})
}

func (r *CmdProcessor) medIndicatorSetCommand(
	userID int64,
	ts time.Time,
	medKey string,
	value float64,
) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetMedicineIndicator(ctx, userID, &storage.MedicineIndicator{
		MedicineKey: medKey,
		Timestamp:   storage.NewTimestamp(ts),
		Value:       value,
	}); err != nil {
		if errors.Is(err, storage.ErrMedicineIndicatorInvalid) {
			return NewSingleCmdResponse(m.MsgErrInvalidCommand)
		}

		if errors.Is(err, storage.ErrMedicineNotFound) {
			return NewSingleCmdResponse(m.MsgErrMedicineNotFound)
		}

		r.logger.Error(
			"medicine indicator set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) medIndicatorDelCommand(userID int64, ts time.Time, medKey string) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteMedicineIndicator(ctx, userID, storage.NewTimestamp(ts), medKey); err != nil {
		r.logger.Error(
			"medicine indicator del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) medIndicatorReportCommand(userID int64, tsFrom, tsTo time.Time) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	dbRes, err := r.stg.GetMedicineIndicatorReport(ctx, userID, storage.NewTimestamp(tsFrom), storage.NewTimestamp(tsTo))
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(m.MsgErrEmptyResult)
		}

		r.logger.Error(
			"medicine indicator report command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	type grpItem struct {
		medName string
		value   float64
	}
	grpData := make(map[storage.Timestamp][]grpItem, len(dbRes))
	graphData := make(map[string]map[storage.Timestamp]float64, len(dbRes))

	for _, d := range dbRes {
		grpData[d.Timestamp] = append(grpData[d.Timestamp], grpItem{
			medName: d.MedicineName,
			value:   d.Value,
		})

		_, ok := graphData[d.MedicineName]
		if !ok {
			graphData[d.MedicineName] = make(map[storage.Timestamp]float64)
		}
		graphData[d.MedicineName][d.Timestamp] = d.Value
	}

	keys := make([]storage.Timestamp, 0, len(grpData))
	for k := range grpData {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Build html
	htmlBuilder := html.NewBuilder("Медицинские показатели за период")
	tsFromStr, tsToStr := formatTimestamp(tsFrom), formatTimestamp(tsTo)
	accordion := html.NewAccordion("accordionMI")

	// Table
	tbl := html.NewTable([]string{"Дата", "Медицина", "Значение"})

	for _, key := range keys {
		first := true
		rows := grpData[key]
		for _, row := range rows {
			tr := html.NewTr(nil)

			if first {
				tr.AddTd(
					html.NewTd(
						html.NewS(formatTimestamp(key.ToTime(r.tz))),
						html.Attrs{"rowspan": strconv.Itoa(len(rows))},
					))
				first = false
			}

			tr.
				AddTd(html.NewTd(html.NewS(row.medName), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", row.value)), nil))

			tbl.AddRow(tr)
		}
	}
	accordion.AddItem(html.HewAccordionItem(
		"tbl",
		"Таблица динамики показателей",
		tbl,
	))

	// Sort medicint names
	medNames := make([]string, 0, len(graphData))
	for k := range graphData {
		medNames = append(medNames, k)
	}
	slices.Sort(medNames)

	// Graph
	var chartSnippets []html.IELement

	for i, medName := range medNames {
		sportTs := make([]storage.Timestamp, 0, len(graphData[medName]))
		for k := range graphData[medName] {
			sportTs = append(sportTs, k)
		}
		slices.Sort(sportTs)

		xlabels := make([]string, 0, len(sportTs))
		data := make([]float64, 0, len(sportTs))

		for _, k := range sportTs {
			xlabels = append(xlabels, formatTimestamp(k.ToTime(r.tz)))
			data = append(data, float64(graphData[medName][k]))
		}

		chartID := fmt.Sprintf("chart%d", i)
		chart := html.NewCanvas(chartID)
		accordion.AddItem(html.HewAccordionItem(
			fmt.Sprintf("sport%d", i),
			fmt.Sprintf("График показателя: %s", medName),
			chart,
		))

		snippet, err := GetChartSnippet(&ChartData{
			PlotFunc: fmt.Sprintf("plot%d", i),
			ElemID:   chartID,
			XLabels:  xlabels,
			Type:     "line",
			Datasets: []ChartDataset{
				{
					Data:  data,
					Label: medName,
					Color: ChartColorBlue,
				},
			},
		})
		if err != nil {
			r.logger.Error(
				"medicine indicator report error",
				zap.Int64("userID", userID),
				zap.Error(err),
			)

			return NewSingleCmdResponse(m.MsgErrInternal)
		}
		chartSnippets = append(chartSnippets, html.NewS(snippet))
	}

	// Doc
	totalElements := []html.IELement{
		html.NewH(
			fmt.Sprintf("Медицинские показатели за %s - %s", tsFromStr, tsToStr),
			5,
			html.Attrs{"align": "center"},
		),
		accordion,
		html.NewScript(_jsBootstrapURL),
		html.NewScript(_jsChartURL),
		html.NewS(GetStartPlotSnippet()),
	}
	totalElements = append(totalElements, chartSnippets...)
	totalElements = append(totalElements, html.NewS(GetEndPlotSnippet()))
	htmlBuilder.
		Add(
			html.NewContainer().Add(totalElements...),
		)

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("medicine_ind_%s_%s.html", tsFromStr, tsToStr),
	})
}
