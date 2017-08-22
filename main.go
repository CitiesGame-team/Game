package main

import (
	"flag"

	"log"

	"./config"
	"./db"
	"./game"

	"fmt"

	"github.com/astaxie/beego/orm"
)

func main() {
	confFile := flag.String("config", "./config.yml", "Configuration file")
	isInit := flag.Bool("init", false, "Init database")
	flag.Parse()

	conf, err := config.ReadProjectConfig(*confFile)
	if err != nil {
		panic(err)
	}

	err = db.InitDB(*conf.Db)
	if err != nil {
		panic(fmt.Sprintf("can't connect to db: %s", err))
	}

	if *isInit {
		err = orm.RunSyncdb("default", true, true)

		if err != nil {
			log.Printf("%s", err)
		}
	} else {
		game.RunGame(conf)
	}
}
