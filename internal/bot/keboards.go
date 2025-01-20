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

func (k *RealBot) initKeyboard() {
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
	buttonsMap["Help"] = buttons{
		ID:   "Help",
		Text: "Подсказка",
	}
	buttonsMap["FullState"] = buttons{
		ID:   "FullState",
		Text: "Вся статистика и поддсказки",
	}
	buttonsMap["LinksForUtil"] = buttons{
		ID:   "LinksForUtil",
		Text: "Список ссылок в работе (для утилиты)",
	}
	buttonsMap["EagerMode"] = buttons{
		ID:   "EagerMode",
		Text: "Жадный режим " + k.downloader.EagerModeState(),
	}
}

func (k *RealBot) keyboardDefault() tgbotapi.InlineKeyboardMarkup {
	k.initKeyboard()
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["CleanHistory"].Text, buttonsMap["CleanHistory"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["FullState"].Text, buttonsMap["FullState"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["LinksForUtil"].Text, buttonsMap["LinksForUtil"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ViewQueue"].Text, buttonsMap["ViewQueue"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ActualState"].Text, buttonsMap["ActualState"].ID),
		),
	)
}

func (k *RealBot) keyboardAdmins() tgbotapi.InlineKeyboardMarkup {
	k.initKeyboard()
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Shutdown"].Text, buttonsMap["Shutdown"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Restart"].Text, buttonsMap["Restart"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["AutoConnect"].Text, buttonsMap["AutoConnect"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Help"].Text, buttonsMap["Help"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["EagerMode"].Text, buttonsMap["EagerMode"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["LinksForUtil"].Text, buttonsMap["LinksForUtil"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["FullState"].Text, buttonsMap["FullState"].ID),
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
