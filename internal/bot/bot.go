package bot

import (
	"fmt"
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
	"strings"
)

type KekBot struct {
	informer   informer.Informer
	finder     finder.Finder
	token      string
	timeout    int
	allowedIPs map[string]struct{}
}

type KekBotConf struct {
	Informer   informer.Informer
	Token      string
	Timeout    int
	Finder     finder.Finder
	AllowedIPs map[string]struct{}
}

func NewKekBot(conf KekBotConf) Bot {
	return &KekBot{
		informer:   conf.Informer,
		token:      conf.Token,
		finder:     conf.Finder,
		timeout:    conf.Timeout,
		allowedIPs: conf.AllowedIPs,
	}
}

func (k *KekBot) welcomeMSG(chatID int64) string {
	welcome := strings.Builder{}
	welcome.WriteString(fmt.Sprintf("Access is allowed for: %d", int(chatID)))
	welcome.WriteString("\n")
	welcome.WriteString("Write /open to open menu-keyboard")
	welcome.WriteString("\n")

	welcome.WriteString(k.informer.IPFormatter())

	return welcome.String()
}
