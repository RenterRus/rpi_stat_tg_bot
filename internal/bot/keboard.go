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
	buttonsMap["DownloadQueue"] = buttons{
		ID:   "DownloadQueue",
		Text: "DownloadQueue",
	}
	buttonsMap["GetFile"] = buttons{
		ID:   "GetFile",
		Text: "GetFile",
	}
	buttonsMap["SendFile"] = buttons{
		ID:   "SendFile",
		Text: "SendFile",
	}
	buttonsMap["DirecoryView"] = buttons{
		ID:   "DirecoryView",
		Text: "DirecoryView",
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
		// unimplemented
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["DownloadQueue"].ID, buttonsMap["DownloadQueue"].Text),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["DirecoryView"].ID, buttonsMap["DirecoryView"].Text),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["GetFile"].ID, buttonsMap["GetFile"].Text),
			tgbotapi.NewInlineKeyboardButtonData(buttonsMap["SendFile"].ID, buttonsMap["SendFile"].Text),
		),
	)
}
