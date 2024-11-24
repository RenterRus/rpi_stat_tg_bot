package bot

import (
	"fmt"
	"log"
	"rpi_stat_tg_bot/internal/cmd"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (k *KekBot) Run() {
	bot, err := tgbotapi.NewBotAPI(k.token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = k.timeout

	cmd := cmd.NewCMD(k.informer, k.finder)
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			var msg tgbotapi.MessageConfig
			if _, ok := k.allowedIPs[fmt.Sprintf("%d", int(update.Message.Chat.ID))]; ok {
				welcome := strings.Builder{}
				welcome.WriteString(fmt.Sprintf("Access is allowed for: %d", int(update.Message.Chat.ID)))
				welcome.WriteString("\n")

				welcome.WriteString("write /open to open main menu")
				welcome.WriteString("\n")

				m, cmd, err := k.informer.Basic()
				if err == nil {
					welcome.WriteString(m)
					welcome.WriteString("\n")

					welcome.WriteString(cmd)
					welcome.WriteString("\n")
				}

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, welcome.String())

				switch update.Message.Text {
				case "/open":
					msg.ReplyMarkup = keyboard()
				}
			} else {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Access is denied: %d", int(update.Message.Chat.ID)))
			}
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			shutdown := false
			var msg tgbotapi.MessageConfig
			m := ""
			command := ""

			switch update.CallbackQuery.Data {
			case buttonsMap["Shutdown"].Text:
				m, shutdown = cmd.Shutdown()
			case buttonsMap["Restart"].Text:
				m, shutdown = cmd.Restart()
			case buttonsMap["AutoConnect"].Text:
				m = cmd.Auto()
			case buttonsMap["Info"].Text:
				m, command = cmd.Info()
			default:
				m = "Press one of button:\nshutdown - shutdown server\nrestart - restart server\nauto - attempt auto connection\ninfo - show info\n"
			}

			// And finally, send a message containing the data received.
			msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, m)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}

			if command != "" {
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, command)
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}

			if shutdown {
				break
			}
		}
	}
}
