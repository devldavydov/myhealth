package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) processJournal(cmdParts []string, userID int64) []CmdResponse {
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
		resp = r.journalSetCommand(cmdParts[1:], userID)
	case "sb":
		resp = r.journalSetBundleCommand(cmdParts[1:], userID)
	case "del":
		resp = r.journalDelCommand(cmdParts[1:], userID)
	case "dm":
		resp = r.journalDelMealCommand(cmdParts[1:], userID)
	case "cp":
		resp = r.journalCopyCommand(cmdParts[1:], userID)
	case "rd":
		resp = r.journalReportDayCommand(cmdParts[1:], userID)
	case "tm":
		resp = r.journalTemplateMealCommand(cmdParts[1:], userID)
	case "fa":
		resp = r.journalFoodAvgWeightCommand(cmdParts[1:], userID)
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

func (r *CmdProcessor) journalSetCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	meal, err := storage.NewMealFromString(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	food_key := cmdParts[2]

	weight, err := strconv.ParseFloat(cmdParts[3], 64)
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetJournal(ctx, userID, &storage.Journal{
		Timestamp:  storage.NewTimestamp(ts),
		Meal:       meal,
		FoodKey:    food_key,
		FoodWeight: weight,
	}); err != nil {
		if errors.Is(err, storage.ErrJournalInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(MsgErrFoodNotFound)
		}

		r.logger.Error(
			"journal set command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalSetBundleCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 3 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	meal, err := storage.NewMealFromString(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	bndlKey := cmdParts[2]

	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetJournalBundle(ctx, userID, storage.NewTimestamp(ts), meal, bndlKey); err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(MsgErrFoodNotFound)
		}

		if errors.Is(err, storage.ErrBundleNotFound) {
			return NewSingleCmdResponse(MsgErrBundleNotFound)
		}

		r.logger.Error(
			"journal set bundle command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalDelCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 3 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	meal, err := storage.NewMealFromString(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	foodKey := cmdParts[2]

	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournal(ctx, userID, storage.NewTimestamp(ts), meal, foodKey); err != nil {
		r.logger.Error(
			"journal del command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalDelMealCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	meal, err := storage.NewMealFromString(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournalMeal(ctx, userID, storage.NewTimestamp(ts), meal); err != nil {
		r.logger.Error(
			"journal del meal command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalCopyCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 4 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse
	tsFrom, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	tsTo, err := r.parseTimestamp(cmdParts[2])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	mealFrom, err := storage.NewMealFromString(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	mealTo, err := storage.NewMealFromString(cmdParts[3])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	cnt, err := r.stg.CopyJournal(ctx,
		userID,
		storage.NewTimestamp(tsFrom),
		mealFrom,
		storage.NewTimestamp(tsTo),
		mealTo)
	if err != nil {
		r.logger.Error(
			"journal copy command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("Скопировано записей: %d", cnt))
}

func (r *CmdProcessor) journalReportDayCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}
	tsStr := formatTimestamp(ts)

	// Get list from DB, user settings
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	var us *storage.UserSettings
	us, err = r.stg.GetUserSettings(ctx, userID)
	if err != nil && !errors.Is(err, storage.ErrUserSettingsNotFound) {
		r.logger.Error(
			"user settings get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	lst, err := r.stg.GetJournalReport(ctx, userID, storage.NewTimestamp(ts), storage.NewTimestamp(ts))
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"journal report day command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	// Report table
	htmlBuilder := html.NewBuilder("Журнал приема пищи")
	tbl := html.NewTable([]string{
		"Наименование", "Вес", "ККал", "Белки", "Жиры", "Углеводы",
	})

	var totalCal, totalProt, totalFat, totalCarb float64
	var subTotalCal, subTotalProt, subTotalFat, subTotalCarb float64
	lastMeal := storage.Meal(-1)
	for i := 0; i < len(lst); i++ {
		j := lst[i]

		// Add meal divider
		if j.Meal != lastMeal {
			tbl.AddRow(
				html.NewTr(html.Attrs{"class": "table-active"}).
					AddTd(html.NewTd(
						html.NewB(j.Meal.MustToString(), nil),
						html.Attrs{"colspan": "6", "align": "center"},
					)),
			)
			lastMeal = j.Meal
		}

		// Add meal rows
		foodLbl := j.FoodName
		if j.FoodBrand != "" {
			foodLbl = fmt.Sprintf("%s - %s", foodLbl, j.FoodBrand)
		}
		foodLbl = fmt.Sprintf("%s [%s]", foodLbl, j.FoodKey)

		tbl.AddRow(
			html.NewTr(nil).
				AddTd(html.NewTd(html.NewS(foodLbl), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.1f", j.FoodWeight)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Cal)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Prot)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Fat)), nil)).
				AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", j.Carb)), nil)))

		totalCal += j.Cal
		totalProt += j.Prot
		totalFat += j.Fat
		totalCarb += j.Carb

		subTotalCal += j.Cal
		subTotalProt += j.Prot
		subTotalFat += j.Fat
		subTotalCarb += j.Carb

		// Add subtotal row
		if i == len(lst)-1 || lst[i+1].Meal != j.Meal {
			tbl.AddRow(
				html.NewTr(nil).
					AddTd(html.NewTd(html.NewB("Всего", nil), html.Attrs{"align": "right", "colspan": "2"})).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalCal)), nil)).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalProt)), nil)).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalFat)), nil)).
					AddTd(html.NewTd(html.NewS(fmt.Sprintf("%.2f", subTotalCarb)), nil)))

			subTotalCal, subTotalProt, subTotalFat, subTotalCarb = 0, 0, 0, 0
		}
	}

	// Footer
	totalPFC := totalProt + totalFat + totalCarb

	tbl.
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего потреблено, ккал: ", nil),
						html.NewS(fmt.Sprintf("%.2f", totalCal)),
					),
					html.Attrs{"colspan": "6"})))

	if us != nil {
		tbl.
			AddFooterElement(html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Лимит, ккал: ", nil),
						html.NewS(fmt.Sprintf("%.2f", us.CalLimit)),
					),
					html.Attrs{"colspan": "6"}))).
			AddFooterElement(html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Разница, ккал: ", nil),
						calDiffSnippet(us.CalLimit-totalCal),
					),
					html.Attrs{"colspan": "6"})))
	}

	tbl.
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего, Б: ", nil),
						pfcSnippet(totalProt, totalPFC),
					),
					html.Attrs{"colspan": "6"}))).
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего, Ж: ", nil),
						pfcSnippet(totalFat, totalPFC),
					),
					html.Attrs{"colspan": "6"}))).
		AddFooterElement(
			html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Всего, У: ", nil),
						pfcSnippet(totalCarb, totalPFC),
					),
					html.Attrs{"colspan": "6"})))

	// Doc
	htmlBuilder.Add(
		html.NewContainer().Add(
			html.NewH(
				fmt.Sprintf("Журнал приема пищи за %s", tsStr),
				5,
				html.Attrs{"align": "center"},
			),
			tbl,
		),
	)

	// Response
	return NewSingleCmdResponse(&tele.Document{
		File:     tele.FromReader(bytes.NewBufferString(htmlBuilder.Build())),
		MIME:     "text/html",
		FileName: fmt.Sprintf("report_%s.html", tsStr),
	})
}

