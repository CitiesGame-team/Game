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

// InitCityDB inits Cities Data Base
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

// CheckCityDB checks city in Cities Data Base
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

// checkInOSM gets info about city from OSM nominatim
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

// CloseCityDB closes Cities Data Base
//
func CloseCityDB() (err error) {
	err = cityDB.Close()
	if err != nil {
		return err
	}

	log.Println("Cities Data Base was successfully closed")
	return nil
}
