package downloader

import (
	"context"
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

const DEFAULT_TIMEOUT = 17

func (d *DLP) ToDownload(url string) error {
	if err := d.qdb.Insert(url); err != nil {
		return fmt.Errorf("to download: %w", err)
	}
	return nil
}

func (d *DLP) downloader(link string) {
	defer func() {
		delete(d.worker.Actual, link)
	}()

	progressInfo := map[string]FileInfo{}
	name := ""
	duration := float64(DEFAULT_TIMEOUT)
	d.dl.ProgressFunc(time.Duration(time.Millisecond*750), func(update ytdlp.ProgressUpdate) {
		size := (float64(update.DownloadedBytes) / 1024) / 1024 // К мегабайтам
		totalSize := (float64(update.TotalBytes) / 1024) / 1024 // К мегабайтам
		fmt.Println(update.Status, update.PercentString(), fmt.Sprintf("[%.2f/%.2f]mb", size, totalSize), update.Filename)
		status := string(update.Status)
		if strings.Contains(status, "finished") {
			status = "converting"
		}

		if update.Info.Duration != nil {
			duration = *update.Info.Duration
		} else {
			duration = DEFAULT_TIMEOUT * 60 // 17 минут
		}

		progressInfo[update.Filename] = FileInfo{
			Name:         d.path + "/" + update.Filename,
			DownloadSize: fmt.Sprintf("%.2f", size),
			TotalSize:    fmt.Sprintf("%.2f", totalSize),
			Proc:         update.PercentString(),
			Status:       status,
		}
		d.Lock()
		d.worker.Actual[link] = progressInfo
		d.Unlock()
		if name != update.Filename {
			name = update.Filename
			if err := d.qdb.Update(link, db.StatusWORK, &name); err != nil {
				fmt.Println("Update name (into work)", err)
			}
		}
	})

	_, err := d.dl.Run(context.TODO(), link)

	baseMessage := FileInfo{
		Name:         progressInfo[name].Name,
		DownloadSize: progressInfo[name].DownloadSize,
		TotalSize:    progressInfo[name].TotalSize,
		Proc:         "100%",
		Status:       "done",
	}
	if err != nil {
		baseMessage.Status = fmt.Sprintf("error: [%s]", err.Error())
		d.Lock()
		d.worker.Actual[link][name] = baseMessage
		d.Unlock()

		fmt.Printf("\ndownload error: %s\n", err.Error())

		d.totalRetry.Add(1)
		if err := d.qdb.Update(link, db.StatusNEW, &name); err != nil {
			fmt.Printf("\ndownloader set video to queue (retry): [%s]\n", err.Error())
		}
		baseMessage.Status += "\n- - - - - - -Returned to the queue."
		d.Lock()
		d.worker.Actual[link][name] = baseMessage
		d.Unlock()

		if !d.eagerMode {

			baseMessage.Status += fmt.Sprintf("\n- - - - - - - -EagleMode: %s\n- - - - - - - - -Waiting %d second to next", d.EagerModeState(), DEFAULT_TIMEOUT)
			d.Lock()
			d.worker.Actual[link][name] = baseMessage
			d.Unlock()
			time.Sleep(time.Second * time.Duration(DEFAULT_TIMEOUT))
		}
		return
	}

	baseMessage.Status = "download and compiling complete"
	d.Lock()
	d.worker.Actual[link][name] = baseMessage
	d.Unlock()
	if err := d.qdb.Update(link, db.StatusDONE, &name); err != nil {
		baseMessage.Status += fmt.Sprintf("\n- - - - - - -update to done status failed: [%s]", err.Error())
		d.Lock()
		d.worker.Actual[link][name] = baseMessage
		d.Unlock()
		fmt.Printf("\ndownloader update status db error: %s\n", err.Error())
	}
	if !d.eagerMode {
		baseMessage.Status += fmt.Sprintf("\n- - - - - - - - -EagleMode: %s\n- - - - - - - - - -Waiting %.2f second to next", d.EagerModeState(), duration)
		d.Lock()
		d.worker.Actual[link][name] = baseMessage
		d.Unlock()
		time.Sleep(time.Second * time.Duration(duration))
		return
	}

	baseMessage.Status += fmt.Sprintf("\n- - - - - - - - -EagleMode: %s\n- - - - - - - - - -Waiting %.2f second to next", d.EagerModeState(), duration)
	d.Lock()
	d.worker.Actual[link][name] = baseMessage
	d.Unlock()
	time.Sleep(time.Second * DEFAULT_TIMEOUT)
}
