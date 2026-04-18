package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	m "github.com/devldavydov/myhealth/internal/common/messages"

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) sportSetCommand(userID int64, key, name, unit, comment string) []CmdResponse {
	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetSport(ctx, userID, &storage.Sport{
		Key:     key,
		Name:    name,
		Unit:    unit,
		Comment: comment,
	}); err != nil {
		if errors.Is(err, storage.ErrSportInvalid) {
			return NewSingleCmdResponse(m.MsgErrInvalidCommand)
		}

		r.logger.Error(
			"sport set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) sportSetTemplateCommand(userID int64, key string) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	sport, err := r.stg.GetSport(ctx, userID, key)
	if err != nil {
		if errors.Is(err, storage.ErrSportNotFound) {
			return NewSingleCmdResponse(m.MsgErrSportNotFound)
		}

		r.logger.Error(
			"sport get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("s,set,%s,%s,%s,%s", sport.Key, sport.Name, sport.Unit, sport.Comment))
}

func (r *CmdProcessor) sportDelCommand(userID int64, key string) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteSport(ctx, userID, key); err != nil {
		if errors.Is(err, storage.ErrSportIsUsed) {
			return NewSingleCmdResponse(m.MsgErrSportIsUsed)
		}

		r.logger.Error(
			"sport del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) sportListCommand(userID int64) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	sportList, err := r.stg.GetSportList(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(m.MsgErrEmptyResult)
		}

		r.logger.Error(
			"sport list command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	// Build html
	htmlBuilder := html.NewBuilder("Список спорта")

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
				"Список спорта",
				5,
				html.Attrs{"align": "center"},
			),
			tbl))

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: "sport.html",
	})
}

func (r *CmdProcessor) sportActivitySetCommand(
	userID int64,
	ts time.Time,
	sportKey string,
	sets []float64,
	comment string,
) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetSportActivity(ctx, userID, &storage.SportActivity{
		SportKey:  sportKey,
		Timestamp: storage.NewTimestamp(ts),
		Sets:      sets,
		Comment:   comment,
	}); err != nil {
		if errors.Is(err, storage.ErrSportActivityInvalid) {
			return NewSingleCmdResponse(m.MsgErrInvalidCommand)
		}

		if errors.Is(err, storage.ErrSportNotFound) {
			return NewSingleCmdResponse(m.MsgErrSportNotFound)
		}

		r.logger.Error(
			"sport activity set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) sportActivityDelCommand(userID int64, ts time.Time, sportKey string) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteSportActivity(ctx, userID, storage.NewTimestamp(ts), sportKey); err != nil {
		r.logger.Error(
			"sport activity del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	return NewSingleCmdResponse(m.MsgOK)
}

func (r *CmdProcessor) sportActivityReportCommand(userID int64, tsFrom, tsTo time.Time) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	dbRes, err := r.stg.GetSportActivityReport(ctx, userID, storage.NewTimestamp(tsFrom), storage.NewTimestamp(tsTo))
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(m.MsgErrEmptyResult)
		}

		r.logger.Error(
			"sport activity report command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(m.MsgErrInternal)
	}

	type grpItem struct {
		sportName string
		sets      string
		total     float64
		comment   string
	}
	grpData := make(map[storage.Timestamp][]grpItem, len(dbRes))
	graphData := make(map[string]map[storage.Timestamp]float64, len(dbRes))
	totalData := make(map[string]float64, len(dbRes))

	for _, d := range dbRes {
		total := float64(0)
		sets := make([]string, 0, len(d.Sets))

		for _, item := range d.Sets {
			total += item
			sets = append(sets, strconv.FormatFloat(item, 'g', -1, 64))
		}

		grpData[d.Timestamp] = append(grpData[d.Timestamp], grpItem{
			sportName: d.SportName,
			sets:      strings.Join(sets, ","),
			total:     total,
			comment:   d.Comment,
		})

		_, ok := graphData[d.SportName]
		if !ok {
			graphData[d.SportName] = make(map[storage.Timestamp]float64)
		}
		graphData[d.SportName][d.Timestamp] = total

		totalData[d.SportName] = totalData[d.SportName] + total
	}

	keys := make([]storage.Timestamp, 0, len(grpData))
	for k := range grpData {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Build html
	htmlBuilder := html.NewBuilder("Спортивная активность за период")
	tsFromStr, tsToStr := formatTimestamp(tsFrom), formatTimestamp(tsTo)
	accordion := html.NewAccordion("accordionSA")

	// Table
	tbl := html.NewTable([]string{"Дата", "Спорт", "Подходы", "Итого", "Комментарий"})

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
				AddTd(html.NewTd(html.NewS(row.sportName), nil)).
				AddTd(html.NewTd(html.NewS(row.sets), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", row.total)), nil)).
				AddTd(html.NewTd(html.NewS(row.comment), nil))

			tbl.AddRow(tr)
		}
	}
	accordion.AddItem(html.HewAccordionItem(
		"tbl",
		"Таблица активности",
		tbl,
	))

	// Sort sport names
	sportNames := make([]string, 0, len(graphData))
	for k := range graphData {
		sportNames = append(sportNames, k)
	}
	slices.Sort(sportNames)

	// Total table
	tblTotal := html.NewTable([]string{"Спорт", "Итого"})
	for _, sportName := range sportNames {
		tblTotal.AddRow(html.
			NewTr(nil).
			AddTd(html.NewTd(html.NewS(sportName), nil)).
			AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", totalData[sportName])), nil)))
	}
	accordion.AddItem(html.HewAccordionItem(
		"tblTotal",
		"Таблица ИТОГО",
		tblTotal,
	))

	// Graph
	var chartSnippets []html.IELement

	for i, sportName := range sportNames {
		sportTs := make([]storage.Timestamp, 0, len(graphData[sportName]))
		for k := range graphData[sportName] {
			sportTs = append(sportTs, k)
		}
		slices.Sort(sportTs)

		xlabels := make([]string, 0, len(sportTs))
		data := make([]float64, 0, len(sportTs))

		for _, k := range sportTs {
			xlabels = append(xlabels, formatTimestamp(k.ToTime(r.tz)))
			data = append(data, float64(graphData[sportName][k]))
		}

		chartID := fmt.Sprintf("chart%d", i)
		chart := html.NewCanvas(chartID)
		accordion.AddItem(html.HewAccordionItem(
			fmt.Sprintf("sport%d", i),
			fmt.Sprintf("График спорта: %s", sportName),
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
					Label: sportName,
					Color: ChartColorBlue,
				},
			},
		})
		if err != nil {
			r.logger.Error(
				"sport activity report error",
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
			fmt.Sprintf("Спортивная активность за %s - %s", tsFromStr, tsToStr),
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
		FileName: fmt.Sprintf("sport_act_%s_%s.html", tsFromStr, tsToStr),
	})
}
