package downloader

import (
	"context"
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

const (
	BASE_BUF_QUEUE_SIZE = 10
	MAX_THREADS         = 2
)

type FileInfo struct {
	Name         string
	DownloadSize string
	TotalSize    string
	Proc         string
	Status       string
}

type WorkerStatus struct {
	Actual map[string]map[string]FileInfo
}

type DLP struct {
	failedQueue   chan string
	worker        WorkerStatus
	path          string
	dl            *ytdlp.Command
	totalComplete atomic.Int64
	qdb           db.Queue
}

func NewDownloader(path string, db db.Queue) Downloader {
	return &DLP{
		failedQueue: make(chan string, BASE_BUF_QUEUE_SIZE),
		worker: WorkerStatus{
			Actual: make(map[string]map[string]FileInfo),
		},
		path: path,
		qdb:  db,
	}
}

func (d *DLP) getHistory(mode string) (string, int) {
	res := ""
	ch := 0
	links, err := d.qdb.SelectAll(mode)
	if err != nil {
		res += fmt.Sprintf("Error get '%s' links: %s\n", mode, err.Error())
	} else {
		res += strings.ToUpper(mode)
		ch = len(links)
		for k, v := range links {
			res += fmt.Sprintf("%d. %s\n", k, v)
		}
	}

	return res, ch
}

func (d *DLP) DownloadHistory() string {
	res := "\n"

	queueCH := 0
	workCH := 0
	doneCH := 0

	history := ""
	history, queueCH = d.getHistory(db.StatusNEW)
	res += fmt.Sprintf("%s\n", history)

	history, workCH = d.getHistory(db.StatusWORK)
	res += fmt.Sprintf("%s\n", history)

	history, doneCH = d.getHistory(db.StatusDONE)
	res += fmt.Sprintf("%s\n", history)

	res += fmt.Sprintf("\nTotal:\n--In queue: %d\n--In work: %d\n--Is done: %d\n", queueCH, workCH, doneCH)

	return res
}

func (d *DLP) CleanHistory() string {
	for k := range d.worker.Actual {
		delete(d.worker.Actual, k)
	}

	if err := d.qdb.Delete(); err != nil {
		return fmt.Errorf("CleanHistory: %w", err).Error()
	}

	return fmt.Sprintf("The history has been cleared\n\n%s\n", d.DownloadHistory())
}

func (d *DLP) ActualStatus() string {
	res := ""
	total_file := 0
	file_finished := 0

	for k, v := range d.worker.Actual {
		res += fmt.Sprintf("\nLink: %s", k)
		for k_file, v_file := range v {
			total_file++
			res += fmt.Sprintf("\n- File: %s\n- - -Info:\n- - - - - -Name: %s\n- - - - - -DownloadedSize: %s\n- - - - - -TotalSize: %s\n- - - - - -Proc: %s\n- - - - - -Status: %s", k_file,
				v_file.Name, v_file.DownloadSize, v_file.TotalSize, v_file.Proc, v_file.Status)
			if strings.Contains(v_file.Proc, "100") {
				file_finished++
			}
		}
	}

	res += fmt.Sprintf("\n\nFiles in work right now: %d\nFailed queue size (to repeat): %d of %d\n", total_file, len(d.failedQueue), BASE_BUF_QUEUE_SIZE)

	if len(d.failedQueue) == BASE_BUF_QUEUE_SIZE {
		res += "\nThe queue for downloading failed attempts is full, new videos will be queued for download and processing as the queue is released"
	}

	_, queueCH := d.getHistory(db.StatusNEW)
	_, workCH := d.getHistory(db.StatusWORK)
	_, doneCH := d.getHistory(db.StatusDONE)

	res += fmt.Sprintf("\nTotal:\n--In queue: %d\n--In work: %d\n--Is done: %d\n", queueCH, workCH, doneCH)

	return res

}

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

	go func() {
		d.fromFailed(ctx)
	}()

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
					d.downloader(link)
				}
			}()

		default:
			time.Sleep(time.Second * 17)
		}
	}
}
