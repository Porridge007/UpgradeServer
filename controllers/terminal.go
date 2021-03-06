package controllers

import (
	"UpgraderServer/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

type UpdateLatestController struct {
	beego.Controller
}

type UpdateGivenController struct {
	beego.Controller
}


func (c *QueryLatestController) Get(){
	deviceName := c.GetString("device")
	o := orm.NewOrm()
	device := models.Device{Device:deviceName}
	var pack models.Package

	o.QueryTable("device").Filter("device", deviceName).One(&device)
	o.QueryTable("package").Filter("device", device.Id).OrderBy("-id").One(&pack)

	c.Ctx.ResponseWriter.Write([]byte(pack.Version))
}

func (c *UpdateLatestController) Post() {
	deviceName := c.GetString("device")
	o := orm.NewOrm()
	device := models.Device{Device:deviceName}
	var pack models.Package

	o.QueryTable("device").Filter("device", deviceName).One(&device)
	o.QueryTable("package").Filter("device", device.Id).OrderBy("-id").One(&pack)
	len := len(strings.Split(pack.Address, "/"))
	transferAddr := "upload/"+strings.Split(pack.Address,"/")[len-1:][0]
	fmt.Println(transferAddr)

	f, err := os.Open(transferAddr)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/octect-stream")
	c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment;filename=\""+pack.Name+"\"")
	c.Ctx.ResponseWriter.Write(data)
}

func (c *UpdateGivenController) Post(){
	deviceName := c.GetString("device")
	version := c.GetString("version")
	o := orm.NewOrm()
	device := models.Device{Device:deviceName}
	var pack models.Package
	o.QueryTable("device").Filter("device", deviceName).One(&device)
	o.QueryTable("package").Filter("device",device.Id).Filter("version",version).One(&pack)

	len := len(strings.Split(pack.Address, "/"))
	transferAddr := "upload/"+strings.Split(pack.Address,"/")[len-1:][0]
	fmt.Println(transferAddr)

	f, err := os.Open(transferAddr)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/octect-stream")
	c.Ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment;filename=\""+pack.Name+"\"")
	c.Ctx.ResponseWriter.Write(data)
}

