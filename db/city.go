package db

import (
	"time"

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

func CityExists(name string) bool {
	o := orm.NewOrm()
	city := CityModel{Name: name}

	err := o.Read(&city, "Name")
	return err == nil
}

func CityAdd(name string) (bool, error) {
	o := orm.NewOrm()
	city := CityModel{Name: name}
	created, _, err := o.ReadOrCreate(&city, "Name")
	return created, err
}
