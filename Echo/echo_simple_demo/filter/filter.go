package filter

import (
	"echodemo/model"
	"echodemo/service"
	"github.com/donnie4w/go-logger/logger"
	"github.com/labstack/echo"
	"net/http"
)

var (
	User *model.User
)

/**
此方法为公共拦截器
这里做是否登录、权限校验等操作
给全局变量设值可以在Controller 中获取到
*/

func OpenTracing() echo.MiddlewareFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//从请求头中获取请求参数
			//token := context.Request().Header.Get("t")
			//uid := context.Request().Header.Get("u")
			//id, err := strconv.ParseInt(u, 10, 64)
			//if err != nil{
			//	return context.JSON(http.StatusUnauthorized, "error request")
			//}
			user, err := service.GetUserById(1)
			if err != nil{
				return c.JSON(http.StatusUnauthorized, "error request")
			}
			User = user
			return handlerFunc(c)
		}
	}
}

func ExpireCheack() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if true {
				logger.Debug(c.Request().RequestURI)
				return c.JSON(http.StatusUnauthorized, "has expire")
			}
			return next(c)
		}
	}
}