func (r *CmdProcessor) journalTemplateMealCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 2 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Parse
	ts, err := r.parseTimestamp(cmdParts[0])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	meal, err := storage.NewMealFromString(cmdParts[1])
	if err != nil {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	rep, err := r.stg.GetJournalReport(ctx, userID, storage.NewTimestamp(ts), storage.NewTimestamp(ts))
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"journal template meal command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	tsStr := formatTimestamp(ts)
	resp := make([]CmdResponse, 0)

	resp = append(resp, NewCmdResponse("<b>Изменение еды</b>", optsHTML))
	for _, item := range rep {
		if item.Meal != meal {
			continue
		}

		resp = append(resp, NewCmdResponse(
			fmt.Sprintf("j,set,%s,%s,%s,%.1f", tsStr, item.Meal.MustToString(), item.FoodKey, item.FoodWeight),
		))
	}
	resp = append(resp, NewCmdResponse("<b>Удаление еды</b>", optsHTML))
	for _, item := range rep {
		if item.Meal != meal {
			continue
		}

		resp = append(resp, NewCmdResponse(
			fmt.Sprintf("j,del,%s,%s,%s", tsStr, item.Meal.MustToString(), item.FoodKey),
		))
	}

	return resp
}

func (r *CmdProcessor) journalFoodAvgWeightCommand(cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) != 1 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	tsTo := time.Now().In(r.tz)
	tsFrom := tsTo.AddDate(-1, 0, 0)

	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	avgW, err := r.stg.GetJournalFoodAvgWeight(ctx,
		userID,
		storage.NewTimestamp(tsFrom),
		storage.NewTimestamp(tsTo),
		cmdParts[0],
	)
	if err != nil {
		r.logger.Error(
			"journal food avg command DB error",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("Средний вес прима пищи за год: %.1fг.", avgW))
}

func pfcSnippet(val, totalVal float64) html.IELement {
	var s string

	if totalVal == 0 {
		s = fmt.Sprintf("%.2f", val)
	} else {
		s = fmt.Sprintf("%.2f (%.2f%%)", val, val/totalVal*100)
	}

	return html.NewS(s)
}

func calDiffSnippet(diff float64) html.IELement {
	switch {
	case diff < 0 && math.Abs(diff) > 0.01:
		return html.NewSpan(
			html.NewB(fmt.Sprintf("%+.2f", diff), html.Attrs{"class": "text-danger"}),
		)
	case diff >= 0 && math.Abs(diff) > 0.01:
		return html.NewSpan(
			html.NewB(fmt.Sprintf("%+.2f", diff), html.Attrs{"class": "text-success"}),
		)
	default:
		return html.NewS(fmt.Sprintf("%.2f", diff))
	}
}
