package downloader

import (
	"context"
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func (d *DLP) ytInit(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("===RECOVER===")
			fmt.Println(r)
			fmt.Println("===RECOVER===")

			d.retryInit.Add(1)
			d.ytInit(ctx)
		}
	}()
	ytdlp.MustInstall(ctx, nil)
}

func (d *DLP) Run(ctx context.Context) {
	d.ytInit(ctx)
	res, err := ytdlp.New().Update(context.Background())
	var updateStat strings.Builder
	updateStat.WriteString(fmt.Sprintf("Версия на сервере: %s\n", ytdlp.Version))
	if err != nil {
		updateStat.WriteString(fmt.Sprintf("Определение актульной версии провалено: %s\n", err.Error()))
	}
	if res != nil {
		updateStat.WriteString(fmt.Sprintf("Актуальная версия: %s\n\n", res.Stdout))
		updateStat.WriteString(fmt.Sprintf("Полное сообщение %s\n", res.String()))
	} else {
		updateStat.WriteString("Обновлений не найдено\n")
	}
	updateStat.WriteString(fmt.Sprintf("Количество повторов инициализации загрузчика: %d", d.retryInit.Load()))
	d.retryInit.Store(0)

	d.updateStat = updateStat.String()

	d.dl = ytdlp.New().
		RmCacheDir().
		//ExtractorArgs("youtube:player_client=default,ios").
		SetWorkDir(d.path).
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		Output("%(title)s.%(ext)s").
		NoRestrictFilenames().
		Fixup(ytdlp.FixupForce).
		IgnoreErrors().
		IgnoreNoFormatsError().
		NoAbortOnError().
		CookiesFromBrowser("chrome").
		MarkWatched().
		EmbedChapters().
		ExtractorArgs("youtube:player-client=default,-tv,web_safari,web_embedded").
		ExtractorArgs("youtube:formats=missing_pot")

	d.worker.Actual = make(map[string]map[string]FileInfo)
	doubleWay := make(chan struct{}, d.maxWorkers)

	for {
		select {
		case <-ctx.Done():
			return
		case doubleWay <- struct{}{}:
			go func() {
				defer func() {
					<-doubleWay
				}()

				link, err := d.qdb.SelectOne()
				if err != nil {
					fmt.Printf("\nERROR get link: %v", err)
					return
				}
				if link != "" {

					if err := d.qdb.Update(link, db.StatusWORK, nil); err != nil {
						fmt.Printf("\ndownloader update db error(run): %s\n", err.Error())
					}

					d.downloader(link)
					return
				}
				time.Sleep(time.Second * DEFAULT_TIMEOUT)
			}()

		default:
			time.Sleep(time.Second * DEFAULT_TIMEOUT)
		}
	}
}
