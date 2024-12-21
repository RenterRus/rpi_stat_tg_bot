package downloader

import (
	"context"
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func (d *DLP) Run(ctx context.Context) {
	ytdlp.MustInstall(context.TODO(), nil)
	res, err := ytdlp.New().Update(context.Background())
	var updateStat strings.Builder
	updateStat.WriteString(fmt.Sprintf("Текущая версия: %s\n", ytdlp.Version))
	if err != nil {
		updateStat.WriteString(fmt.Sprintf("check update error: %s\n", err.Error()))
	}
	if res != nil {
		extinfo, errs := res.GetExtractedInfo()

		updateStat.WriteString(fmt.Sprintf("res.Stdout: %s\n", res.Stdout))
		updateStat.WriteString(fmt.Sprintf("res.GetExtractedInfo(): %v\n", extinfo))
		updateStat.WriteString(fmt.Sprintf("res.GetExtractedInfo() error: %v\n", errs))
		updateStat.WriteString(fmt.Sprintf("res.String(): %s\n", res.String()))
		updateStat.WriteString(fmt.Sprintf("res.Stderr: %s\n", res.Stderr))
		updateStat.WriteString(fmt.Sprintf("res.OutputLogs: %v\n", res.OutputLogs))
		updateStat.WriteString(fmt.Sprintf("res.ExitCode: %d\n", res.ExitCode))
		updateStat.WriteString(fmt.Sprintf("res.Executable: %s\n", res.Executable))
		updateStat.WriteString(fmt.Sprintf("res.Args: %v\n", res.Args))
	} else {
		updateStat.WriteString("Обновлений не найдено\n")
	}

	d.updateStat = updateStat.String()

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
