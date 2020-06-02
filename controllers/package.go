package controllers

import (
	"UpgraderServer/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"os"
	"strings"
)

//  PackageController operations for Package
type PackageController struct {
	beego.Controller
}

// URLMapping ...
func (c *PackageController) URLMapping() {
	c.Mapping("Post", c.Post)
	//c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	//c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Package
// @Param   device   query   string  false       "device belongs to"
// @Success 201 {int} models.Package
// @Failure 403 body is empty
// @router / [post]
func (c *PackageController) Post() {
	os.Mkdir("upload",os.ModePerm)
	var v models.Package
	deviceId, _ :=  c.GetInt64("device")
	//var  device models.Device
	//_ =  orm.NewOrm().QueryTable("device").Filter("id", deviceId).One(&device)

	f, h, err := c.GetFile("file")
	if err != nil {
		log.Fatal("getfile err ", err)
	}
	defer f.Close()
	path :=  "upload/"+ h.Filename

	err = c.SaveToFile("file", path)
	if  err != nil{
		ret := models.Resp{
			Code: 401,
			Msg:  "Save Package Failure",
			Data: err.Error(),
		}
		c.Data["json"] =  ret
		c.ServeJSON()
		return
	}
	v.Name = h.Filename
	v.Version = getVersion(v.Name)
	v.Device = deviceId
	v.Address = strings.Split(c.Ctx.Request.RemoteAddr,":")[0]+ beego.AppConfig.String("httpport") +"/download/"+v.Name
	if _, err := models.AddPackage(&v); err == nil {
		updateLastestField(deviceId,v.Version)
		c.Ctx.Output.SetStatus(201)
		ret := models.Resp{
			Code: 200,
			Msg:  "Upload Package Success",
			Data: v,
		}
		c.Data["json"] = ret
	} else {
		ret := models.Resp{
			Code: 402,
			Msg:  "Upload Package Failure",
			Data: err.Error(),
		}
		c.Data["json"] = ret
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Package by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Package
// @Failure 403 :id is empty
// @router /:id [get]
//func (c *PackageController) GetOne() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.ParseInt(idStr, 0, 64)
//	v, err := models.GetPackageById(id)
//	if err != nil {
//		c.Data["json"] = err.Error()
//	} else {
//		c.Data["json"] = v
//	}
//	c.ServeJSON()
//}

// GetAll ...
// @Title Get All
// @Description get Package
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Package
// @Failure 403
// @router / [get]
func (c *PackageController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	deviceId, _ :=  c.GetInt64("device")
	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllPackage(query, fields, sortby, order, offset, limit, deviceId)
	if err != nil {
		ret := models.Resp{
			Code: 403,
			Msg:  "Get Package List Failure",
			Data: err.Error(),
		}
		c.Data["json"] = ret
	} else {
		ret := models.Resp{
			Code: 200,
			Msg:  "Get Package List Success",
			Data: l,
		}
		c.Data["json"] = ret
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Package
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Package	true		"body for Package content"
// @Success 200 {object} models.Package
// @Failure 403 :id is not int
// @router /:id [put]
//func (c *PackageController) Put() {
//	idStr := c.Ctx.Input.Param(":id")
//	id, _ := strconv.ParseInt(idStr, 0, 64)
//	v := models.Package{Id: id}
//	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
//	if err := models.UpdatePackageById(&v); err == nil {
//		c.Data["json"] = "OK"
//	} else {
//		c.Data["json"] = err.Error()
//	}
//	c.ServeJSON()
//}

// Delete ...
// @Title Delete
// @Description delete the Package
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /delete [post]
func (c *PackageController) Delete() {
	var toDelete []int64
	json.Unmarshal(c.Ctx.Input.RequestBody, &toDelete)
	var undeleted []int64
	for  _,id := range toDelete {
		var pack models.Package
		orm.NewOrm().QueryTable("package").Filter("id", id).One(&pack)
		if err := models.DeletePackage(int64(id));
			err != nil {
			undeleted = append(undeleted, id)
		}
		os.Remove("upload/"+pack.Name)
	}

		if undeleted == nil {
			ret := models.Resp{
				Code: 200,
				Msg:  "Delete Package Success",
				Data: "OK",
			}
			c.Data["json"] = ret
		} else {
			ret := models.Resp{
				Code: 404,
				Msg:  "Delete Package Failure",
				Data: undeleted,
			}
			c.Data["json"] = ret
		}
	c.ServeJSON()
}

func getVersion(filename string)  string{
	version_split := strings.Split(filename, ".")
	version_split = version_split[1:len(version_split)-1]
	version :=strings.Join(version_split,".")
	return version
}

func updateLastestField(deviceId int64, latestVersion string)  {
	o := orm.NewOrm()
	device := models.Device{Id:deviceId}
	if o.Read(&device) == nil {
		device.Latest = latestVersion
		if num, err := o.Update(&device); err == nil {
			fmt.Println(num)
		}
	}
}