package downloader

import "fmt"

func (d *DLP) CleanHistory() string {
	for k := range d.worker.Actual {
		delete(d.worker.Actual, k)
	}

	if err := d.qdb.Delete(); err != nil {
		return fmt.Errorf("CleanHistory: %w", err).Error()
	}

	return fmt.Sprintf("The history has been cleared\n\n%s\n", d.DownloadHistory())
}
