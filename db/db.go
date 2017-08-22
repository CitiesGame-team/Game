package db

import (
	"../config"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB(dbConf config.DatabaseConfig, connName string) error {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	return orm.RegisterDataBase(connName, "mysql",
		dbConf.ConnStr, dbConf.PoolMaxIdle, dbConf.PoolMaxOpen)
}
