package downloader

import (
	"fmt"
	"rpi_stat_tg_bot/internal/db"
)

func (d *DLP) ActualStatus() string {
	res := ""
	total_file := 0

	// Формируем основное сообщение
	for k, v := range d.worker.Actual {
		res += fmt.Sprintf("\nLink: %s", k)
		for k_file, v_file := range v {
			total_file++
			res += fmt.Sprintf("\n- File: %s\n- - -Info:\n- - - - - -Name: %s\n- - - - - -Downloaded: [%s/%s]mb or %s\n- - - - - -Status: %s\n", k_file,
				v_file.Name, v_file.DownloadSize, v_file.TotalSize, v_file.Proc, v_file.Status)
		}
	}

	// Докидываем сахара
	res += fmt.Sprintf("\n\nFiles in work right now: %d", total_file)

	// Включаем кусочек общей статистики, т.к. это дает дполонительное понимание в каком месте находится сервис
	_, queueCH := d.getHistory(db.StatusNEW)
	_, workCH := d.getHistory(db.StatusWORK)
	_, doneCH := d.getHistory(db.StatusDONE)

	res += fmt.Sprintf("\nTotal:\n--In queue: %d\n--In work: %d\n--Is done: %d\n--Retry: %d", queueCH, workCH, doneCH, d.totalRetry.Load())

	return res

}
