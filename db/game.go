package db

import (
	"time"

	"fmt"

	"github.com/astaxie/beego/orm"
)

type GameState int16

const (
	GAME_CREATED GameState = iota
	GAME_STARTED
	GAME_FINISHED
)

type GameModel struct {
	Id      int        `orm:"auto"`
	State   GameState  `orm:"default(0)"`
	User1   *UserModel `orm:"rel(fk)"`
	User2   *UserModel `orm:"rel(fk)"`
	Created time.Time  `orm:"auto_now_add;type(datetime)"`

	Cities []*CityModel `orm:"rel(m2m)"`
}

func init() {
	orm.RegisterModel(new(GameModel))
}

func GameGet(id int) (GameModel, error) {
	o := orm.NewOrm()
	game := GameModel{Id: id}

	err := o.Read(&game)
	return game, err
}

func GameGetById(id int) (GameModel, error) {
	return GameGet(id)
}

func GameAdd(u1 UserModel, u2 UserModel, state GameState) (int, error) {
	o := orm.NewOrm()
	game := &GameModel{
		State: state,
		User1: &u1,
		User2: &u2,
	}

	id, err := o.Insert(game)
	return int(id), err
}

func (game GameModel) AddCity(city CityModel) error {
	if game.HasCity(city) {
		return fmt.Errorf("city %q already exists in game with id=%d", city.Name, game.Id)
	}

	o := orm.NewOrm()

	m2m := o.QueryM2M(&game, "Cities")
	_, err := m2m.Add(city)
	return err
}

func (game GameModel) HasCity(city CityModel) bool {
	o := orm.NewOrm()

	m2m := o.QueryM2M(&game, "Cities")
	return m2m.Exist(&city)
}

func (game GameModel) ChangeState(state GameState) error {
	o := orm.NewOrm()

	game.State = state
	_, err := o.Update(&game, "State")
	return err
}
