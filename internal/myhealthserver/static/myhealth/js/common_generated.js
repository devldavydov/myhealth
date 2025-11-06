
var Constants = {
	"Common_Calc": "Рассчитать",
	"Common_Comment": "Комментарий",
	"Common_Date": "Дата",
	"Common_DateGToday": "Дата должна быть больше сегодняшнего дня",
	"Common_DateRange": "Диапазон дат:",
	"Common_Search": "Поиск",
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
	"Food_Brand": "Бренд",
	"Food_C": "К",
	"Food_CPFC": "КБЖУ, 100г.",
	"Food_Cal100": "ККал, 100г.",
	"Food_Carb100": "Углеводы, 100г.",
	"Food_Cb": "У",
	"Food_DeletePrompt": "Удалить еду?",
	"Food_F": "Ж",
	"Food_Fat100": "Жиры, 100г.",
	"Food_Name": "Наименование",
	"Food_NotFound": "Еда не найдена",
	"Food_P": "Б",
	"Food_Prot100": "Белки, 100г.",
	"Food_Saved": "Еда успешно сохранена",
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
	"Page_Food_FoodCreate": "Создание еды",
	"Page_Food_FoodEdit": "Редактирование еды",
	"Page_Food_FoodList": "Управление едой",
	"Page_Main": "Главная страница",
	"Page_NotFound": "Страница не найдена",
	"Page_Settings_CalcCal": "Расчет лимита ккал.",
	"Page_Weight_WeightCreate": "Добавить вес",
	"Page_Weight_WeightEdit": "Редактирование веса тела",
	"Page_Weight_WeightList": "Вес тела",
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
	"Weight_DeletePrompt": "Удалить вес?",
	"Weight_NotFound": "Вес не найден",
	"Weight_Saved": "Вес успешно сохранен",
	"Weight_Value": "Вес, кг.",
};

function createPage(tmpl) {
	$('#page').html(tmpl);
}

function hideElement(sel) {
	el = $(sel);
	el.addClass('d-none');
}

function showElement(sel) {
	el = $(sel);
	el.removeClass('d-none');
}

function showAlert(id, msg) {
	$('#' + id + ' .msg').text(msg);
	showElement('#' + id);
}

function getQueryParams() {
	return new URLSearchParams(window.location.search);
}

function tmplLoader() {
	return `
	<div id="loader" class="spinner-border" role="status">
		<span class="visually-hidden">Loading...</span>
	</div>	
	`;
}

function tmplSearch() {
	return `
	<div class="input-group">
    	<span class="input-group-text"><i class="bi bi-search"></i></span>
    	<input id="search" type="text" class="form-control" placeholder="Поиск">
    	<button class="btn btn-outline-secondary" type="button" id="btnSearchClear"><i class="bi bi-x-lg"></i></button>
	</div>
	`;
}

function tmplToast() {
	return `
	<div class="toast-container position-fixed top-0 end-0 p-3">
	<div id="liveToast" class="toast bg-light" role="alert" aria-live="assertive" aria-atomic="true">
		<div class="d-flex">
		<div id="toastBody" class="toast-body">
		</div>
		<button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
		</div>
	</div>
	</div>	
	`;
}

function tmplAlert(alClass, alID, urlBack) {
	return `
	<div id="${alID}" class="d-none">
		<div class="msg alert ${alClass}" role="alert">
		</div>
		<div class="mb-3">
    		<a href="${urlBack}" class="btn btn-primary"><i class="bi bi-arrow-90deg-left"></i></a>
		</div>
	</div>
	`;
}
