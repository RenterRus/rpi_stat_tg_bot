package bot

import (
	"context"
	"fmt"
	"log"
	"rpi_stat_tg_bot/internal/cmd"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (k *RealBot) Run() {
	bot, err := tgbotapi.NewBotAPI(k.token)
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	validate := validator.New(validator.WithRequiredStructEnabled())
	u := tgbotapi.NewUpdate(0)
	u.Timeout = k.timeout
	cmd := cmd.NewCMD(k.informer, k.finder)
	updates := bot.GetUpdatesChan(u)
	ctx := context.Background()

	go func() {
		k.downloader.Run(ctx)
	}()

	go func() {
		for k := range k.admins {
			id, err := strconv.Atoi(k)
			if err != nil {
				fmt.Println("ATOI:", err)
			}
			bot.Send(tgbotapi.NewMessage(int64(id), "Бот перезапущен. Через 3 минуты придет информация по обновлению yt-dlp"))
		}

		time.Sleep(time.Minute * 3)
		updInfo := k.downloader.UpdateInfo()
		for k := range k.admins {
			id, err := strconv.Atoi(k)
			if err != nil {
				fmt.Println("ATOI:", err)
			}
			bot.Send(tgbotapi.NewMessage(int64(id), updInfo))
		}
	}()
	for update := range updates {
		// Обработка простых сообщений
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// Убеждаемся, что пользователь из разрешенного пула
			var msg tgbotapi.MessageConfig
			if _, ok := k.allowedIPs[fmt.Sprintf("%d", int(update.Message.Chat.ID))]; ok {
				// Этот блок должен идти до валидации на url, т.к. в очереди, теоретически, может оказаться вообще не ссылка (ручной ввод)
				// Если режим удаления
				if k.isDelete {
					k.isDelete = false
					err := k.queue.DeleteByLink(update.Message.Text)
					if err != nil {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Не удалось удалить из очереди. Причина: %s", err.Error())))
					} else {
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] удалена из очереди", update.Message.Text)))
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, k.downloader.DownloadHistory()))
					}

					continue
				}

				// Не получилось обновружить ссылку
				if err := validate.Var(update.Message.Text, "url"); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, k.welcomeMSG(update.Message.Chat.ID))
					if _, ok := k.admins[fmt.Sprintf("%d", int(update.Message.Chat.ID))]; ok {
						msg.ReplyMarkup = keyboardAdmins()
					} else {
						msg.ReplyMarkup = keyboardDefault()
					}
					// Это ссылка, но вставка не удалась
				} else if err := k.downloader.ToDownload(update.Message.Text); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Не удалось вставить в очередь ссылку %s. Причина: %v", update.Message.Text, err.Error()))
					//Ссылка встала в очередь
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] взята в работу", update.Message.Text))
				}

			} else { // Если нет, то даем ответ о запрещенном доступе
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Доступ запрещен: %d", int(update.Message.Chat.ID)))
			}

			// Отправляем сообщение
			if _, err = bot.Send(msg); err != nil {
				fmt.Println("Send", err)
			}
		} else if update.CallbackQuery != nil { // Если пришел колбэк
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				fmt.Println("update.CallbackQuery", err)
			}

			shutdown := false
			m := ""

			switch update.CallbackQuery.Data {
			case buttonsMap["Shutdown"].ID:
				ctx.Done()
				time.Sleep(time.Second * 10)
				m, shutdown = cmd.Shutdown()
			case buttonsMap["Restart"].ID:
				ctx.Done()
				time.Sleep(time.Second * 10)
				m, shutdown = cmd.Restart()
			case buttonsMap["RemoveFromQueue"].ID:
				m = "Вставьте ссылку, которую надо удалить"
				k.isDelete = true
			case buttonsMap["AutoConnect"].ID:
				m = cmd.Auto()
			case buttonsMap["CleanHistory"].ID:
				m = k.downloader.CleanHistory()
			case buttonsMap["ActualState"].ID:
				m = k.downloader.ActualStatus()
			case buttonsMap["ViewQueue"].ID:
				m = k.downloader.DownloadHistory()
			case buttonsMap["Sensors"].ID:
				m = cmd.Sensors()
			case buttonsMap["Info"].ID:
				command := ""
				m, command = cmd.Info()
				if _, err := bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Команда для быстрого обновления бота на сервере")); err != nil {
					fmt.Println("Info(send)", err)
				}

				if _, err := bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "cd pets/rpi_stat_tg_bot/ && sudo rm main && git pull && sudo systemctl stop runbot.service && go build cmd/main.go && sudo systemctl start runbot.service && sudo systemctl enable runbot.service && sudo systemctl status runbot.service")); err != nil {
					fmt.Println("Info(send2)", err)
				}

				if _, err := bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Команда для подключения RAID массива")); err != nil {
					fmt.Println("Info(send3)", err)
				}
				if _, err := bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, command)); err != nil {
					fmt.Println("Info(send4)", err)
				}
			default:
				m = "Неожиданная команда"
			}

			// Отправляем сообщение, полученное в результате обработки данных выше
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, m)
			if _, err := bot.Send(msg); err != nil {
				fmt.Println("NewMessage", err)
			}

			// Если вызвано выключение или перезапуск - выходим из бесконечного цикла, что б бот корректно завершидл работу
			if shutdown {
				break
			}
		}
	}
}
