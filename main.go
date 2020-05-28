package main

import (
	_ "UpgraderServer/routers"
	"github.com/astaxie/beego/orm"

	"github.com/astaxie/beego"
	_ "github.com/mattn/go-sqlite3"
)

func init()  {
    orm.RegisterDriver("sqlite", orm.DRSqlite)
    orm.RegisterDataBase("default", "sqlite3", "database/upgraderserver.db")
    orm.RunSyncdb("default", false, true)
}


func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
