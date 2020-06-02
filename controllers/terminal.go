package controllers

import (
	"UpgraderServer/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"net/http"
	"time"
)

type File struct{
	Id int64
	File_sha1 string `orm:"unique"`
	File_name string
	File_size int64
	File_addr string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now_add;type(datetime)"`
	Status int64 `orm:"null"`
	Device string
	Version string
}

type QueryLatestController struct {
	beego.Controller
}

func (c *QueryLatestController) Get(){
	deviceName := c.GetString("device")
	o := orm.NewOrm()
	device := models.Device{Device:deviceName}
	var pack models.Package

	o.QueryTable("device").Filter("device", deviceName).One(&device)
	fmt.Println(device.Id)
	o.QueryTable("package").Filter("device", device.Id).OrderBy("-id").One(&pack)
	file := File{
		Id:        pack.Id,
		File_sha1: "",
		File_name: pack.Name,
		File_size: 0,
		File_addr: pack.Address,
		Created:   time.Time{},
		Updated:   time.Time{},
		Status:    0,
		Device:    deviceName,
		Version:   pack.Version,
	}
	data, err := json.Marshal(file)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Ctx.ResponseWriter.Write(data)

}