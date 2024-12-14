package downloader

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

type FileInfo struct {
	Name         string
	DownloadSize string
	TotalSize    string
	Proc         string
	Status       string
}

type WorkerStatus struct {
	IsIdle  bool
	History map[string]map[string]FileInfo
	Actual  map[string]map[string]FileInfo
}

type DLP struct {
	queue       chan string
	failedQueue chan string
	worker      WorkerStatus
	path        string
}

func NewDownloader(path string) Downloader {
	return &DLP{
		queue:       make(chan string, 1),
		failedQueue: make(chan string, 100),
		worker: WorkerStatus{
			History: make(map[string]map[string]FileInfo),
			Actual:  make(map[string]map[string]FileInfo),
		},
		path: path,
	}
}

func (d *DLP) DownloadHistory() string {
	res := ""
	total_file := 0

	for k, v := range d.worker.History {
		res += fmt.Sprintf("\nLink: %s", k)
		for k_file, v_file := range v {
			total_file++
			res += fmt.Sprintf("\n- File: %s\n- - -Info:\n- - - - - -Name: %s\n- - - - - -DownloadedSize: %s\n- - - - - -TotalSize: %s\n- - - - - -Proc: %s\n- - - - - -Status: %s", k_file,
				v_file.Name, v_file.DownloadSize, v_file.TotalSize, v_file.Proc, v_file.Status)
		}
	}

	res += fmt.Sprintf("\n\nTotal files: %d\nTotal video: %d", total_file, len(d.worker.History))
	if d.worker.IsIdle {
		res += "\nDownloader is idle"
	} else {
		res += "\nDownloader is run"
	}

	return res
}

func (d *DLP) ActualStatus() string {
	if d.worker.IsIdle {
		return "Downloader is idle"
	}

	res := "Downloader is run"
	total_file := 0
	file_finished := 0

	for k, v := range d.worker.Actual {
		res += fmt.Sprintf("\nLink: %s", k)
		for k_file, v_file := range v {
			total_file++
			res += fmt.Sprintf("\n- File: %s\n- - -Info:\n- - - - - -Name: %s\n- - - - - -DownloadedSize: %s\n- - - - - -TotalSize: %s\n- - - - - -Proc: %s\n- - - - - -Status: %s", k_file,
				v_file.Name, v_file.DownloadSize, v_file.TotalSize, v_file.Proc, v_file.Status)
			if v_file.Proc == "100" {
				file_finished++
			}
		}
	}

	status := "downloading"
	if total_file == file_finished {
		status = "converting"
		if file_finished == 0 {
			status = "preparing"
		}
	}

	res += fmt.Sprintf("\n\nTotal files: %d\nStatus for this download: %s\nQueue size: %d\nFailed queue size (to repeat): %d", total_file, status, len(d.queue), len(d.failedQueue))

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

func (d *DLP) Run(ctx context.Context) {
	//ytdlp.MustInstall(context.TODO(), nil)
	ytdlp.Install(context.TODO(), nil)
	dl := ytdlp.New().
		UnsetCacheDir().
		ProgressDelta(3).
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
	d.worker.History = make(map[string]map[string]FileInfo)

	for {
		select {
		case <-ctx.Done():
			return
		case link := <-d.queue:
			d.worker.IsIdle = false

			// Ультра чистка
			for k := range d.worker.Actual {
				delete(d.worker.Actual, k)
			}

			progressInfo := map[string]FileInfo{}

			dl.ProgressFunc(time.Duration(time.Second*3), func(update ytdlp.ProgressUpdate) {
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

				d.worker.History[link] = progressInfo
				d.worker.Actual[link] = progressInfo
			})
			_, err := dl.Run(context.TODO(), link)
			if err != nil {
				d.failedQueue <- link
				fmt.Println(err)
			}
			time.Sleep(time.Second * 17)
		default:
			d.worker.IsIdle = true
			time.Sleep(time.Second * 3)
		}
	}
}
