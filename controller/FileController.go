package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type FileController struct {

}

func (file *FileController) Router (engine *gin.Engine) {
	engine.GET("/index",file.Index)
}

func (file *FileController) Index (context *gin.Context) {
	context.JSON(http.StatusOK,gin.H{
		"status" : 1,
		"message" : "请求成功",
	})
}
