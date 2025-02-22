package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) processSport(cmdParts []string, userID int64) []CmdResponse {
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
	// Sport
	case "set":
		resp = r.sportSetCommand(cmdParts[1:], userID)
	case "st":
		resp = r.sportSetTemplateCommand(cmdParts[1:], userID)
	case "del":
		resp = r.sportDelCommand(cmdParts[1:], userID)
	case "list":
		resp = r.sportListCommand(cmdParts[1:], userID)
	// Sport activity
	case "as":
		resp = r.sportActivitySetCommand(cmdParts[1:], userID)
	case "ad":
		resp = r.sportActivityDelCommand(cmdParts[1:], userID)
	case "ar":
		resp = r.sportActivityReportCommand(cmdParts[1:], userID)
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

func (r *CmdProcessor) sportSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 3 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetSport(ctx, userID, &storage.Sport{
		Key:     cmdParts[0],
		Name:    cmdParts[1],
		Comment: cmdParts[2],
	}); err != nil {
		if errors.Is(err, storage.ErrSportInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		r.logger.Error(
			"sport set command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) sportSetTemplateCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Get from DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	sport, err := r.stg.GetSport(ctx, userID, cmdParts[0])
	if err != nil {
		if errors.Is(err, storage.ErrSportNotFound) {
			return NewSingleCmdResponse(MsgErrSportNotFound)
		}

		r.logger.Error(
			"sport get command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("s,set,%s,%s,%s", sport.Key, sport.Name, sport.Comment))
}

func (r *CmdProcessor) sportDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteSport(ctx, userID, cmdParts[0]); err != nil {
		if errors.Is(err, storage.ErrSportIsUsed) {
			return NewSingleCmdResponse(MsgErrSportIsUsed)
		}

		r.logger.Error(
			"sport del command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) sportListCommand(cmdParts []string, userID int64) []CmdResponse {
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
			zap.Strings("cmdParts", cmdParts),
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

func (r *CmdProcessor) sportActivitySetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) < 3 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	timestamp, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	sets := []int64{}
	for _, part := range cmdParts[2:] {
		s, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			r.logger.Error(
				"invalid command",
				zap.Strings("cmdParts", cmdParts),
				zap.Int64("userID", userID),
			)
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		sets = append(sets, s)
	}

	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetSportActivity(ctx, userID, &storage.SportActivity{
		SportKey:  cmdParts[1],
		Timestamp: storage.NewTimestamp(timestamp),
		Sets:      sets,
	}); err != nil {
		if errors.Is(err, storage.ErrSportNotFound) {
			return NewSingleCmdResponse(MsgErrSportNotFound)
		}

		r.logger.Error(
			"sport activity set command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) sportActivityDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	timestamp, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Call storage
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteSportActivity(ctx, userID, storage.NewTimestamp(timestamp), cmdParts[1]); err != nil {
		r.logger.Error(
			"sport activity del command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) sportActivityReportCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	tsFrom, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	tsTo, err := r.parseTimestamp(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

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
			zap.Strings("cmdParts", cmdParts),
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
