const ViewsNames = {
    MainPage: "Главная страница",
    FoodList: "Управление едой",
    FoodJournal: "Журнал приема пищи",
    WeightList: "Вес тела",
    ActivityJournal: "Журнал активности",
    SportList: "Управление спортом",
    SettingsCalcCal: "Расчет лимита ккал.",
    SettingsUser: "Пользователь",
    Statistics: "Статистика",
    FinanceBondCalc: "Расчет облигаций"
};

const CommonLabels = {
    ValueG0: "Значение должно быть больше нуля",
    DateGToday: "Дата должна быть больше сегодняшнего дня"
};

const FoodLabels = {};
const WeightLabels = {};
const ActivityLabels = {};

const SettingsLabels = {
    Gender: "Пол",
    GenderMale: "Мужчина",
    GenderFemale: "Женщина",
    Weight: "Вес, кг.",
    Height: "Рост, см.",
    Age: "Возраст, лет.",
    Calculate: "Рассчитать",
    Ubm: "Уровень Базального Метаболизма (УБМ)",
    Activity1: "Сидячая активность",
    Activity2: "Легкая активность",
    Activity3: "Средняя активность",
    Activity4: "Полноценная активность",
    Activity5: "Cупер активность"
};

const FinanceLabels = {
    Nominal: "Номинал, руб.",
    BondCurrentPrice: "Текущая цена облигации, руб.",
    BondMaturityDate: "Дата погашения облигации",
    BondNcd: "НКД, руб.",
    BondCd: "КД, руб.",
    BondCdCount: "Количество выплат КД",
    BondSum: "Сумма приобретения, руб.",
    TotalBoughtCnt: "Количество купленных облигаций",
    TotalBoughtSum: "Сумма купленных облигаций, руб.",
    DaysToMaturity: "Дней до погашения",
    YTM: "Доходность к погашению, %г.",
    TotalSum: "Итоговая сумма, руб.",
    DiffSum: "Разница, руб."
};

export const StringConstants = {
    ...ViewsNames,
    ...CommonLabels,
    ...FoodLabels,
    ...WeightLabels,
    ...ActivityLabels,
    ...FinanceLabels,
    ...SettingsLabels,
};