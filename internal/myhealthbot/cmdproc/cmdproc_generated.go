package cmdproc

// Code generated by "go generate". DO NOT EDIT!

import (
	"fmt"
	"time"
	"strconv"
	"strings"	

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)	

func (r *CmdProcessor) process(c tele.Context, cmd string, userID int64) error {
	cmdParts := []string{}
	for _, part := range strings.Split(cmd, ",") {
		cmdParts = append(cmdParts, strings.Trim(part, " "))
	}

	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.String("command", cmd),
			zap.Int64("userID", userID),
		)
		return c.Send(MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	case "w":
		resp = r.process_w("w", cmdParts[1:], userID)
	case "u":
		resp = r.process_u("u", cmdParts[1:], userID)
	case "f":
		resp = r.process_f("f", cmdParts[1:], userID)
	case "h":
		resp = r._processHelp()
	default:
		r.logger.Error(
			"unknown command",
			zap.String("command", cmd),
			zap.Int64("userID", userID),
		)
		resp = NewSingleCmdResponse(MsgErrInvalidCommand)
	}	

	if r.debugMode {
		if err := c.Send("!!! ОТЛАДОЧНЫЙ РЕЖИМ !!!"); err != nil {
			return err
		}
	}

	for _, rItem := range resp {
		if err := c.Send(rItem.what, rItem.opts...); err != nil {
			return err
		}
	}

	return nil	
}

