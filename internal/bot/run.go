package bot

import (
	"fmt"
	"log"
	"rpi_stat_tg_bot/internal/cmd"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (k *KekBot) Run() {
	bot, err := tgbotapi.NewBotAPI(k.token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = k.timeout

	cmd := cmd.NewCMD(k.informer, k.finder)
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		// Обработка простых сообщений
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Убеждаемся, что пользователь из разрешенного пула
			var msg tgbotapi.MessageConfig
			if _, ok := k.allowedIPs[fmt.Sprintf("%d", int(update.Message.Chat.ID))]; ok {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, k.welcomeMSG(update.Message.Chat.ID))

				switch update.Message.Text {
				case "/open":
					msg.ReplyMarkup = keyboard()
				}
			} else { // Если нет, то даем ответ о запрещенном доступе
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Access is denied: %d", int(update.Message.Chat.ID)))
			}

			// Отправляем сообщение
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil { // Если пришел колбэк
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			shutdown := false
			m := ""

			switch update.CallbackQuery.Data {
			case buttonsMap["Shutdown"].Text:
				m, shutdown = cmd.Shutdown()
			case buttonsMap["Restart"].Text:
				m, shutdown = cmd.Restart()
			case buttonsMap["AutoConnect"].Text:
				m = cmd.Auto()
			case buttonsMap["Info"].Text:
				command := ""
				m, command = cmd.Info()
				if _, err := bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, command)); err != nil {
					panic(err)
				}
			default:
				m = "Press one of button:\nshutdown - shutdown server\nrestart - restart server\nauto - attempt auto connection\ninfo - show info\n"
			}

			// Отправляем сообщение, полученное в результате обработки данных выше
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, m)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}

			// Если вызвано выключение или перезапуск - выходим из бесконечного цикла, что б бот корректно завершидл работу
			if shutdown {
				break
			}
		}
	}
}
