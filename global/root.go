package global

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

/*
用于推断当前项目路径
*/

var (
	RootDir string           //项目根路径
	once    = new(sync.Once) //确保并发程序中某段代码只被执行一次
)

// 判断文件是否存在
func exist(filePath string) bool {
	_, err := os.Stat(filePath)                      //获取文件的基本信息
	return err == nil || errors.Is(err, os.ErrExist) //os.ErrExist 表示文件或目录已经存在的错误
}

// 计算项目路径
func inferRootDir() string {
	cwd, err := os.Getwd() //获取当前工作目录的路径
	if err != nil {
		panic(err)
	}
	//通过本项目根目录下的子目录 /config 来对根目录进行定位
	var infer func(string) string      //定义函数变量
	infer = func(path string) string { // 给函数变量赋值
		if exist(path + "/config") {
			return path
		}
		return infer(filepath.Dir(path))
	}
	return infer(cwd) //一定要调用 infer 函数
}

// 项目跟目录初始化函数
func init() {
	once.Do(func() { //确保查找项目根目录只查找一次
		RootDir = inferRootDir()
	})
}
