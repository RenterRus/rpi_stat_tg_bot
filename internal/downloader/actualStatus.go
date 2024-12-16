package downloader

import (
	"fmt"
	"rpi_stat_tg_bot/internal/db"
	"strings"
)

func (d *DLP) ActualStatus() string {
	res := ""
	total_file := 0
	file_finished := 0

	// Формируем основное сообщение
	for k, v := range d.worker.Actual {
		res += fmt.Sprintf("\nLink: %s", k)
		for k_file, v_file := range v {
			total_file++
			res += fmt.Sprintf("\n- File: %s\n- - -Info:\n- - - - - -Name: %s\n- - - - - -Downloaded: [%s/%s]mb\n- - - - - -Proc: %s\n- - - - - -Status: %s", k_file,
				v_file.Name, v_file.DownloadSize, v_file.TotalSize, v_file.Proc, v_file.Status)
			if strings.Contains(v_file.Proc, "100") {
				file_finished++
			}
		}
	}

	// Докидываем сахара
	res += fmt.Sprintf("\n\nFiles in work right now: %d\nFailed queue size (to repeat): %d of %d\n", total_file, len(d.failedQueue), BASE_BUF_QUEUE_SIZE)

	if len(d.failedQueue) == BASE_BUF_QUEUE_SIZE {
		res += "\nThe queue for downloading failed attempts is full, new videos will be queued for download and processing as the queue is released"
	}

	// Включаем кусочек общей статистики, т.к. это дает дполонительное понимание в каком месте находится сервис
	_, queueCH := d.getHistory(db.StatusNEW)
	_, workCH := d.getHistory(db.StatusWORK)
	_, doneCH := d.getHistory(db.StatusDONE)

	res += fmt.Sprintf("\nTotal:\n--In queue: %d\n--In work: %d\n--Is done: %d\n", queueCH, workCH, doneCH)

	return res

}
