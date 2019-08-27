package router

import (
	"echodemo/conf"
	"echodemo/controller"
	"echodemo/filter"
	"github.com/donnie4w/go-logger/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"strings"
)

var e *echo.Echo

func Start()  {
	logger.Debug("start router ...")
	e = echo.New()
	e.Static("/", "static")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(filter.OpenTracing())
	//e.Use(filter.ExpireCheack())
	//e.Group("/v1/", filter.ExpireCheack())
	//api.Use(filter.ExpireCheack())
	logger.Debug("666666")
	initRouter()
	e.Logger.Fatal(e.Start(":"+conf.Config.Port))
}


/**
api地址 拼接版本
*/
func genUrl(url string) string {
	if strings.HasPrefix(url, "/") {
		return "/" + conf.Config.ApiVersion + url
	} else {
		return "/" + conf.Config.ApiVersion + "/" + url
	}
}


func initRouter(){
	api := e.Group("/v2", filter.ExpireCheack())
	// 装饰器早于具体路由写入
	api.GET("/user/InserUser", controller.InserUser)

	e.GET(genUrl("/user/InserUser"), controller.InserUser)
	e.GET(genUrl("/user/UpdatePassword"), controller.UpdatePassword)
	e.GET(genUrl("/user/SelectUser"), controller.SelectUser)
	e.GET(genUrl("/user/DeleteUser"), controller.DeleteUser)
}