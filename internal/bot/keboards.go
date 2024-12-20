package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type buttons struct {
	ID   string
	Text string
}

var (
	buttonsMap = make(map[string]buttons)
)

func init() {
	buttonsMap["ActualState"] = buttons{
		ID:   "ActualState",
		Text: "Текущие задачи",
	}
	buttonsMap["CleanHistory"] = buttons{
		ID:   "CleanHistory",
		Text: "Очистить историю",
	}
	buttonsMap["RemoveFromQueue"] = buttons{
		ID:   "RemoveFromQueue",
		Text: "Удалить из очереди",
	}
	buttonsMap["ViewQueue"] = buttons{
		ID:   "ViewQueue",
		Text: "Показать очередь",
	}
	buttonsMap["Shutdown"] = buttons{
		ID:   "Shutdown",
		Text: "Выключить сервер",
	}
	buttonsMap["Restart"] = buttons{
		ID:   "Restart",
		Text: "Перезапустить сервер",
	}
	buttonsMap["AutoConnect"] = buttons{
		ID:   "AutoConnect",
		Text: "Автоматическое подключение RAID-массива",
	}
	buttonsMap["Info"] = buttons{
		ID:   "Info",
		Text: "Статистика памяти",
	}
	buttonsMap["Sensors"] = buttons{
		ID:   "Sensors",
		Text: "Показания датчиков",
	}
}

func keyboardDefault() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["CleanHistory"].Text, buttonsMap["CleanHistory"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ViewQueue"].Text, buttonsMap["ViewQueue"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ActualState"].Text, buttonsMap["ActualState"].ID),
		),
	)
}

func keyboardAdmins() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Shutdown"].Text, buttonsMap["Shutdown"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Restart"].Text, buttonsMap["Restart"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["AutoConnect"].Text, buttonsMap["AutoConnect"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["CleanHistory"].Text, buttonsMap["CleanHistory"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["RemoveFromQueue"].Text, buttonsMap["RemoveFromQueue"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Info"].Text, buttonsMap["Info"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ViewQueue"].Text, buttonsMap["ViewQueue"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Sensors"].Text, buttonsMap["Sensors"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ActualState"].Text, buttonsMap["ActualState"].ID),
		),
	)
}
