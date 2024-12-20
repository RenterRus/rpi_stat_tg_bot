package downloader

import (
	"context"
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func (d *DLP) ToDownload(url string) error {
	if err := d.qdb.Insert(url); err != nil {
		return fmt.Errorf("to download: %w", err)
	}
	return nil
}

func (d *DLP) downloader(link string) {
	defer func() {
		name := ""
		for i := range d.worker.Actual[link] {
			if name == "" {
				name = d.worker.Actual[link][i].Name
				break
			}
		}

		fmt.Printf("\n\nVIDEO %s\nLINK: %s\nIS DONE\n\n", name, link)
		delete(d.worker.Actual, link)
	}()

	progressInfo := map[string]FileInfo{}
	name := ""
	d.dl.ProgressFunc(time.Duration(time.Millisecond*750), func(update ytdlp.ProgressUpdate) {
		size := (float64(update.DownloadedBytes) / 1024) / 1024 // К мегабайтам
		totalSize := (float64(update.TotalBytes) / 1024) / 1024 // К мегабайтам
		fmt.Println(update.Status, update.PercentString(), fmt.Sprintf("[%.2f/%.2f]mb", size, totalSize), update.Filename)
		status := string(update.Status)
		if strings.Contains(status, "finished") {
			status = "converting"
		}
		progressInfo[update.Filename] = FileInfo{
			Name:         d.path + "/" + update.Filename,
			DownloadSize: fmt.Sprintf("%.2f", size),
			TotalSize:    fmt.Sprintf("%.2f", totalSize),
			Proc:         update.PercentString(),
			Status:       status,
		}
		d.worker.Actual[link] = progressInfo

		if name != update.Filename {
			name = update.Filename
			if err := d.qdb.Update(link, db.StatusWORK, &name); err != nil {
				fmt.Println("Update name (into work)", err)
			}
		}
	})

	_, err := d.dl.Run(context.TODO(), link)
	if err != nil {
		d.totalRetry.Add(1)
		if err := d.qdb.Update(link, db.StatusNEW, &name); err != nil {
			fmt.Printf("\ndownloader set video to queue (retry): %s\n", err.Error())
		}
	} else {
		if err := d.qdb.Update(link, db.StatusDONE, &name); err != nil {
			fmt.Printf("\ndownloader update db error: %s\n", err.Error())
		}
	}

	// Даем процессору "отдохнуть". Ему реально было не просто
	time.Sleep(time.Second * 7)
}
