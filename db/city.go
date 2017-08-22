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
