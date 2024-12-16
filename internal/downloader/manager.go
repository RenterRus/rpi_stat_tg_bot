package downloader

import (
	"rpi_stat_tg_bot/internal/db"
	"sync"
	"sync/atomic"

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
	sync.Mutex
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
