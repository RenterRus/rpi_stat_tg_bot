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
		Text: "Actual State",
	}
	buttonsMap["CleanHistory"] = buttons{
		ID:   "CleanHistory",
		Text: "Clean History",
	}
	buttonsMap["RemoveFromQueue"] = buttons{
		ID:   "RemoveFromQueue",
		Text: "Remove From Queue",
	}
	buttonsMap["ViewQueue"] = buttons{
		ID:   "ViewQueue",
		Text: "View Queue",
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
		Text: "Auto Connect",
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
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Shutdown"].Text, buttonsMap["Shutdown"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Restart"].Text, buttonsMap["Restart"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["AutoConnect"].Text, buttonsMap["AutoConnect"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Info"].Text, buttonsMap["Info"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["Sensors"].Text, buttonsMap["Sensors"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["CleanHistory"].Text, buttonsMap["CleanHistory"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["RemoveFromQueue"].Text, buttonsMap["RemoveFromQueue"].ID),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ViewQueue"].Text, buttonsMap["ViewQueue"].ID),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["ActualState"].Text, buttonsMap["ActualState"].ID),
		),
	)
}
