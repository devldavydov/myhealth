package messages

const (
	MsgErrInvalidCommand   = "Неправильная команда (h для помощи)"
	MsgErrInvalidArgsCount = "Неверное количество аргументов"
	MsgErrInvalidArg       = "Неверный аргумент"

	MsgErrInternal    = "Внутренняя ошибка"
	MsgErrEmptyResult = "Пустой результат"

	MsgErrSportNotFound = "Спорт не найден"
	MsgErrSportIsUsed   = "Спорт используется в активностях"

	MsgErrMedicineNotFound = "Медицина не найдена"
	MsgErrMedicineIsUsed   = "Медицина используется в показателях"

	MsgErrFoodNotFound = "Еда не найдена"
	MsgErrFoodIsUsed   = "Еда уже используется в журнале приема пищи или бандле"

	MsgErrBundleDepBundleNotFound  = "Зависимый бандл не найден в базе данных"
	MsgErrBundleDepFoodNotFound    = "Зависимая еда не найдена в базе данных"
	MsgErrBundleDepBundleRecursive = "Зависимый бандл не может быть рекурсивным"
	MsgErrBundleNotFound           = "Бандл не найден в базе данных"
	MsgErrBundleIsUsed             = "Бандл уже используется в другом бандле"

	MsgErrUserSettingsNotFound = "Настройки пользователя не найдены"

	MsgOK = "OK"
)
