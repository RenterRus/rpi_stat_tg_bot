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
	res, err := ytdlp.New().Update(context.Background())
	if err != nil {
		fmt.Println("check update error", err)
	}
	if res != nil {
		extinfo, errs := res.GetExtractedInfo()
		fmt.Println("res.Stdout", res.Stdout)
		fmt.Println("res.GetExtractedInfo()", extinfo)
		fmt.Println("res.GetExtractedInfo() error", errs)
		fmt.Println("res.String()", res.String())
		fmt.Println("res.Stderr", res.Stderr)
		fmt.Println("res.OutputLogs", res.OutputLogs)
		fmt.Println("res.ExitCode", res.ExitCode)
		fmt.Println("res.Executable", res.Executable)
		fmt.Println("res.Args", res.Args)
	}

	d.dl = ytdlp.New().
		UnsetCacheDir().
		SetWorkDir(d.path).
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		Output("%(title)s.%(ext)s").
		NoRestrictFilenames().
		Fixup(ytdlp.FixupForce).AbortOnError()

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
				} else {

					if err := d.qdb.Update(link, db.StatusWORK, nil); err != nil {
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