func (r *CmdProcessor) process_w(baseCmd string, cmdParts []string, userID int64) []CmdResponse {
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
		if len(cmdParts[1:]) != 2 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseTimestamp(r.tz, cmdParts[0])
		if err != nil {
			return argError("Дата")
		}
		
		val1, err := parseFloatG0(cmdParts[1])
		if err != nil {
			return argError("Значение")
		}
		
		resp = r.weightSetCommand(
			userID,
			val0,
			val1,
			)
				
	case "del":
		if len(cmdParts[1:]) != 1 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseTimestamp(r.tz, cmdParts[0])
		if err != nil {
			return argError("Дата")
		}
		
		resp = r.weightDelCommand(
			userID,
			val0,
			)
				
	case "list":
		if len(cmdParts[1:]) != 2 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseTimestamp(r.tz, cmdParts[0])
		if err != nil {
			return argError("С")
		}
		
		val1, err := parseTimestamp(r.tz, cmdParts[1])
		if err != nil {
			return argError("По")
		}
		
		resp = r.weightListCommand(
			userID,
			val0,
			val1,
			)
				
	case "h":
		return NewSingleCmdResponse(
			newCmdHelpBuilder(baseCmd, "Управление весом").
			addCmd(
				"Установка",
				"set",
				"Дата [Дата]",
				"Значение [Дробное>0]",
				).
			addCmd(
				"Удаление",
				"del",
				"Дата [Дата]",
				).
			addCmd(
				"Отчет",
				"list",
				"С [Дата]",
				"По [Дата]",
				).
			build(),
		optsHTML)

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

func (r *CmdProcessor) process_u(baseCmd string, cmdParts []string, userID int64) []CmdResponse {
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
		if len(cmdParts[1:]) != 1 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseFloatG0(cmdParts[0])
		if err != nil {
			return argError("Лимит калорий")
		}
		
		resp = r.userSettingsSetCommand(
			userID,
			val0,
			)
				
	case "st":
		resp = r.userSettingsSetTemplateCommand(userID)
				
	case "get":
		resp = r.userSettingsGetCommand(userID)
				
	case "h":
		return NewSingleCmdResponse(
			newCmdHelpBuilder(baseCmd, "Управление настройками пользователя").
			addCmd(
				"Установка",
				"set",
				"Лимит калорий [Дробное>0]",
				).
			addCmd(
				"Шаблон команды установки",
				"st",
				).
			addCmd(
				"Получение",
				"get",
				).
			build(),
		optsHTML)

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

func (r *CmdProcessor) process_f(baseCmd string, cmdParts []string, userID int64) []CmdResponse {
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
		if len(cmdParts[1:]) != 8 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseStringG0(cmdParts[0])
		if err != nil {
			return argError("Ключ")
		}
		
		val1, err := parseStringG0(cmdParts[1])
		if err != nil {
			return argError("Наименование")
		}
		
		val2, err := parseStringGE0(cmdParts[2])
		if err != nil {
			return argError("Бренд")
		}
		
		val3, err := parseFloatGE0(cmdParts[3])
		if err != nil {
			return argError("ККал 100г")
		}
		
		val4, err := parseFloatGE0(cmdParts[4])
		if err != nil {
			return argError("Б 100г")
		}
		
		val5, err := parseFloatGE0(cmdParts[5])
		if err != nil {
			return argError("Ж 100г")
		}
		
		val6, err := parseFloatGE0(cmdParts[6])
		if err != nil {
			return argError("У 100г")
		}
		
		val7, err := parseStringGE0(cmdParts[7])
		if err != nil {
			return argError("Комментарий")
		}
		
		resp = r.foodSetCommand(
			userID,
			val0,
			val1,
			val2,
			val3,
			val4,
			val5,
			val6,
			val7,
			)
				
	case "st":
		if len(cmdParts[1:]) != 1 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseStringG0(cmdParts[0])
		if err != nil {
			return argError("Ключ")
		}
		
		resp = r.foodSetTemplateCommand(
			userID,
			val0,
			)
				
	case "find":
		if len(cmdParts[1:]) != 1 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseStringGE0(cmdParts[0])
		if err != nil {
			return argError("Подстрока")
		}
		
		resp = r.foodFindCommand(
			userID,
			val0,
			)
				
	case "calc":
		if len(cmdParts[1:]) != 2 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseStringG0(cmdParts[0])
		if err != nil {
			return argError("Ключ")
		}
		
		val1, err := parseFloatGE0(cmdParts[1])
		if err != nil {
			return argError("Вес")
		}
		
		resp = r.foodCalcCommand(
			userID,
			val0,
			val1,
			)
				
	case "list":
		resp = r.foodListCommand(userID)
				
	case "del":
		if len(cmdParts[1:]) != 1 {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		
		val0, err := parseStringG0(cmdParts[0])
		if err != nil {
			return argError("Ключ")
		}
		
		resp = r.foodDelCommand(
			userID,
			val0,
			)
				
	case "h":
		return NewSingleCmdResponse(
			newCmdHelpBuilder(baseCmd, "Управление едой").
			addCmd(
				"Установка",
				"set",
				"Ключ [Строка>0]",
				"Наименование [Строка>0]",
				"Бренд [Строка>=0]",
				"ККал 100г [Дробное>=0]",
				"Б 100г [Дробное>=0]",
				"Ж 100г [Дробное>=0]",
				"У 100г [Дробное>=0]",
				"Комментарий [Строка>=0]",
				).
			addCmd(
				"Шаблон команды установки",
				"st",
				"Ключ [Строка>0]",
				).
			addCmd(
				"Поиск",
				"find",
				"Подстрока [Строка>=0]",
				).
			addCmd(
				"Расчет КБЖУ",
				"calc",
				"Ключ [Строка>0]",
				"Вес [Дробное>=0]",
				).
			addCmd(
				"Список",
				"list",
				).
			addCmd(
				"Удаление",
				"del",
				"Ключ [Строка>0]",
				).
			build(),
		optsHTML)

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

func (r *CmdProcessor) _processHelp() []CmdResponse {
	return nil
}

func parseTimestamp(tz *time.Location, arg string) (time.Time, error) {
	var t time.Time
	var err error

	if arg == "" {
		t = time.Now().In(tz)
	} else {
		t, err = time.Parse("02.01.2006", arg)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tz), nil
}

func parseFloatG0(arg string) (float64, error) {
	val, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, err
	}

	if val <= 0 {
		return 0, fmt.Errorf("not above zero")
	}

	return val, nil
}

func parseFloatGE0(arg string) (float64, error) {
	val, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, err
	}

	if val < 0 {
		return 0, fmt.Errorf("not above or equal zero")
	}

	return val, nil
}

func parseStringG0(arg string) (string, error) {
	if len(arg) == 0 {
		return "", fmt.Errorf("empty string")
	}
	
	return arg, nil
}

func parseStringGE0(arg string) (string, error) {
	return arg, nil
}

func argError(argName string) []CmdResponse {
	return NewSingleCmdResponse(fmt.Sprintf("%s: %s", MsgErrInvalidArg, argName))
}
