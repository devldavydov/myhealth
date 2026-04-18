package cmdproc

import (
	"fmt"
	"strings"
)

func (r *CmdProcessor) calcCalCalcCommand(userID int64, gender string, weight float64, height float64, age float64) []CmdResponse {
	ubm := 10*weight + 6.25*height - 5*age
	if gender == "m" {
		ubm += 5
	} else {
		ubm -= 161
	}

	var sb strings.Builder
	sb.WriteString("<b>Уровень Базального Метаболизма (УБМ)</b>\n")
	sb.WriteString(fmt.Sprintf("%d ккал\n\n", int64(ubm)))

	sb.WriteString("<b>Усредненные значения по активностям</b>\n")
	for _, i := range []struct {
		name string
		k    float64
	}{
		{name: "\u2022 Сидячая активность", k: 1.2},
		{name: "\u2022 Легкая активность", k: 1.375},
		{name: "\u2022 Средняя активность", k: 1.55},
		{name: "\u2022 Полноценная активность", k: 1.725},
		{name: "\u2022 Супер активность", k: 1.9},
	} {
		sb.WriteString(fmt.Sprintf("<b>%s</b>\n", i.name))
		sb.WriteString(fmt.Sprintf("ККал: %d\n", int64(ubm*i.k)))
		sb.WriteString("\n")
	}

	return NewSingleCmdResponse(sb.String(), optsHTML)
}
