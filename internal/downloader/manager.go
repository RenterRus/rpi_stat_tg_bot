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

var (
	qualityMapping = map[int]int{
		0: 4320,
		1: 2160,
		2: 1440,
		3: 1080,
		4: 720,
	}
)

type DLP struct {
	worker     WorkerStatus
	path       string
	dl         *ytdlp.Command
	totalRetry atomic.Int64
	retryInit  atomic.Int64
	qdb        db.Queue
	maxWorkers int
	updateStat string
	eagerMode  bool
	quality    int
}

func NewDownloader(path string, db db.Queue, maxWorkers int, eagerMode bool) Downloader {
	return &DLP{
		worker: WorkerStatus{
			Actual: make(map[string]map[string]FileInfo),
		},
		path:       path,
		qdb:        db,
		maxWorkers: maxWorkers,
		eagerMode:  eagerMode,
		quality:    0,
	}
}
