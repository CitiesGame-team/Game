package db

import (
	"../config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB(dbConf config.DatabaseConfig) error {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	return orm.RegisterDataBase("default", "mysql",
		dbConf.ConnStr, dbConf.PoolMaxIdle, dbConf.PoolMaxOpen)
}

func GetDB() {

}
