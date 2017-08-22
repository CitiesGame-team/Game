package main

import (
	"flag"

	"./config"
	"./game"
)

func main() {
	confFile := flag.String("config", "./config.yml", "Configuration file")
	flag.Parse()

	conf, err := config.ReadProjectConfig(*confFile)
	if err != nil {
		panic(err)
	}

	game.RunGame(conf)
}
