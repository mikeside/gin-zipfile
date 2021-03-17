package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
)

type FileController struct {
}

func (file *FileController) Router(engine *gin.Engine) {
	engine.GET("/index", file.Index)
}

func (file *FileController) Index(context *gin.Context) {

	list, err := file.GetAllDir("public/file")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(list)

	context.HTML(http.StatusOK, "index.html", gin.H{
		"status": 1,
		"data":   list,
	})
}

func (file *FileController) GetAllDir(pathname string) ([]map[string]string, error) {

	var tmpDir []map[string]string

	list, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return tmpDir, err
	}

	for _, fi := range list {

		if fi.IsDir() {

			demoList, _ := ioutil.ReadDir(pathname + "/" + fi.Name())
			newTmpDir := map[string]string{
				"name":  fi.Name(),
				"count": strconv.Itoa(len(demoList)),
			}

			tmpDir = append(tmpDir, newTmpDir)
		}
	}

	return tmpDir, nil
}
