package downloader

import "context"

type Downloader interface {
	ToDownload(url string) error
	Run(ctx context.Context)
	DownloadHistory() string
	ActualStatus() string
	CleanHistory() string
}
