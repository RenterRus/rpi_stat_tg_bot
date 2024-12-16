package downloader

import (
	"context"
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strconv"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

func (d *DLP) ToDownload(url string) error {
	if err := d.qdb.Insert(url); err != nil {
		return fmt.Errorf("to download: %w", err)
	}
	return nil
}

func (d *DLP) fromFailed(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case link := <-d.failedQueue:
			d.ToDownload(link)
		default:
			time.Sleep(time.Minute * 3)
		}
	}
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

		if err := d.qdb.Update(link, db.StatusDONE); err != nil {
			fmt.Printf("\ndownloader updare db error: %s\n", err.Error())
		}

		d.totalComplete.Add(1)

		fmt.Printf("\n\nVIDEO %s\nLINK: %s\nIS DONE\n\n", name, link)
		delete(d.worker.Actual, link)
	}()

	progressInfo := map[string]FileInfo{}

	d.dl.ProgressFunc(time.Duration(time.Millisecond*500), func(update ytdlp.ProgressUpdate) {
		size := (float64(update.DownloadedBytes) / 1024) / 1024
		totalSize := (float64(update.TotalBytes) / 1024) / 1024
		fmt.Println(update.Status, update.PercentString(), fmt.Sprintf("[%d/%d] mb", int(size), int(totalSize)), update.Filename)
		progressInfo[update.Filename] = FileInfo{
			Name:         d.path + "/" + update.Filename,
			DownloadSize: strconv.Itoa(int(size)),
			TotalSize:    strconv.Itoa(int(totalSize)),
			Proc:         update.PercentString(),
			Status:       string(update.Status),
		}

		d.worker.Actual[link] = progressInfo
	})

	if err := d.qdb.Update(link, db.StatusWORK); err != nil {
		fmt.Printf("\ndownloader updare db error: %s\n", err.Error())
	}

	_, err := d.dl.Run(context.TODO(), link)
	if err != nil {
		d.failedQueue <- link
		fmt.Println(err)
	}
	// Даем процессору "отдохнуть". Ему реально было не просто
	time.Sleep(time.Second * 7)
}
