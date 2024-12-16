package downloader

import "fmt"

func (d *DLP) CleanHistory() string {
	if err := d.qdb.Delete(); err != nil {
		return fmt.Errorf("CleanHistory: %w", err).Error()
	}

	return fmt.Sprintf("The history has been cleared\n\n%s\n", d.DownloadHistory())
}