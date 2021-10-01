package gui

import (
	"go-pong/game"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var window fyne.Window

func Init(a fyne.App) error {
	window = a.NewWindow("go-pong")
	window.Resize(fyne.NewSize(500, 500))
	window.SetFixedSize(true)
	mainContainer := createMainContainer()
	window.SetContent(mainContainer)
	window.ShowAndRun()
	return nil
}
func createMainContainer() *container.AppTabs {
	offlineItem := container.NewTabItem("Offline", createOfflineCanvas())
	onlineItem := container.NewTabItem("LAN", widget.NewLabel("Coming soon"))
	tabs := container.NewAppTabs(offlineItem, onlineItem)
	tabs.SetTabLocation(container.TabLocationLeading)
	return tabs
}
func createOfflineCanvas() fyne.CanvasObject {
	vContainer := container.NewVBox()
	runButton := widget.NewButton("Run", func() { game.Run(window, false) })
	vContainer.Add(runButton)
	return vContainer
}
