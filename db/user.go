package db

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type UserModel struct {
	Id       int       `orm:"auto"`
	Name     string    `orm:"unique"`
	PassHash string    `orm:"size(53)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`

	Games []*GameModel `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(UserModel))
}
