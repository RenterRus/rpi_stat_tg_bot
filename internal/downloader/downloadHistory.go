package downloader

import (
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strings"
)

func (d *DLP) getHistory(mode string) (string, int) {
	res := ""
	ch := 0
	links, err := d.qdb.SelectAll(mode)
	if err != nil {
		res += fmt.Sprintf("Error get '%s' links: %s\n", mode, err.Error())
	} else {
		res += "\n" + strings.ToUpper(mode) + "\n"
		ch = len(links)
		for k, v := range links {
			name := "Coming soon"
			if v.Name != nil {
				name = *v.Name
			}
			res += fmt.Sprintf("%d. %s [%s]\n", (k + 1), name, v.Link)
		}
	}

	return res, ch
}

func (d *DLP) DownloadHistory() string {
	res := ""
	for _, v := range []string{db.StatusDONE, db.StatusWORK, db.StatusNEW} {
		history, _ := d.getHistory(v)
		res += fmt.Sprintf("%s\n", history)
	}

	res += fmt.Sprintf("\nRetry: %d", d.totalRetry.Load())

	return res
}
