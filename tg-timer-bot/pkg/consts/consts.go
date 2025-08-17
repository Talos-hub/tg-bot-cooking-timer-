package consts

// for handle user commands
const (
	HELP        = "help"
	MEAT        = "meat"
	EGG         = "egg"
	SETTINGS    = "set"
	START_TIMER = "start_timer"
	START       = "start"
	SHOW        = "show"
)

// TEXTS
const (
	TEXT_HELP = `Привет! Я бот-таймер для приготовления еды.
Доступные команды:
/set - установить таймер
/start_timer - запустить таймер
/show - показать текущие настройки`

	TEXT_SETTINGS = `Вы можете установить пользовательские настройки таймера.
Выберите тип еды:
/meat - для мяса
/egg - для яиц`
	TEXT_MEAT = `Введите время для мяса в формате:
"час минута секунда"
Например: "0 20 0" для 20 минут`
	TEXT_EGG = `Введите время для яиц в формате:
"час минута секунда"
Например: "0 8 0" для 8 минут`
	TEXT_DEFAULt = "Я не понимаю эту команду. Напишите /help для списка команд."
)

// Paths
const (
	JSON_NAME = ".json" // name json file
	//this is paths for default settings
	DEFAULT_MEAT_PATH = "defaultMeat.json" // default path for settings meat
	DEFAULT_EGG_PATH  = "defaultEGG.json"  // default path for settings egg
)
