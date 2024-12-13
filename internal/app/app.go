package app

import (
	"rpi_stat_tg_bot/internal/bot"
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
	"sync"
	"time"
)

const time_to_start = 7

type App struct {
	Conf Config
	TTS  int // time to real start in seconds
}

func NewApp(path string) App {
	conf, err := loadConfig(path)
	if err != nil {
		panic(err)
	}

	config := Config{
		Token:          conf.Token,
		Timeout:        conf.Timeout,
		AllowedIDs:     conf.AllowedIDs,
		FTPuser:        conf.FTPuser,
		DevSearch:      conf.DevSearch,
		PathToDownload: conf.PathToDownload,
	}

	return App{
		Conf: config,
		TTS:  time_to_start,
	}
}

func (a *App) Run() {
	finder := finder.NewFinder(finder.KekFinderConf{
		FileSearch: a.Conf.DevSearch,
	})

	allowedIPs := make(map[string]struct{})
	for _, v := range a.Conf.AllowedIDs {
		allowedIPs[v] = struct{}{}
	}

	bot := bot.NewKekBot(bot.KekBotConf{
		Token:      a.Conf.Token,
		Timeout:    a.Conf.Timeout,
		AllowedIPs: allowedIPs,
		Finder:     finder,
		Informer: informer.NewKekInformer(informer.KekInformerConf{
			Finder: finder,
			User:   a.Conf.FTPuser,
		}),
		PathDownload: a.Conf.PathToDownload,
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()

		// Ленивое ожидание, пока устройство загрузится
		time.Sleep(time.Duration(a.TTS) * time.Second)
		bot.Run()
	}()

	wg.Wait()
}
