package game

import "fyne.io/fyne/v2"

var isOnline bool

func Run(w fyne.Window, on bool) {
	isOnline = on
	w.Hide()

}
