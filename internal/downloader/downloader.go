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
}

func NewDownloader() Downloader {
	return &DLP{
		queue:       make(chan string, 1),
		failedQueue: make(chan string, 100),
		worker: WorkerStatus{
			History: make(map[string]map[string]FileInfo),
			Actual:  make(map[string]map[string]FileInfo),
		},
	}
}

func (d *DLP) DownloadHistory() string {
	res := ""
	total_file := 0

	for k, v := range d.worker.History {
		res += fmt.Sprintf("link: %s", k)
		for k_file, v_file := range v {
			total_file++
			res += fmt.Sprintf("- File: %s\n- -Info:\n- - -Name: %s\n- - -DownloadedSize: %s\n- - -TotalSize: %s\n- - -Proc: %s\n- - -Status: %s", k_file,
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

	for k, v := range d.worker.History {
		res += fmt.Sprintf("link: %s", k)
		for k_file, v_file := range v {
			total_file++
			res += fmt.Sprintf("- File: %s\n- -Info:\n- - -Name: %s\n- - -DownloadedSize: %s\n- - -TotalSize: %s\n- - -Proc: %s\n- - -Status: %s", k_file,
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

	res += fmt.Sprintf("\n\nTotal files: %d\nStatus for this download: %s", total_file, status)

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
			time.Sleep(time.Minute * 17)
		}
	}
}

func (d *DLP) Run(ctx context.Context) {
	dl := ytdlp.New().SetWorkDir("/home/ftppi/raid/Посмотреть/bot_download").
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		Output("%(title)s.%(ext)s").
		NoRestrictFilenames().Fixup(ytdlp.FixupForce).AbortOnError()

	go func() {
		d.fromFailed(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case link := <-d.queue:
			d.worker.IsIdle = false
			d.worker.Actual = make(map[string]map[string]FileInfo)

			dl.ProgressFunc(time.Duration(time.Second*3), func(update ytdlp.ProgressUpdate) {
				size := (float64(update.DownloadedBytes) / 1024) / 1024
				totalSize := (float64(update.TotalBytes) / 1024) / 1024
				fmt.Println(update.Status, update.PercentString(), fmt.Sprintf("[%d/%d] mb", int(size), int(totalSize)), update.Filename)
				d.worker.History[link][update.Filename] = FileInfo{
					Name:         update.Filename,
					DownloadSize: strconv.Itoa(int(size)),
					TotalSize:    strconv.Itoa(int(totalSize)),
					Proc:         update.PercentString(),
					Status:       string(update.Status),
				}

				d.worker.Actual[link][update.Filename] = FileInfo{
					Name:         update.Filename,
					DownloadSize: strconv.Itoa(int(size)),
					TotalSize:    strconv.Itoa(int(totalSize)),
					Proc:         update.PercentString(),
					Status:       string(update.Status),
				}
			})
			_, err := dl.Run(context.TODO(), link)
			if err != nil {
				d.failedQueue <- link
				fmt.Println(err)
				return
			}
			time.Sleep(time.Second * 7)
		default:
			d.worker.IsIdle = true
			d.worker.Actual = nil
			time.Sleep(time.Second * 3)
		}
	}
}
