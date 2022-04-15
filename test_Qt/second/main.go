package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
)

func main() {
	// 创建应用程序
	app := widgets.NewQApplication(len(os.Args), os.Args)
	// 创建主窗口
	window := widgets.NewQWidget(nil, 0)
	// 设置窗口最小尺寸
	window.SetMinimumSize2(400, 200)
	// 设置标题
	window.SetWindowTitle("hello QT, hello GO")
	mainLayout := widgets.NewQVBoxLayout()
	// 创建垂直布局
	window.SetLayout(mainLayout)

	lineEidt := widgets.NewQLineEdit(nil)

	//创建一个label，用于存放go logo
	icon := gui.NewQPixmap3("go.jpg", "", core.Qt__AutoColor)
	iconLabel := widgets.NewQLabel(nil, 0)
	iconLabel.SetPixmap(icon)

	btn := widgets.NewQPushButton2("点我", nil)
	btn.ConnectClicked(func(bool) {
		widgets.QMessageBox_Information(nil, "我是对话框", "hello go，hello qt", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})

	// !!!! 注意是AddWidget 不是AddChildWidget
	mainLayout.AddWidget(lineEidt, 0, 0)
	mainLayout.AddWidget(iconLabel, 0, 0)
	window.Layout().AddWidget(btn)

	// 显示窗口
	window.Show()
	// 进入消息循环
	app.Exec()
}
