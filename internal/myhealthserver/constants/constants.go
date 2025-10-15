package constants

import "maps"

type Constants map[string]string

var TotalConstants Constants

var menus = Constants{
	"Menu_Food":               "Еда",
	"Menu_Food_Journal":       "Журнал приема пищи",
	"Menu_Food_FoodList":      "Управление едой",
	"Menu_Food_Stats":         "Статистика",
	"Menu_Weight":             "Вес",
	"Menu_Weight_WeightList":  "Вес тела",
	"Menu_Weight_Stats":       "Статистика",
	"Menu_Activity":           "Активность",
	"Menu_Activity_Journal":   "Журнал активности",
	"Menu_Activity_SportList": "Управление спортом",
	"Menu_Activity_Stats":     "Статистика",
	"Menu_Finance":            "Финансы",
	"Menu_Finance_BondCalc":   "Расчет доходности облигаций",
	"Menu_Settings":           "Настройки",
	"Menu_Settings_CalcCal":   "Расчет лимита ккал.",
	"Menu_Settings_User":      "Пользователь",
}

var pages = Constants{
	"Page_Main":             "Главная страница",
	"Page_NotFound":         "Страница не найдена",
	"Page_Food_FoodList":    "Управление едой",
	"Page_Settings_CalcCal": "Расчет лимита ккал.",
}

var settings = Constants{
	"Settings_Gender":  "Пол",
	"Settings_GenderM": "Мужской",
	"Settings_GenderF": "Женский",
	"Settings_Weight":  "Вес, кг.",
	"Settings_Height":  "Рост, см.",
	"Settings_Age":     "Возраст, лет",
}

var common = Constants{
	"Common_Calc":    "Рассчитать",
	"Common_ValueG0": "Значение должно быть больше нуля",
}

func init() {
	TotalConstants = make(Constants)
	maps.Copy(TotalConstants, menus)
	maps.Copy(TotalConstants, pages)
	maps.Copy(TotalConstants, settings)
	maps.Copy(TotalConstants, common)
}
