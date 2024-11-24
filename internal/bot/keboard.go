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
		),
	)
}
