package cmdproc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/devldavydov/myhealth/internal/common/html"
	"github.com/devldavydov/myhealth/internal/storage"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)

func (r *CmdProcessor) journalSetCommand(
	userID int64,
	ts time.Time,
	meal storage.Meal,
	foodKey string,
	foodWeight float64,
) []CmdResponse {
	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetJournal(ctx, userID, &storage.Journal{
		Timestamp:  storage.NewTimestamp(ts),
		Meal:       meal,
		FoodKey:    foodKey,
		FoodWeight: foodWeight,
	}); err != nil {
		if errors.Is(err, storage.ErrJournalInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(MsgErrFoodNotFound)
		}

		r.logger.Error(
			"journal set command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalSetBundleCommand(
	userID int64,
	ts time.Time,
	meal storage.Meal,
	bndlKey string,
) []CmdResponse {
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
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalDelCommand(
	userID int64,
	ts time.Time,
	meal storage.Meal,
	foodKey string,
) []CmdResponse {
	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournal(ctx, userID, storage.NewTimestamp(ts), meal, foodKey); err != nil {
		r.logger.Error(
			"journal del command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalDelMealCommand(userID int64, ts time.Time, meal storage.Meal) []CmdResponse {
	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteJournalMeal(ctx, userID, storage.NewTimestamp(ts), meal); err != nil {
		r.logger.Error(
			"journal del meal command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalDelBundleCommand(
	userID int64,
	ts time.Time,
	meal storage.Meal,
	bndlKey string,
) []CmdResponse {
	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DelJournalBundle(ctx, userID, storage.NewTimestamp(ts), meal, bndlKey); err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(MsgErrFoodNotFound)
		}

		if errors.Is(err, storage.ErrBundleNotFound) {
			return NewSingleCmdResponse(MsgErrBundleNotFound)
		}

		r.logger.Error(
			"journal del bundle command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalCopyCommand(
	userID int64,
	tsFrom time.Time,
	mealFrom storage.Meal,
	tsTo time.Time,
	mealTo storage.Meal,
) []CmdResponse {
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
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(fmt.Sprintf("Скопировано записей: %d", cnt))
}

func (r *CmdProcessor) journalReportDayCommand(userID int64, ts time.Time) []CmdResponse {
	tsStr := formatTimestamp(ts)

	// Get list from DB, user settings
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	// Get total burned calories
	var us *storage.UserSettings
	us, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil && !errors.Is(err, storage.ErrUserSettingsNotFound) {
		r.logger.Error(
			"user settings get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	burnedCal, err := r.stg.GetTotalBurnedCal(ctx, userID, storage.NewTimestamp(ts))
	if err != nil && !errors.Is(err, storage.ErrTotalBurnedCalNotFound) {
		r.logger.Error(
			"total burned cal get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	var totalBurnedCal float64
	if burnedCal != 0 {
		totalBurnedCal = burnedCal
	} else if us != nil {
		totalBurnedCal = us.CalLimit
	}

	// Generate report
	lst, err := r.stg.GetJournalReport(ctx, userID, storage.NewTimestamp(ts), storage.NewTimestamp(ts))
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"journal report day command DB error",
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

	if totalBurnedCal != 0 {
		tbl.
			AddFooterElement(html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Потрачено, ккал: ", nil),
						html.NewS(fmt.Sprintf("%.2f", totalBurnedCal)),
					),
					html.Attrs{"colspan": "6"}))).
			AddFooterElement(html.NewTr(nil).
				AddTd(html.NewTd(
					html.NewSpan(
						html.NewB("Разница, ккал: ", nil),
						calDiffSnippet(totalBurnedCal-totalCal),
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

func (r *CmdProcessor) journalReportDayCalloriesCommand(userID int64, ts time.Time) []CmdResponse {
	// Get list from DB, user settings
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	// Get total burned calories
	var us *storage.UserSettings
	us, err := r.stg.GetUserSettings(ctx, userID)
	if err != nil && !errors.Is(err, storage.ErrUserSettingsNotFound) {
		r.logger.Error(
			"user settings get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	burnedCal, err := r.stg.GetTotalBurnedCal(ctx, userID, storage.NewTimestamp(ts))
	if err != nil && !errors.Is(err, storage.ErrTotalBurnedCalNotFound) {
		r.logger.Error(
			"total burned cal get command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	var totalBurnedCal float64
	if burnedCal != 0 {
		totalBurnedCal = burnedCal
	} else if us != nil {
		totalBurnedCal = us.CalLimit
	}

	// Generate report
	lst, err := r.stg.GetJournalReport(ctx, userID, storage.NewTimestamp(ts), storage.NewTimestamp(ts))
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"journal report day calories command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	var sb strings.Builder
	sb.WriteString("<b>Отчет по ккал за день:</b>\n\n")

	mealCal := map[string]float64{}
	var mealOrder []string
	var totalCal float64
	for _, j := range lst {
		mealStr := j.Meal.MustToString()
		_, ok := mealCal[mealStr]
		if !ok {
			mealOrder = append(mealOrder, mealStr)
		}
		mealCal[mealStr] += j.Cal
		totalCal += j.Cal
	}

	for _, mealStr := range mealOrder {
		sb.WriteString(fmt.Sprintf("%s, ккал: %.2f\n", mealStr, mealCal[mealStr]))
	}

	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Всего потреблено, ккал: %.2f\n", totalCal))
	if totalBurnedCal != 0 {
		sb.WriteString(fmt.Sprintf("Потрачено, ккал: %.2f\n", totalBurnedCal))
		sb.WriteString(fmt.Sprintf("Разница, ккал: <b>%+.2f</b>\n", totalBurnedCal-totalCal))
	}

	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func (r *CmdProcessor) journalTemplateMealCommand(userID int64, ts time.Time, meal storage.Meal) []CmdResponse {
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

func (r *CmdProcessor) journalFoodStatCommand(userID int64, foodKey string) []CmdResponse {
	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	food, err := r.stg.GetFood(ctx, userID, foodKey)
	if err != nil {
		if errors.Is(err, storage.ErrFoodNotFound) {
			return NewSingleCmdResponse(MsgErrFoodNotFound)
		}

		r.logger.Error(
			"journal food stat command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInternal)
	}

	foodStat, err := r.stg.GetJournalFoodStat(ctx,
		userID,
		foodKey,
	)
	if err != nil {
		if errors.Is(err, storage.ErrEmptyResult) {
			return NewSingleCmdResponse(MsgErrEmptyResult)
		}

		r.logger.Error(
			"journal food stat command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)
		return NewSingleCmdResponse(MsgErrInternal)
	}

	var sb strings.Builder
	foodName := food.Name
	if food.Brand != "" {
		foodName = fmt.Sprintf("%s [%s]", foodName, food.Brand)
	}
	sb.WriteString(fmt.Sprintf("<b>Наименование:</b> %s\n", foodName))
	sb.WriteString(fmt.Sprintf("<b>Итого съедено:</b> %.1fг. (%.1fкг.)\n", foodStat.TotalWeight, foodStat.TotalWeight/1000))
	sb.WriteString(fmt.Sprintf("<b>Средний вес за приём пищи:</b> %.1fг.\n", foodStat.AvgWeight))
	sb.WriteString(fmt.Sprintf("<b>Количество раз:</b> %d\n", foodStat.TotalCount))
	sb.WriteString(fmt.Sprintf("<b>Первый раз:</b> %s\n", formatTimestamp(foodStat.FirstTimestamp.ToTime(r.tz))))
	sb.WriteString(fmt.Sprintf("<b>Последний раз:</b> %s\n", formatTimestamp(foodStat.LastTimestamp.ToTime(r.tz))))

	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func (r *CmdProcessor) journalSetDayTotalCal(userID int64, ts time.Time, totalCal float64) []CmdResponse {
	// Save in DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.SetTotalBurnedCal(ctx, userID, storage.NewTimestamp(ts), totalCal); err != nil {
		if errors.Is(err, storage.ErrDayTotalCalInvalid) {
			return NewSingleCmdResponse(MsgErrInvalidCommand)
		}

		r.logger.Error(
			"journal set total burned cal command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
}

func (r *CmdProcessor) journalDeleteDayTotalCal(userID int64, ts time.Time) []CmdResponse {
	// Call DB
	ctx, cancel := context.WithTimeout(context.Background(), storage.StorageOperationTimeout)
	defer cancel()

	if err := r.stg.DeleteTotalBurnedCal(ctx, userID, storage.NewTimestamp(ts)); err != nil {
		r.logger.Error(
			"journal del total burned cal command DB error",
			zap.Int64("userID", userID),
			zap.Error(err),
		)

		return NewSingleCmdResponse(MsgErrInternal)
	}

	return NewSingleCmdResponse(MsgOK)
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
