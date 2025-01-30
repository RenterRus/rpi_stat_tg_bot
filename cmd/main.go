package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"rpi_stat_tg_bot/internal/app"
	"time"
)

var path *string

func init() {
	path = flag.String("config", "../config.yaml", "path to config. Example: ../config.yaml")

	flag.Parse()
	if path == nil || len(*path) < 6 {
		log.Fatal("config flag not found")
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Lazy waiting for full load device")
	time.Sleep(time.Second * 2)
	a := app.NewApp(*path)
	a.Run()
}
