package workWithDB

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Check city in Cities Data Base
//
func CheckCityBase(cityDB *sql.DB, cityName string, filePath string) (flag bool, err error) {
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
	/* Планирую сделать так, чтобы лишние пробелы из City удалялясиь. При запросе пробелы буду менять на +, чтобы иметь возможность
	искать такие города как "New York", "Nizhniy Tagil" и тд */
	err = getJSON(filePath, fmt.Sprintf("http://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1&featuretype=city", cityName))
	if err != nil {
		return false, err
	}
	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	str := string(f)
	/* Планирую не искать подстроку в строке, а нормально парсить  JSON */
	if strings.Contains(str, "\"type\":\"city\"") || strings.Contains(str, "\"type\":\"town\"") {
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
		return true, nil
	} else {
		return false, nil
	}
}

func getJSON(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
