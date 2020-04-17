package router

import (
	"m/models"
	"m/models/db"
	"m/router/api"
	"m/router/middleware"

	"github.com/gin-gonic/gin"
)

func Server() {
	db.Initmysql()
	r := gin.Default()

	log := middleware.InitLog()
	log.Info("服务器开启")

	r.Use(middleware.Cors())
	r.Use(middleware.Logger())

	defer db.DEFAULTDB.Close()
	// 路由组
	ApiGroup := r.Group("") // 方便统一添加路由组前缀 多服务器上线使用
	UserRouter := ApiGroup.Group("user")
	{
		UserRouter.POST("regist", api.Regist)
		UserRouter.POST("login", api.Login)
		UserRouter.GET("get", api.GetUsers)
		UserRouter.PATCH("update", api.UpdateUser)
		UserRouter.POST("add", api.AddUsers)
	}
	test := ApiGroup.Group("test") //.Use(middleware.JWTAuth())
	{
		test.GET("hello", middleware.JWTAuth(), middleware.CheckPerm(models.PermDataRead), api.Test)
	}
	r.Run()
}
