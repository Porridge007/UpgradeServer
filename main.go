package main

import (
	_ "UpgraderServer/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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
	beego.SetStaticPath("/download", "upload")
	beego.Run()
}
