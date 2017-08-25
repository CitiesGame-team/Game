package db

import (
	"time"

	"fmt"

	"math/rand"

	"github.com/astaxie/beego/orm"
)

type CityModel struct {
	Id      int       `orm:"auto"`
	Name    string    `orm:"unique"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`

	Games []*GameModel `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(CityModel))
}

func CityExists(name string) (bool, error) {
	_, err := CityGet(name)
	return err == nil, err
}

func CityGet(name string) (CityModel, error) {
	o := orm.NewOrm()
	city := CityModel{Name: name}

	err := o.Read(&city, "Name")
	return city, err
}

func CityAdd(name string) (bool, error) {
	o := orm.NewOrm()
	city := CityModel{Name: name}

	created, _, err := o.ReadOrCreate(&city, "Name")
	return created, err
}

func CityHintForGame(letter string, game GameModel) (CityModel, error) {
	var cities []*CityModel
	o := orm.NewOrm()

	qs := o.QueryTable(new(CityModel))
	qs.Filter("name__istartswith", letter).Limit(500).All(&cities)

	shuffleCities(cities)
	for _, city := range cities {
		if game.HasCity(*city) {
			continue
		}

		return *city, nil
	}

	return CityModel{}, fmt.Errorf("cannot find next city for game with id=%d starting with letter %q!", game.Id, letter)
}

func shuffleCities(s []*CityModel) {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}
