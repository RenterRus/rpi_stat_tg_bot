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
		Text: "ActualState",
	}
	buttonsMap["CleanHistory"] = buttons{
		ID:   "CleanHistory",
		Text: "CleanHistory",
	}
	buttonsMap["RemoveFromQueue"] = buttons{
		ID:   "RemoveFromQueue",
		Text: "RemoveFromQueue",
	}
	buttonsMap["ViewQueue"] = buttons{
		ID:   "ViewQueue",
		Text: "ViewQueue",
	}
	buttonsMap["Shutdown"] = buttons{
		ID:   "Shutdown",
		Text: "Shutdown",
	}
	buttonsMap["Restart"] = buttons{
		ID:   "Restart",
		Text: "Restart",
	}
	buttonsMap["AutoConnect"] = buttons{
		ID:   "AutoConnect",
		Text: "AutoConnect",
	}
	buttonsMap["Info"] = buttons{
		ID:   "Info",
		Text: "Info",
	}
	buttonsMap["Sensors"] = buttons{
		ID:   "Sensors",
		Text: "Sensors",
	}
}

func keyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Shutdown"].ID, buttonsMap["Shutdown"].Text),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Restart"].ID, buttonsMap["Restart"].Text),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["AutoConnect"].ID, buttonsMap["AutoConnect"].Text),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Info"].ID, buttonsMap["Info"].Text),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Sensors"].ID, buttonsMap["Sensors"].Text),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["CleanHistory"].ID, buttonsMap["CleanHistory"].Text),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["RemoveFromQueue"].ID, buttonsMap["RemoveFromQueue"].Text),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ViewQueue"].ID, buttonsMap["ViewQueue"].Text),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ActualState"].ID, buttonsMap["ActualState"].Text),
		),
	)
}
