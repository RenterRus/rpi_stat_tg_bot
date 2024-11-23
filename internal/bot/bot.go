package bot

import (
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
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
