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
		res += strings.ToUpper(mode)
		ch = len(links)
		for k, v := range links {
			res += fmt.Sprintf("%d. %s\n", k, v)
		}
	}

	return res, ch
}

func (d *DLP) DownloadHistory() string {
	res := "\n"

	queueCH := 0
	workCH := 0
	doneCH := 0

	history := ""
	history, queueCH = d.getHistory(db.StatusNEW)
	res += fmt.Sprintf("%s\n", history)

	history, workCH = d.getHistory(db.StatusWORK)
	res += fmt.Sprintf("%s\n", history)

	history, doneCH = d.getHistory(db.StatusDONE)
	res += fmt.Sprintf("%s\n", history)

	res += fmt.Sprintf("\nTotal:\n--In queue: %d\n--In work: %d\n--Is done: %d\n", queueCH, workCH, doneCH)

	return res
}
