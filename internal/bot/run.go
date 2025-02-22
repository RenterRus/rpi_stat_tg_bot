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

const DEFAULT_SLEEP = 7

func (k *RealBot) Run() {
	var err error
	k.bot, err = tgbotapi.NewBotAPI(k.token)
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", k.bot.Self.UserName)

	validate := validator.New(validator.WithRequiredStructEnabled())
	u := tgbotapi.NewUpdate(0)
	u.Timeout = k.timeout
	cmd := cmd.NewCMD(k.informer, k.finder)
	updates := k.bot.GetUpdatesChan(u)
	ctx := context.Background()

	go func() {
		k.downloader.Run(ctx)
	}()

	autoConnect := cmd.Auto()

	go func() {
		k.toAdmins(fmt.Sprintf("Бот запущен. Через минуту придет информация по обновлению yt-dlp.\n%s", autoConnect))
		time.Sleep(time.Minute)
		updInfo := k.downloader.UpdateInfo()
		k.toAdmins(updInfo)
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
				if k.allowedIPs[fmt.Sprintf("%d", update.Message.Chat.ID)].Remove {
					k.allowedIPs[fmt.Sprintf("%d", update.Message.Chat.ID)] = &UserMode{
						Remove: false,
					}
					err := k.queueDB.DeleteByLink(update.Message.Text)
					if err != nil {
						k.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Не удалось удалить из очереди. Причина: %s", err.Error())))
					} else {
						k.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] удалена из очереди", update.Message.Text)))
						if _, ok := k.admins[strconv.Itoa(int(update.Message.Chat.ID))]; !ok {
							k.toAdmins(fmt.Sprintf("Ссылка [%s] удалена из очереди", update.Message.Text))
						}

						k.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, k.downloader.DownloadHistory()))
					}

					continue
				}

				// Не получилось обновружить ссылку
				if err := validate.Var(update.Message.Text, "url"); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, k.welcomeMSG(update.Message.Chat.ID))
					// Это ссылка, но вставка не удалась
				} else if err := k.downloader.ToDownload(update.Message.Text); err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Не удалось вставить в очередь ссылку %s. Причина: %v", update.Message.Text, err.Error()))
					//Ссылка встала в очередь
				} else {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ссылка [%s] взята в работу", update.Message.Text))
					if _, ok := k.admins[strconv.Itoa(int(update.Message.Chat.ID))]; !ok {
						k.toAdmins(fmt.Sprintf("Ссылка [%s] взята в работу", update.Message.Text))
					}
				}

				// Всегда прикрепляем клавиатуру
				if _, ok := k.admins[fmt.Sprintf("%d", int(update.Message.Chat.ID))]; ok {
					msg.ReplyMarkup = k.keyboardAdmins()
				} else {
					msg.ReplyMarkup = k.keyboardDefault()
				}

			} else { // Если нет, то даем ответ о запрещенном доступе
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Доступ запрещен: %d", int(update.Message.Chat.ID)))
			}

			/*if pinned := update.FromChat().PinnedMessage; pinned != nil {
				fmt.Println("pinned")
				if video := pinned.Video; video != nil {
					fmt.Println("video")
					go k.saveVideo(update.Message.Chat.ID, video.Thumbnail.FileID)
				}
			}

			if media := update.Message.Video; media != nil {
				fmt.Println("video")

				file, err := k.bot.GetFile(tgbotapi.FileConfig{
					FileID: media.FileUniqueID,
				})
				if err != nil {
					k.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Не получилось сохранить файл на сервере: %s", err.Error())))
				}

				fmt.Println(file.Link(k.token))

				//		go k.saveVideo(update.Message.Chat.ID, media.Thumbnail.FileID)
			}
			*/

			// Отправляем сообщение
			if _, err = k.bot.Send(msg); err != nil {
				fmt.Println("Send", err)
			}
		} else if update.CallbackQuery != nil { // Если пришел колбэк
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := k.bot.Request(callback); err != nil {
				fmt.Println("update.CallbackQuery", err)
			}

			shutdown := false
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

			mode := func(m MODE) {
				files, err := k.getAllowedFiles()
				if err != nil {
					msg.Text = "Не удается получить список файлов: " + err.Error()
					return
				}
				if len(files) == 0 {
					msg.Text = "Нет файлов пригодных под передачу через телеграмм"
					return
				}

				k.allowedIPs[fmt.Sprintf("%d", int(update.CallbackQuery.Message.Chat.ID))] = &UserMode{
					Mode:     m,
					Download: true,
					Files:    files,
				}

				msg.Text += "Выберите видео для продолжения:\n"
				var rows []tgbotapi.InlineKeyboardButton
				for i, v := range files {
					msg.Text += fmt.Sprintf("%d. %s\n", i, v)

					buttonsMap[v] = buttons{
						ID:         v,
						IsDownload: true,
						Text:       v,
					}
					rows = append(rows, tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(buttonsMap[v].Text, buttonsMap[v].ID),
					)...)

				}

				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardRow(rows...)
			}

			switch update.CallbackQuery.Data {
			case buttonsMap["RestartBot"].ID:
				ctx.Done()
				time.Sleep(time.Second * DEFAULT_SLEEP)
				k.toAdmins("Вызван перезапуск бота")
				msg.Text = cmd.RestartBot(k.botName)
			case buttonsMap["Restart"].ID:
				ctx.Done()
				time.Sleep(time.Second * DEFAULT_SLEEP)
				msg.Text, shutdown = cmd.Restart()
			case buttonsMap["RemoveFromQueue"].ID:
				k.allowedIPs[fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)] = &UserMode{
					Remove: true,
				}
				msg.Text = "Вставьте ссылку, которую надо удалить"
			case buttonsMap["Download"].ID:
				mode(DownloadMode)
			case buttonsMap["EraseDownload"].ID:
				mode(EraseMode)
			case buttonsMap["RemoveDownload"].ID:
				mode(RemoveMode)
			case buttonsMap["AutoConnect"].ID:
				msg.Text = cmd.Auto()
			case buttonsMap["EagerMode"].ID:
				k.downloader.EagerModeToggle()

				msg.Text = "Жадный режим " + k.downloader.EagerModeState()
			case buttonsMap["LinksForUtil"].ID:
				msg.Text = k.queueDB.WorkList()
			case buttonsMap["Help"].ID:
				command := ""
				_, command = cmd.Info()
				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Команда для быстрого обновления бота на сервере")); err != nil {
					fmt.Println("Info(send)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "cd pets/rpi_stat_tg_bot/ && sudo rm main && git pull && sudo systemctl stop runbot.service && go build cmd/main.go && sudo systemctl start runbot.service && sudo systemctl enable runbot.service && sudo systemctl status runbot.service")); err != nil {
					fmt.Println("Info(send2)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Команда для подключения RAID массива")); err != nil {
					fmt.Println("Info(send3)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, command)); err != nil {
					fmt.Println("Info(send4)", err)
				}
			case buttonsMap["FullState"].ID:
				command := ""
				msg.Text, command = cmd.Info()
				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Команда для быстрого обновления бота на сервере")); err != nil {
					fmt.Println("Info(send5)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "cd pets/rpi_stat_tg_bot/ && sudo rm main && git pull && sudo systemctl stop runbot.service && go build cmd/main.go && sudo systemctl start runbot.service && sudo systemctl enable runbot.service && sudo systemctl status runbot.service")); err != nil {
					fmt.Println("Info(send6)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Команда для подключения RAID массива")); err != nil {
					fmt.Println("Info(send7)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, command)); err != nil {
					fmt.Println("Info(send8)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, k.downloader.DownloadHistory())); err != nil {
					fmt.Println("Info(send9)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, k.downloader.ActualStatus())); err != nil {
					fmt.Println("Info(send10)", err)
				}

				if _, err := k.bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, cmd.Sensors())); err != nil {
					fmt.Println("Info(send11)", err)
				}
			case buttonsMap["CleanHistory"].ID:
				msg.Text = k.downloader.CleanHistory()
			case buttonsMap["Quality"].ID:
				k.downloader.QualityModeToggle()

				msg.Text = "Максимальное качество видео: " + k.downloader.QualityModeState()
			case buttonsMap["ActualState"].ID:
				msg.Text = k.downloader.ActualStatus()
			case buttonsMap["ViewQueue"].ID:
				msg.Text = k.downloader.DownloadHistory()
			case buttonsMap["Sensors"].ID:
				msg.Text = cmd.Sensors()
			case buttonsMap["Info"].ID:
				msg.Text, _ = cmd.Info()
			default:
				if v, ok := buttonsMap[update.CallbackQuery.Data]; ok && v.IsDownload {
					k.allowedIPs[fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)] = &UserMode{
						Download: false,
					}

					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Видео [%s] поставлено в загрузку", v.Text))
					go func() {
						if k.allowedIPs[fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)].Mode == DownloadMode || k.allowedIPs[fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)].Mode == EraseMode {
							k.loadVideo(update.CallbackQuery.Message.Chat.ID, v.Text)
						}
						if k.allowedIPs[fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)].Mode == EraseMode || k.allowedIPs[fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)].Mode == RemoveMode {
							k.removeVideo(update.CallbackQuery.Message.Chat.ID, v.Text)
						}
					}()
				}
				msg.Text = "Неожиданная команда"
			}

			// Отправляем сообщение, полученное в результате обработки данных выше

			if _, err := k.bot.Send(msg); err != nil {
				fmt.Println("NewMessage", err)
			}

			// Если вызвано выключение или перезапуск - выходим из бесконечного цикла, что б бот корректно завершидл работу
			if shutdown {
				break
			}
		}
	}
}
