package databases

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Cities Data Base
var cityDB *sql.DB

// Init Cities Data Base
//
func InitCityDB(dataSourceName string, maxIdleConns int, maxOpenConns int) (err error) {
	cityDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	cityDB.SetMaxIdleConns(maxIdleConns)
	cityDB.SetMaxOpenConns(maxOpenConns)

	err = cityDB.Ping()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Cities Data Base was successfully initialized on \"%s\"", dataSourceName))
	return nil
}

// Check city in Cities Data Base
//
func CheckCityDB(cityName string) (flag bool, err error) {
	// Search in DB
	log.Println(fmt.Sprintf("Searching in Cities Data Base \"%s\"", cityName))

	rows, err := cityDB.Query(fmt.Sprintf("SELECT name FROM cities WHERE name = '%s'", cityName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	// Search in OSM
	log.Println(fmt.Sprintf("Searching in OSM \"%s\"", cityName))

	flag, err = checkInOSM(fmt.Sprintf("http://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1&featuretype=city", strings.Replace(cityName, " ", "+", -1)))
	if err != nil {
		return false, err
	}

	if flag {
		rows, err := cityDB.Query("SELECT MAX(id) FROM cities")
		if err != nil {
			return false, err
		}

		id := 0
		if rows.Next() {
			err = rows.Scan(&id)
			if err != nil {
				return false, err
			}
		} else {
			// return // Ошибка
		}
		id++
		_, err = cityDB.Exec(fmt.Sprintf("INSERT INTO cities VALUES(%d, '%s', 1)", id, cityName))
		if err != nil {
			return false, err
		}
	}

	return flag, nil
}

// Get info about city from OSM nominatim
func checkInOSM(url string) (flag bool, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	str := string(body)
	if strings.Contains(str, "\"type\":\"city\"") || strings.Contains(str, "\"type\":\"town\"") {
		return true, nil
	} else {
		return false, nil
	}
}

// Close Cities Data Base
//
func CloseCityDB() (err error) {
	err = cityDB.Close()
	if err != nil {
		return err
	}

	log.Println("Cities Data Base was successfully closed")
	return nil
}

// Players Data Base
var playerDB *sql.DB

// Init Players Data Base
//
func InitPlayerDB(dataSourceName string, maxIdleConns int, maxOpenConns int) (err error) {
	playerDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	playerDB.SetMaxIdleConns(maxIdleConns)
	playerDB.SetMaxOpenConns(maxOpenConns)

	err = playerDB.Ping()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Players Data Base was successfully initialized on \"%s\"", dataSourceName))
	return nil
}

// To parse player's rows
type player struct {
	id   int
	name string
	pass string
}

// Check player in Players Data Base
//
func CheckPlayerDB(name string, pass string) (flagName bool, flagPass bool, err error) {
	log.Println(fmt.Sprintf("Searching in Players Data Base \"%s : %s\"", name, pass))

	rows, err := playerDB.Query(fmt.Sprintf("SELECT * FROM players WHERE name = '%s'", name))
	if err != nil {
		return false, false, err
	}
	defer rows.Close()

	plr := new(player)
	if rows.Next() {
		err := rows.Scan(&plr.id, &plr.name, &plr.pass)
		if err != nil {
			return false, false, err
		}
	} else {
		return false, false, nil
	}

	if plr.name == name {
		flagName = true
	}
	if plr.pass == pass {
		flagPass = true
	}

	return flagName, flagPass, nil
}

// Add player in Players Data Base
//
func AddPlayerDB(name string, pass string) (flag bool, err error) {
	log.Println(fmt.Sprintf("Adding in Players Data Base \"%s : %s\"", name, pass))

	ex, _, err := CheckPlayerDB(name, pass)
	if err != nil {
		return false, err
	}
	if ex {
		return false, nil
	}

	rows, err := playerDB.Query("SELECT MAX(id) FROM players")
	if err != nil {
		return false, err
	}

	id := 0
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return false, err
		}
	} else {
		// return // Ошибка
	}
	id++
	_, err = playerDB.Exec(fmt.Sprintf("INSERT INTO players VALUES(%d, '%s', '%s')", id, name, pass))
	if err != nil {
		return false, err
	}

	return true, nil

}

// Close Players Data Base
//
func ClosePlayerDB() (err error) {
	err = playerDB.Close()
	if err != nil {
		return err
	}

	log.Println("Players Data Base was successfully closed")
	return nil
}
