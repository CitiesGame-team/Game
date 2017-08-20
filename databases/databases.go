package databases

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type cityJSON struct {
	ClassPlace string `json:"class"`
	TypePlace  string `json:"type"`
}

// Init Cities Data Base
//
func InitCityBase(dataSourceName string) (cityDB *sql.DB, err error) {
	cityDB, err = sql.Open("mysql ", dataSourceName)
	return cityDB, err
}

// Check city in Cities Data Base
//
func CheckCityBase(cityDB *sql.DB, cityName string) (flag bool, err error) {
	// Search in DB
	rows, err := cityDB.Query(fmt.Sprintf("SELECT name FROM cities WHERE name = '%s'", cityName))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	// Search in OSM
	log.Println("Searching in OSM")

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
//
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

	var city []cityJSON
	err = json.Unmarshal(body, &city)
	if err != nil {
		return false, err
	}
	if (city[0].TypePlace == "city") || (city[0].TypePlace == "town") {
		return true, nil
	} else {
		return false, nil
	}
}
