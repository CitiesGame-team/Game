package db

import (
	"time"

	"fmt"

	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Id       int       `orm:"auto"`
	Name     string    `orm:"unique"`
	PassHash string    `orm:"size(62)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
	Updated  time.Time `orm:"auto_now;type(datetime)"`

	Games []*GameModel `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(UserModel))
}

func UserExists(name string) bool {
	_, err := UserGet(name)
	return err == nil
}

func UserGet(name string) (UserModel, error) {
	o := orm.NewOrm()
	user := UserModel{Name: name}

	err := o.Read(&user, "Name")
	return user, err
}

func UserExistsById(id int) bool {
	_, err := UserGetById(id)
	return err == nil
}

func UserGetById(id int) (UserModel, error) {
	o := orm.NewOrm()
	user := UserModel{Id: id}

	err := o.Read(&user)
	return user, err
}

func UserAdd(name string, password []byte) (bool, int, error) {
	hash, err := hashPass(password)

	if err != nil {
		return false, -1, err
	}

	o := orm.NewOrm()
	user := UserModel{Name: name, PassHash: fmt.Sprintf("%s", hash)}
	created, id, err := o.ReadOrCreate(&user, "Name")
	return created, int(id), err
}

func UserAuth(name string, password []byte) error {
	user, err := UserGet(name)

	if err != nil {
		return err
	}

	return hashPassCheck([]byte(user.PassHash), password)
}

func hashPass(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 10)
}

func hashPassCheck(hash []byte, pass []byte) error {
	return bcrypt.CompareHashAndPassword(hash, pass)
}
