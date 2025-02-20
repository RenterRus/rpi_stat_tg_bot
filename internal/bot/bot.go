package bot

import (
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"rpi_stat_tg_bot/internal/downloader"
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RealBot struct {
	informer   informer.Informer
	finder     finder.Finder
	token      string
	timeout    int
	allowedIPs map[string]bool
	admins     map[string]struct{}
	downloader downloader.Downloader
	queueDB    db.Queue
	bot        *tgbotapi.BotAPI
	botName    string
}

type BotConf struct {
	Informer   informer.Informer
	Token      string
	Timeout    int
	Finder     finder.Finder
	AllowedIPs map[string]bool
	Admins     map[string]struct{}
	Downloader downloader.Downloader
	Queue      db.Queue
	Name       string
}

func NewBot(conf BotConf) Bot {
	return &RealBot{
		informer:   conf.Informer,
		token:      conf.Token,
		finder:     conf.Finder,
		timeout:    conf.Timeout,
		allowedIPs: conf.AllowedIPs,
		downloader: conf.Downloader,
		queueDB:    conf.Queue,
		admins:     conf.Admins,
		botName:    conf.Name,
	}
}

func (k *RealBot) welcomeMSG(chatID int64) string {
	welcome := strings.Builder{}
	welcome.WriteString(fmt.Sprintf("Доступ разрешен для: %d", int(chatID)))
	welcome.WriteString("\n")
	welcome.WriteString(k.informer.IPFormatter())
	welcome.WriteString("\n")
	welcome.WriteString("\n")
	welcome.WriteString("Вставьте ссылку для отправки ее в очередь на скачивание")
	welcome.WriteString("\n")
	welcome.WriteString("Или выберите одну из опций ниже")

	return welcome.String()
}
