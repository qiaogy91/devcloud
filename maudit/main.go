package main

import (
	//_ "github.com/qiaogy91/ioc/apps/health/restful"
	//_ "github.com/qiaogy91/ioc/apps/metrics/restful"
	//_ "github.com/qiaogy91/ioc/apps/swagger/restful"
	//_ "github.com/qiaogy91/ioc/config/cors/restful"
	//\
	_ "github.com/qiaogy91/ioc/apps/health/gin"
	_ "github.com/qiaogy91/ioc/apps/metrics/gin"
	_ "github.com/qiaogy91/ioc/config/cors/gin"
	"github.com/qiaogy91/ioc/server"

	_ "github.com/qiaogy91/devcloud/maudit/apps"
)

func main() {
	server.Start()
}
