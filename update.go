package main

import (
	"time"

	"github.com/beewit/spread-update/update"
	"fmt"
	"os"
)

func main() {
	defer func() {
		update.Logs("结束更新程序...")
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("《程序出现严重错误，终止运行！》，ERROR：%v", err)
			update.Logs(errStr)
		}
	}()
	update.Logs("启动更新程序...")
	_, err := update.Update(update.Version{Major: 1, Minor: 0, Patch: 0})
	if err != nil {
		update.Logs(fmt.Sprintf("Update 更新失败：%s", err.Error()))
		return
	}
	time.Sleep(time.Millisecond * 500)
	_, err = os.StartProcess("spread.exe", nil, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
	if err != nil {
		update.Logs(fmt.Sprintf("CallEXE 启动失败：%s", err.Error()))
		return
	}
	update.Logs("更新程序完毕")
}