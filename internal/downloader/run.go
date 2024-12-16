package downloader

import (
	"context"
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func (d *DLP) Run(ctx context.Context) {
	ytdlp.MustInstall(context.TODO(), nil)

	d.dl = ytdlp.New().
		UnsetCacheDir().
		SetWorkDir(d.path).
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		Output("%(title)s.%(ext)s").
		NoRestrictFilenames().
		Fixup(ytdlp.FixupForce).AbortOnError()

	d.worker.Actual = make(map[string]map[string]FileInfo)
	doubleWay := make(chan struct{}, MAX_THREADS)

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
				} else {

					if err := d.qdb.Update(link, db.StatusWORK); err != nil {
						fmt.Printf("\ndownloader update db error(run): %s\n", err.Error())
					}

					d.downloader(link)
				}
			}()

		default:
			time.Sleep(time.Second * 17)
		}
	}
}
