
var Constants = {
	"Common_Calc": "Рассчитать",
	"Common_DateGToday": "Дата должна быть больше сегодняшнего дня",
	"Common_ValueG0": "Значение должно быть больше нуля",
	"Finance_BondCd": "КД, руб.",
	"Finance_BondCdCount": "Количество выплат КД",
	"Finance_BondCurrentPrice": "Текущая цена, руб.",
	"Finance_BondMaturityDate": "Дата погашения облигации",
	"Finance_BondNcd": "НКД, руб.",
	"Finance_BondSum": "Сумма приобретения, руб.",
	"Finance_DaysToMaturity": "Дней до погашения",
	"Finance_DiffSum": "Разница, руб.",
	"Finance_Nominal": "Номинал, руб.",
	"Finance_TotalBoughtCnt": "Количество купленных облигаций",
	"Finance_TotalBoughtSum": "Сумма купленных облигаций, руб.",
	"Finance_TotalSum": "Итоговая сумма, руб.",
	"Finance_YTM": "Доходность к погашению, %г.",
	"Menu_Activity": "Активность",
	"Menu_Activity_Journal": "Журнал активности",
	"Menu_Activity_SportList": "Управление спортом",
	"Menu_Activity_Stats": "Статистика",
	"Menu_Finance": "Финансы",
	"Menu_Finance_BondCalc": "Расчет доходности облигаций",
	"Menu_Food": "Еда",
	"Menu_Food_FoodList": "Управление едой",
	"Menu_Food_Journal": "Журнал приема пищи",
	"Menu_Food_Stats": "Статистика",
	"Menu_Settings": "Настройки",
	"Menu_Settings_CalcCal": "Расчет лимита ккал.",
	"Menu_Settings_User": "Пользователь",
	"Menu_Weight": "Вес",
	"Menu_Weight_Stats": "Статистика",
	"Menu_Weight_WeightList": "Вес тела",
	"Page_Finance_BondCalc": "Расчет доходности облигаций",
	"Page_Food_FoodList": "Управление едой",
	"Page_Main": "Главная страница",
	"Page_NotFound": "Страница не найдена",
	"Page_Settings_CalcCal": "Расчет лимита ккал.",
	"Settings_Activity1": "Сидячая активность",
	"Settings_Activity2": "Легкая активность",
	"Settings_Activity3": "Средняя активность",
	"Settings_Activity4": "Полноценная активность",
	"Settings_Activity5": "Cупер активность",
	"Settings_Age": "Возраст, лет",
	"Settings_Gender": "Пол",
	"Settings_GenderF": "Женский",
	"Settings_GenderM": "Мужской",
	"Settings_Height": "Рост, см.",
	"Settings_Ubm": "Уровень базального метаболизма, ккал.",
	"Settings_Weight": "Вес, кг.",
};

function hideElement(sel) {
	el = $(sel);
	el.addClass('d-none');
}

function showElement(sel) {
	el = $(sel);
	el.removeClass('d-none');
}
