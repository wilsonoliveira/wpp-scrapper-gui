package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lxn/win"
)

var image canvas.Image

func loadImage(path string) fyne.CanvasObject {
	image.File = path
	image.FillMode = canvas.ImageFillContain
	image.Refresh()
	return &image

}

func getwindowSize() (float32, float32) {
	return float32(win.GetSystemMetrics(win.SM_CXSCREEN)) * 0.5, float32(win.GetSystemMetrics(win.SM_CYSCREEN)) * 0.5

}

func makeList() fyne.CanvasObject {
	images := []string{`All-Might-1.jpg`, `bakugou.png`}

	list := widget.NewList(func() int { return 2 },
		func() fyne.CanvasObject {
			label := widget.NewLabel("view 1")
			return container.NewHBox(label)
		},
		func(i int, c fyne.CanvasObject) {
			cont := c.(*fyne.Container)
			label := cont.Objects[0].(*widget.Label)
			label.SetText(fmt.Sprintf("view %v", i+1))
		})

	list.OnSelected = func(id int) {
		image.File = images[id]
		image.Refresh()
	}

	return list
}

func makeUI() fyne.CanvasObject {
	image.File = `All-Might-1.jpg`
	image.FillMode = canvas.ImageFillContain
	// content := widget.NewForm(
	// 	widget.NewFormItem("", ),
	// )

	return container.NewBorder(nil, nil, makeList(), nil, loadImage(""))
}

func qrCode() fyne.CanvasObject {
	qr := canvas.NewImageFromFile(`qr.png`)
	qr.FillMode = canvas.ImageFillOriginal

	return container.New(layout.NewCenterLayout(), qr)
}

func main() {
	App := app.New()
	window := App.NewWindow("GUI")

	screen_width, screen_height := getwindowSize()

	window.SetContent(makeUI())
	//window.SetContent(qrCode())
	window.Resize(fyne.NewSize(screen_width, screen_height))
	window.ShowAndRun()

}
