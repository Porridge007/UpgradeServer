package controllers

import (
	"UpgraderServer/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid := models.AddUser(user)
	u.Data["json"] = map[string]string{"uid": uid}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid := u.GetString(":uid")
	if uid != "" {
		user, err := models.GetUser(uid)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)

		uu, err := models.UpdateUser(uid, &user)
		if user.Username == "admin" {
			beego.AppConfig.Set("username", user.Username)
			beego.AppConfig.Set("password", user.Password)
		}else if user.Username == "guest"{
			beego.AppConfig.Set("guest_name", user.Username)
			beego.AppConfig.Set("guest_password", user.Password)
		}
		beego.AppConfig.SaveConfigFile("conf/app.conf")
		if err != nil {
			ret := models.Resp{
				Code: 501,
				Msg:  "Modify User Failure",
				Data:  err.Error(),
			}
			u.Data["json"] = ret
		} else {
			ret := models.Resp{
				Code: 200,
				Msg:  "Modify User Success",
				Data:  uu,
			}
			u.Data["json"] = ret
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")
	//guest_name :=  u.GetString("guest_name")
	//guest_password := u.GetString("guest_password")
	if models.Login(username, password){
		ret := models.Resp{
			Code: 200,
			Msg:  "Login Success",
			Data:  "login success",
		}
		u.Data["json"] = ret
	} else {
		ret := models.Resp{
			Code: 502,
			Msg:  "Login Success",
			Data:  "user not exist",
		}
		u.Data["json"] = ret
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	ret := models.Resp{
		Code: 200,
		Msg:  "Login Success",
		Data:  "logout success",
	}
	u.Data["json"] = ret
	u.ServeJSON()
}

