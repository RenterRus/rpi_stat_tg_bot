package app

import (
	"rpi_stat_tg_bot/internal/bot"
	"rpi_stat_tg_bot/internal/db"
	"rpi_stat_tg_bot/internal/downloader"
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
	"sync"
	"time"
)

const time_to_start = 7

type App struct {
	Conf Config
	TTS  int
}

func NewApp(path string) App {
	conf, err := loadConfig(path)
	if err != nil {
		panic(err)
	}

	config := *conf

	return App{
		Conf: config,
		TTS:  time_to_start,
	}
}

func (a *App) Run() {
	finder := finder.NewFinder(finder.FinderConf{
		FileSearch: a.Conf.DevSearch,
	})

	queue := db.NewManager(a.Conf.PathToDB)
	bot := bot.NewBot(bot.BotConf{
		DownloadPath: a.Conf.PathToDownload,
		Token:        a.Conf.Token,
		Timeout:      a.Conf.Timeout,
		AllowedIPs: func() map[string]*bot.UserMode {
			allowedIPs := make(map[string]*bot.UserMode)
			for _, v := range a.Conf.AllowedIDs {
				allowedIPs[v] = &bot.UserMode{
					Remove: false,
				}
			}
			return allowedIPs
		}(),
		Admins: func() map[string]struct{} {
			admins := make(map[string]struct{})
			for _, v := range a.Conf.Admins {
				admins[v] = struct{}{}
			}
			return admins
		}(),
		Finder: finder,
		Informer: informer.NewInformer(informer.InformerConf{
			Finder: finder,
			User:   a.Conf.FTPuser,
		}),
		Downloader: downloader.NewDownloader(a.Conf.PathToDownload, queue, a.Conf.MaxWorkers, a.Conf.EagerMode),
		Queue:      queue,
		Name:       a.Conf.ServiceName,
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
