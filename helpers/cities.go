package helpers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"log"

	"../db"
)

func CityExists(name string) (bool, error) {
	inDb, err := db.CityExists(name)
	log.Printf("check city in db %q: %s (%s)", name, inDb, err)

	if err != nil || !inDb {
		inOSM, err := checkCityOSM(name)
		log.Printf("check city in osm %q: %s (%s)", name, inOSM, err)

		if err != nil || !inOSM {
			return false, err
		}

		db.CityAdd(name)
	}

	return true, nil
}

func checkCityOSM(name string) (bool, error) {
	url := "http://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1&featuretype=city"
	encName := strings.Replace(name, " ", "+", -1)

	resp, err := http.Get(fmt.Sprintf(url, encName))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	str := string(body)
	isCity, isTown := strings.Contains(str, "\"type\":\"city\""), strings.Contains(str, "\"type\":\"town\"")

	return isCity || isTown, nil
}
