package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type FileController struct {

}

func (file *FileController) Router (engine *gin.Engine) {
	engine.GET("/index",file.Index)
}

func (file *FileController) Index (context *gin.Context) {

	list,err := file.GetAllDir("public/file/")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(list)

	context.HTML(http.StatusOK,"index.html",list)
}

func (file *FileController) GetAllDir (pathname string) ([]string,error) {

	var tmpDir []string

	list,err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return tmpDir,err
	}

	for _,fi := range list {

		if fi.IsDir() {
			tmpDir = append(tmpDir,fi.Name())
		}
	}

	return tmpDir,nil
}
