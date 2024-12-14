package downloader

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

const BASE_BUF_QUEUE_SIZE = 10

type FileInfo struct {
	Name         string
	DownloadSize string
	TotalSize    string
	Proc         string
	Status       string
}

type WorkerStatus struct {
	History []string
	Actual  map[string]map[string]FileInfo
}

type DLP struct {
	queue         chan string
	failedQueue   chan string
	worker        WorkerStatus
	path          string
	dl            *ytdlp.Command
	totalComplete atomic.Int64
}

func NewDownloader(path string) Downloader {
	return &DLP{
		queue:       make(chan string, BASE_BUF_QUEUE_SIZE),
		failedQueue: make(chan string, BASE_BUF_QUEUE_SIZE*BASE_BUF_QUEUE_SIZE),
		worker: WorkerStatus{
			History: make([]string, 0, BASE_BUF_QUEUE_SIZE),
			Actual:  make(map[string]map[string]FileInfo),
		},
		path: path,
	}
}

func (d *DLP) DownloadHistory() string {
	res := "\n"

	for k, v := range d.worker.History {
		res += fmt.Sprintf("%d. %s\n", k, v)
	}

	res += fmt.Sprintf("View %d last done video\nTotal video is done: %d", len(d.worker.History), d.totalComplete.Load())

	return res
}

func (d *DLP) CleanHistory() string {
	d.worker.History = make([]string, 0, BASE_BUF_QUEUE_SIZE)
	for k := range d.worker.Actual {
		delete(d.worker.Actual, k)
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

	res += fmt.Sprintf("\n\nFiles in work right now: %d\nQueue size: %d of %d\nFailed queue size (to repeat): %d of %d\n", total_file, len(d.queue), BASE_BUF_QUEUE_SIZE, len(d.failedQueue), BASE_BUF_QUEUE_SIZE*BASE_BUF_QUEUE_SIZE)
	if len(d.queue) == BASE_BUF_QUEUE_SIZE {
		res += "\nThe download queue is full, new videos will be queued for download and processing as the queue is released"
	}

	if len(d.failedQueue) == BASE_BUF_QUEUE_SIZE*BASE_BUF_QUEUE_SIZE {
		res += "\nThe queue for downloading failed attempts is full, new videos will be queued for download and processing as the queue is released"
	}

	return res

}

func (d *DLP) ToDownload(url string) {
	d.queue <- url
}

func (d *DLP) fromFailed(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case link := <-d.failedQueue:
			d.ToDownload(link)
		default:
			time.Sleep(time.Minute * 7)
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
		d.totalComplete.Add(1)
		d.worker.History = append(d.worker.History, fmt.Sprintf("sDONE: [%s] %s", link, name))

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

	_, err := d.dl.Run(context.TODO(), link)
	if err != nil {
		d.failedQueue <- link
		fmt.Println(err)
	}
	// Даем процессору "отдохнуть". Ему реально было не просто
	time.Sleep(time.Second * 17)
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
	d.worker.History = make([]string, 0, BASE_BUF_QUEUE_SIZE)
	doubleWay := make(chan struct{}, 2)

	for {
		select {
		case <-ctx.Done():
			return
		case doubleWay <- struct{}{}:
			go func() {
				defer func() {
					<-doubleWay
				}()

				d.downloader(<-d.queue)
			}()

		default:
			time.Sleep(time.Second * 3)
		}
	}
}
