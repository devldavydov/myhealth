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

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) sportSetCommand(userID int64, key, name, comment string) []CmdResponse {
	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetSport(ctx, userID, &storage.Sport{
		Key:     key,
		Name:    name,
		Comment: comment,
	}); err != nil {
		if errors.Is(err, storage.ErrSportInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		r.logger.Error(
			"sport set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) sportSetTemplateCommand(userID int64, key string) []CmdResponse {
	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	sport, err := r.stg.GetSport(ctx, userID, key)
	if err != nil {
		if errors.Is(err, storage.ErrSportNotFound) {
			return NewSingleCmdResponse(MsgErrSportNotFound)
		}

		r.logger.Error(
			"sport get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("s,set,%s,%s,%s", sport.Key, sport.Name, sport.Comment))
}

func (r *CmdProcessor) sportDelCommand(userID int64, key string) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteSport(ctx, userID, key); err != nil {
		if errors.Is(err, storage.ErrSportIsUsed) {
			return NewSingleCmdResponse(MsgErrSportIsUsed)
		}

		r.logger.Error(
			"sport del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) sportListCommand(userID int64) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	sportList, err := r.stg.GetSportList(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"sport list command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	// Build html
	htmlBuilder := html.NewBuilder("Список спорта")

	// Table
	tbl := html.NewTable([]string{"Ключ", "Наименование", "Комментарий"})

	for _, item := range sportList {
		tr := html.NewTr(nil)
		tr.
			AddTd(html.NewTd(html.NewS(item.Key), nil)).
			AddTd(html.NewTd(html.NewS(item.Name), nil)).
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
	sets []int64,
) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetSportActivity(ctx, userID, &storage.SportActivity{
		SportKey:  sportKey,
		Timestamp: storage.NewTimestamp(ts),
		Sets:      sets,
	}); err != nil {
		if errors.Is(err, storage.ErrSportActivityInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		if errors.Is(err, storage.ErrSportNotFound) {
			return NewSingleCmdResponse(MsgErrSportNotFound)
		}

		r.logger.Error(
			"sport activity set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
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

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) sportActivityReportCommand(userID int64, tsFrom, tsTo time.Time) []CmdResponse {
	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	dbRes, err := r.stg.GetSportActivityReport(ctx, userID, storage.NewTimestamp(tsFrom), storage.NewTimestamp(tsTo))
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"sport activity report command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	type grpItem struct {
		sportName string
		sets      string
		total     int64
	}
	grpData := make(map[storage.Timestamp][]grpItem, len(dbRes))

	for _, d := range dbRes {
		total := int64(0)
		sets := make([]string, 0, len(d.Sets))

		for _, item := range d.Sets {
			total += item
			sets = append(sets, strconv.FormatInt(item, 10))
		}

		grpData[d.Timestamp] = append(grpData[d.Timestamp], grpItem{
			sportName: d.SportName,
			sets:      strings.Join(sets, ","),
			total:     total,
		})
	}

	keys := make([]storage.Timestamp, 0, len(grpData))
	for k := range grpData {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Build html
	htmlBuilder := html.NewBuilder("Спортивная активность за период")
	tsFromStr, tsToStr := formatTimestamp(tsFrom), formatTimestamp(tsTo)

	// Table
	tbl := html.NewTable([]string{"Дата", "Спорт", "Подходы", "Итого"})

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
				AddTd(html.NewTd(html.NewS(strconv.Itoa(int(row.total))), nil))

			tbl.AddRow(tr)
		}
	}

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				fmt.Sprintf("Спортивная активность за %s - %s", tsFromStr, tsToStr),
				5,
				html.Attrs{"align": "center"},
			),
			tbl))

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("sport_act_%s_%s.html", tsFromStr, tsToStr),
	})
}
