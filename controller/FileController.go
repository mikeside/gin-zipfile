package controller

import (
	"fmt"
	"gin-zipfile/tool"
	"github.com/gin-gonic/gin"
	"github.com/gookit/validate"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	filePath,_ = tool.Mkdir("public/file")
	tempPath,_ = tool.Mkdir("public/temp")
)

type FileController struct {
}

type FileCreate struct {
	Name string `validate:"required|maxLen:50" message:"required:名称不能为空|maxLen:名称输入限制50字符以内"`
}

type UploadDemo struct {
	Name string                `validate:"required" message:"required:项目名称参数错误"`
	Demo *multipart.FileHeader `validate:"required" message:"required:请选择要上传的文件"`
}

func (file *FileController) Router(engine *gin.Engine) {

	engine.GET("/", file.Index)

	engine.POST("/file/create", file.Create)

	engine.POST("/file/del", file.Delete)

	engine.POST("/file/upload/demo", file.UploadDemo)
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

func (file *FileController) UploadDemo(context *gin.Context) {

	name := context.PostForm("name")
	demo, _ := context.FormFile("demo")

	v := validate.Struct(&UploadDemo{
		Name: name,
		Demo: demo,
	})
	if v.Validate() {

		if filepath.Ext(demo.Filename) != ".zip" {
			context.JSON(http.StatusOK, tool.Err.WithMsg("请上传zip格式的文件"))
			return
		}

		// 执行上传
		tempDemoPath := tempPath + "/" + demo.Filename
		if err := context.SaveUploadedFile(demo, tempDemoPath); err != nil {
			fmt.Println(err)
			context.JSON(http.StatusOK, tool.Err.WithMsg("上传失败"))
			return
		}

		// 执行解压
		projectPath := filePath + "/" + name
		if err := file.Unzip(tempDemoPath,projectPath); err != nil {
			fmt.Println(err)
			context.JSON(http.StatusOK, tool.Err.WithMsg("解压失败"))
			return
		}

		context.JSON(http.StatusOK, tool.Ok.WithMsg("上传成功"))

	} else {
		context.JSON(http.StatusOK, tool.Err.WithMsg(v.Errors.One()))
	}
}

func (file *FileController) Unzip(zipPath string, destPath string) error {

	if err := tool.New(zipPath, tempPath).Extract(); err != nil {
		return err
	}

	// 清空旧demo文件
	os.RemoveAll(destPath)

	// 移动文件和删除zip包
	oldDemoPath := tempPath + "/" + strings.TrimSuffix(path.Base(zipPath), path.Ext(zipPath))
	if err := os.Rename(oldDemoPath, destPath); err != nil {
		return err
	}
	os.Remove(zipPath)

	return nil
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
