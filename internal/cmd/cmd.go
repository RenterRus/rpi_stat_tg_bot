package cmd

import (
	"rpi_stat_tg_bot/internal/finder"
	"rpi_stat_tg_bot/internal/informer"
)

func NewCMD(informer informer.Informer, finder finder.Finder) CMD {
	const ttp = 1
	return CMD{
		informer: informer,
		finder:   finder,
		ttp:      ttp,
	}
}
