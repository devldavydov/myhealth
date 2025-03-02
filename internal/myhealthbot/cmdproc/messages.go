package cmdproc

const (
	MsgErrInvalidCommand = "Неправильная команда"
	MsgErrNotImplemented = "В разработке"
	MsgErrInternal       = "Внутренняя ошибка"
	MsgErrEmptyResult    = "Пустой результат"

	MsgErrSportNotFound = "Спорт не найден"
	MsgErrSportIsUsed   = "Спорт используется в активностях"

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
