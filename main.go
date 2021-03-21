package main

import (
	"gin-zipfile/controller"
	"gin-zipfile/tool"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		panic(err.Error())
	}

	engine := gin.Default()

	registerRouter(engine)

	initExtra(engine)

	engine.Run(cfg.AppHost + ":" + cfg.AppPort)
}

// 路由注册
func registerRouter(route *gin.Engine) {
	new(controller.FileController).Router(route)
}

// 额外初始化
func initExtra(route *gin.Engine) {

	route.LoadHTMLGlob("view/**/*")

	route.Static("/file/","./public/file")
}
