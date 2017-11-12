package main

import (
	"os/exec"

	"github.com/beewit/spread-update/update"
)

func main() {
	update.Log.Info("启动更新程序...")
	err := update.Upload(update.Version{Major: 1, Minor: 0, Patch: 0})
	if err != nil {
		update.Log.Error(err.Error())
	}
	err = CallEXE("spread.exe")
	if err != nil {
		update.Log.Error(err.Error())
	}
	update.Log.Info("更新程序完毕")
}

func CallEXE(strGameName string) (err error) {
	cmd := exec.Command(strGameName)
	err = cmd.Start()
	if err != nil {
		return
	}
	return
}
