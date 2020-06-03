package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"net"
	"os"
)

func SeverAddr() string{
	conf_ip := beego.AppConfig.String("hostip")
	if conf_ip != ""{
		return conf_ip+":"+beego.AppConfig.String("httpport")
	}

	addrs, err := net.InterfaceAddrs()
	var ips []string

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}

		}
	}
	return ips[0] +":"+beego.AppConfig.String("httpport")
}