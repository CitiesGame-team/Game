package main

import (
	"flag"
	"fmt"
	"log"

	"Game/config"
	"Game/db"
	"Game/game"

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

	connName := "default"
	err = db.InitDB(*conf.Db, connName)
	if err != nil {
		panic(fmt.Sprintf("can't connect to db: %s", err))
	}

	if *isInit {
		err = orm.RunSyncdb(connName, false, true)
		if err != nil {
			log.Printf("orm sync error: %s", err)
			return
		}

		cities, err := config.ReadCitiesBase(conf.CitiesBaseFile)

		if err != nil {
			log.Printf("cities base init error: %s", err)
			return
		}

		total := len(cities.Cities)
		for index, city := range cities.Cities {
			db.CityAdd(city)
			if (index+1)%500 == 0 || (total-index-1) == 0 {
				log.Printf("%d cities of %d", index+1, total)
			}
		}
	} else {
		game.RunGame(conf)
	}
}
