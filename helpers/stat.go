package helpers

import "github.com/astaxie/beego/orm"

type TopCity struct {
	Name  string
	Count int
}

func GetCitiesStat() ([]*TopCity, error) {
	var cities []*TopCity

	o := orm.NewOrm()

	_, err := o.Raw(
		`SELECT c.name AS name, COUNT(*) AS count
				FROM game_model_city_models m2m
				INNER JOIN city_model c ON c.id = m2m.city_model_id
				GROUP BY c.name
				ORDER BY COUNT(*) DESC
				LIMIT 10`).QueryRows(&cities)

	return cities, err
}

type GameStat struct {
	Games   int
	Players int
	Cities  int
}

func GetGameStat() (GameStat, error) {
	var stat GameStat

	o := orm.NewOrm()
	err := o.Raw(
		`SELECT
		(SELECT COUNT(*) FROM game_model) AS games,
		(SELECT COUNT(*) FROM user_model) AS players,
		(SELECT COUNT(*) FROM city_model) AS cities`).QueryRow(&stat)

	return stat, err
}
