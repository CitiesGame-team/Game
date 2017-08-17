package workWithDB

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	cityDB   *sql.DB
	playerDB *sql.DB
	cityMap  = map[byte][]string{}
)

// Init Cities Data Base
//
func InitCityBase(driverName, dataSourceName string) (err error) {
	cityDB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	for i := 'A'; i <= 'Z'; i++ {
		names, err := cityDB.Query(fmt.Sprintf("SELECT name FROM cities WHERE name LIKE '%c%%'", i))
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer names.Close()
		for names.Next() {
			nm := ""
			err := names.Scan(&nm)
			if err != nil {
				return err
			}
			cityMap[nm[0]] = append(cityMap[nm[0]], nm)
		}
	}
	return nil
}

// Check city in Cities Data Base
//
func CheckCityBase(cityName string) bool {
	for i := 0; i < len(cityMap[cityName[0]]); i++ {
		if cityMap[cityName[0]][i] == cityName {
			return true
		}
	}
	// Find in OSM
	return false
}

// Init Players Data Base
//
func InitPlayerBase(driverName, dataSourceName string) (err error) {
	playerDB, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	return nil
}
