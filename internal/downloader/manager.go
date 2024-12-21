package downloader

import (
	"rpi_stat_tg_bot/internal/db"
	"sync/atomic"

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
	Actual map[string]map[string]FileInfo
}

type DLP struct {
	worker     WorkerStatus
	path       string
	dl         *ytdlp.Command
	totalRetry atomic.Int64
	qdb        db.Queue
	maxWorkers int
	updateStat string
}

func NewDownloader(path string, db db.Queue, maxWorkers int) Downloader {
	return &DLP{
		worker: WorkerStatus{
			Actual: make(map[string]map[string]FileInfo),
		},
		path:       path,
		qdb:        db,
		maxWorkers: maxWorkers,
	}
}
