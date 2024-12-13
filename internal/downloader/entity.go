package downloader

import "context"

type Downloader interface {
	ToDownload(url string)
	Run(ctx context.Context)
	DownloadHistory() string
	ActualStatus() string
}
