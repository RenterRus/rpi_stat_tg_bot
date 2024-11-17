package cmd

import (
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
)

type CMD struct {
	informer informer.Informer
	finder   finder.Finder
}
