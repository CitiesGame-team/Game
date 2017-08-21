package databases

import (
	"fmt"
	"log"
)

// CurrentGame is a type for Current Game "Data Base"
//
type CurrentGame struct {
	ID1    int
	ID2    int
	Used   []string
	Winner int
}

// InitCurrentGameDB inits Current Game "Data Base"
//
func (c *CurrentGame) InitCurrentGameDB(ID1 int, ID2 int) {
	c.ID1 = ID1
	c.ID2 = ID2
	log.Println(fmt.Sprintf("Current Game Data Base was successfully initialized for players \"ID1:%d, ID2:%d\"", ID1, ID2))
}

// CityNameUsed checks if city name was used in Current Game "Data Base"
//
func (c *CurrentGame) CityNameUsed(cityName string) (flag bool) {
	for _, name := range c.Used {
		if name == cityName {
			log.Println(fmt.Sprintf("City name \"%s\" was found in Current Game Data Base for players \"ID1:%d, ID2:%d\"", cityName, c.ID1, c.ID2))
			return true
		}
	}
	log.Println(fmt.Sprintf("City name \"%s\" was not found in Current Game Data Base for players \"ID1:%d, ID2:%d\"", cityName, c.ID1, c.ID2))
	c.Used = append(c.Used, cityName)
	return false
}
