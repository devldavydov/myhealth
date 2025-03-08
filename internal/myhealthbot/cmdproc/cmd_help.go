package cmdproc

import (
	"fmt"
	"strings"
)

func (r *CmdProcessor) processHelp(userID int64) []CmdResponse {
	var sb strings.Builder
	sb.WriteString("<b>Команды помощи по разделам:</b>\n")
	sb.WriteString("<b>\u2022 w,h</b> - вес\n")
	sb.WriteString("<b>\u2022 f,h</b> - еда\n")
	sb.WriteString("<b>\u2022 j,h</b> - журнал приема пищи\n")
	sb.WriteString("<b>\u2022 b,h</b> - бандлы\n")
	sb.WriteString("<b>\u2022 u,h</b> - настройки пользователя\n")
	sb.WriteString("<b>\u2022 s,h</b> - спорт\n")
	sb.WriteString("<b>\u2022 c,h</b> - расчет калорий\n")
	sb.WriteString("<b>\u2022 m,h</b> - служебный раздел\n\n")
	sb.WriteString("<b>Типы данных:</b>\n")
	sb.WriteString("<b>\u2022 Дата</b> - дата в формате DD.MM.YYYY или пустая строка для текущей даты\n")
	sb.WriteString("<b>\u2022 Строка>0</b> - строка длинной >0\n")
	sb.WriteString("<b>\u2022 Строка>=0</b> - строка длинной >=0\n")
	sb.WriteString("<b>\u2022 Дробное>0</b> - дробное число >0\n")
	sb.WriteString("<b>\u2022 Дробное>=0</b> - дробное число >=0\n")
	sb.WriteString("<b>\u2022 Прием пищи</b> - одно из значений\n")
	sb.WriteString("<b>  \u2022</b> завтрак\n")
	sb.WriteString("<b>  \u2022</b> до обеда\n")
	sb.WriteString("<b>  \u2022</b> обед\n")
	sb.WriteString("<b>  \u2022</b> полдник\n")
	sb.WriteString("<b>  \u2022</b> до ужина\n")
	sb.WriteString("<b>  \u2022</b> ужин\n")
	return NewSingleCmdResponse(sb.String(), optsHTML)
}

type cmdHelpItem struct {
	label string
	cmd   string
	args  []string
}

type cmdHelpBuilder struct {
	label string
	items []cmdHelpItem
}

func newCmdHelpBuilder(label string) *cmdHelpBuilder {
	return &cmdHelpBuilder{label: label}
}

func (r *cmdHelpBuilder) addCmd(label, cmd string, args ...string) *cmdHelpBuilder {
	r.items = append(r.items, cmdHelpItem{
		label: label,
		cmd:   cmd,
		args:  args,
	})
	return r
}

func (r *cmdHelpBuilder) build() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>%s</b>\n", r.label))
	for i, item := range r.items {
		sb.WriteString(fmt.Sprintf("<b>\u2022 %s</b>\n", item.label))
		sb.WriteString(fmt.Sprintf("%s,\n", item.cmd))
		for j, arg := range item.args {
			if j == len(item.args)-1 {
				sb.WriteString(fmt.Sprintf(" %s\n", arg))
			} else {
				sb.WriteString(fmt.Sprintf(" %s,\n", arg))
			}
		}

		if i != len(r.items)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
