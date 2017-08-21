package databases

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Players Data Base
var playerDB *sql.DB

// InitPlayerDB inits Players Data Base
//
func InitPlayerDB(dataSourceName string, maxIdleConns int, maxOpenConns int) (err error) {
	playerDB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	playerDB.SetMaxIdleConns(maxIdleConns)
	playerDB.SetMaxOpenConns(maxOpenConns)

	err = playerDB.Ping()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Players Data Base was successfully initialized on \"%s\"", dataSourceName))
	return nil
}

// To parse player's rows
type player struct {
	id   int
	name string
	pass string
}

// CheckPlayerDB checks player in Players Data Base
//
func CheckPlayerDB(name string, pass string) (flagName bool, flagPass bool, err error) {
	log.Println(fmt.Sprintf("Searching in Players Data Base \"%s : %s\"", name, pass))

	rows, err := playerDB.Query(fmt.Sprintf("SELECT * FROM players WHERE name = '%s'", name))
	if err != nil {
		return false, false, err
	}
	defer rows.Close()

	plr := new(player)
	if rows.Next() {
		err := rows.Scan(&plr.id, &plr.name, &plr.pass)
		if err != nil {
			return false, false, err
		}
	} else {
		return false, false, nil
	}

	if plr.name == name {
		flagName = true
	}
	if plr.pass == pass {
		flagPass = true
	}

	return flagName, flagPass, nil
}

// AddPlayerDB adds player in Players Data Base
//
func AddPlayerDB(name string, pass string) (flag bool, err error) {
	log.Println(fmt.Sprintf("Adding in Players Data Base \"%s : %s\"", name, pass))

	ex, _, err := CheckPlayerDB(name, pass)
	if err != nil {
		return false, err
	}
	if ex {
		return false, nil
	}

	rows, err := playerDB.Query("SELECT MAX(id) FROM players")
	if err != nil {
		return false, err
	}

	id := 0
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return false, err
		}
	} else {
		// return // Ошибка
	}
	id++
	_, err = playerDB.Exec(fmt.Sprintf("INSERT INTO players VALUES(%d, '%s', '%s')", id, name, pass))
	if err != nil {
		return false, err
	}

	return true, nil

}

// ClosePlayerDB closes Players Data Base
//
func ClosePlayerDB() (err error) {
	err = playerDB.Close()
	if err != nil {
		return err
	}

	log.Println("Players Data Base was successfully closed")
	return nil
}
