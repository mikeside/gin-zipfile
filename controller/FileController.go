package controller

import (
	"fmt"
	"gin-zipfile/tool"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	filePath = "public/file"
)

type FileController struct {
}

type FileCreate struct {
	Name string `validate:"required|maxLen:50" message:"required:名称不能为空|maxLen:名称输入限制50字符以内"`
}

func (file *FileController) Router(engine *gin.Engine) {
	engine.GET("/", file.Index)
	engine.POST("/file/create", file.Create)
	engine.POST("/file/del", file.Delete)
}

func (file *FileController) Index(context *gin.Context) {

	list, err := file.GetAllDir(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}

	//fmt.Println(list)

	context.HTML(http.StatusOK, "index.html", gin.H{
		"status": 1,
		"data":   list,
	})
}

func (file *FileController) Create(context *gin.Context) {

	// 验证数据
	name := context.PostForm("name")
	v := validate.Struct(&FileCreate{Name: name})
	if v.Validate() { // 验证成功
		projectPath := filePath + "/" + name
		if validate.PathExists(projectPath) {
			context.JSON(http.StatusOK, tool.Err.WithMsg("项目已存在"))
		} else {
			if err := os.MkdirAll(projectPath, os.ModePerm); err != nil {
				fmt.Println(err)
				context.JSON(http.StatusOK, tool.Err.WithMsg("新建项目失败"))
				return
			}

			context.JSON(http.StatusOK, tool.Ok.WithMsg("新建项目成功"))
		}
	} else {
		fmt.Println(v.Errors.One())
		context.JSON(http.StatusOK, tool.Err.WithMsg(v.Errors.One()))
	}
}

func (file *FileController) Delete(context *gin.Context) {

	name := context.PostForm("name")
	if name == "" {
		context.JSON(http.StatusOK, tool.Err.WithMsg("参数错误"))
		return
	}

	projectPath := filePath + "/" + name
	if err := os.RemoveAll(projectPath); err != nil {
		fmt.Println(err)
		context.JSON(http.StatusOK, tool.Err.WithMsg("删除失败"))
		return
	}

	context.JSON(http.StatusOK, tool.Ok.WithMsg("删除成功"))
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
