package db

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type GameModel struct {
	Id      int        `orm:"auto"`
	State   int16      `orm:"default(0)"`
	User1Id *UserModel `orm:"rel(fk)"`
	User2Id *UserModel `orm:"rel(fk)"`
	Created time.Time  `orm:"auto_now_add;type(datetime)"`

	Cities []*CityModel `orm:"rel(m2m)"`
}

func init() {
	orm.RegisterModel(new(GameModel))
}
