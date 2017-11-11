package main

import (
	"github.com/beewit/spread-update/update"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	"os/exec"
)

type MyWindow struct {
	*walk.MainWindow
	label *Label
}

func main() {
	//mw := &MyWindow{label: &Label{}}
	//mw.MainWindow = new(walk.MainWindow)
	//if err := walk.InitWindow(
	//	mw,
	//	nil,
	//	`\o/ Walk_MainWindow_Class \o/`,
	//	win.WS_CAPTION|win.WS_SYSMENU|win.WS_THICKFRAME, 0); err != nil {
	//	update.Log.Error(err.Error())
	//	return
	//}
	//
	//succeeded := false
	//defer func() {
	//	if !succeeded {
	//		mw.Dispose()
	//	}
	//}()
	//
	//mw.SetPersistent(true)
	//icon, err := walk.NewIconFromFile("favicon.ico")
	//if err != nil {
	//}
	//mw.label.Text = "你好啊"
	//m := MainWindow{
	//	Icon:     icon,
	//	AssignTo: &mw.MainWindow,
	//	Title:    "工蜂小智 - 更新程序",
	//	MinSize:  Size{300, 100},
	//	Size:     Size{300, 100},
	//	Layout:   VBox{MarginsZero: true},
	//	Children: []Widget{
	//		mw.label,
	//		PushButton{
	//			Text: "Copy",
	//			OnClicked: func() {
	//				println(fmt.Sprintf("%v", mw.label.Text))
	//
	//				mw.label.Text = "你是王武"
	//			},
	//		},
	//	},
	//}
	//go func() {
	//	println("已经暂停了？")
	//	time.Sleep(time.Second * 3)
	//	mw.label.Text = fmt.Sprintf("你是王武：%v", mw.Handle())
	//	println(fmt.Sprintf("你是王武：%v", mw.Handle()))
	//
	//	flog := win.SetForegroundWindow(mw.Handle())
	//	println(flog)
	//}()
	//if _, err := m.Run(); err != nil {
	//	update.Log.Error(err.Error())
	//}
	//return
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
