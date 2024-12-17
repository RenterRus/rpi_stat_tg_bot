package bot

import (
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"rpi_stat_tg_bot/internal/downloader"
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
	"strings"
)

type RealBot struct {
	informer   informer.Informer
	finder     finder.Finder
	token      string
	timeout    int
	isDelete   bool
	allowedIPs map[string]struct{}
	downloader downloader.Downloader
	queue      db.Queue
}

type BotConf struct {
	Informer   informer.Informer
	Token      string
	Timeout    int
	Finder     finder.Finder
	AllowedIPs map[string]struct{}
	Downloader downloader.Downloader
	Queue      db.Queue
}

func NewBot(conf BotConf) Bot {
	return &RealBot{
		informer:   conf.Informer,
		token:      conf.Token,
		finder:     conf.Finder,
		timeout:    conf.Timeout,
		allowedIPs: conf.AllowedIPs,
		downloader: conf.Downloader,
		queue:      conf.Queue,
	}
}

func (k *RealBot) welcomeMSG(chatID int64) string {
	welcome := strings.Builder{}
	welcome.WriteString(fmt.Sprintf("Access is allowed for: %d", int(chatID)))
	welcome.WriteString("\n")
	welcome.WriteString("Write /open to open menu-keyboard")
	welcome.WriteString("\n")
	welcome.WriteString("Or paste link for send video to download")
	welcome.WriteString("\n")

	welcome.WriteString(k.informer.IPFormatter())

	return welcome.String()
}
