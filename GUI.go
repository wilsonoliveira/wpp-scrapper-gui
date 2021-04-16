package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/lxn/win"
	wppscrapper "github.com/ribeiroferreiralucas/wpp-scrapper"
	"github.com/ribeiroferreiralucas/wpp-scrapper/wppscrapperimp"
)

var App = app.New()
var window = App.NewWindow("GUI")
var scrapper = wppscrapperimp.InitializeConnection()

var image canvas.Image

func loadImage(path string) fyne.CanvasObject {
	image.File = path
	image.FillMode = canvas.ImageFillContain
	image.Refresh()
	return &image

}

type listItem struct {
	grupo  string
	status int
}

type EventLogger struct {
}

func (logger *EventLogger) OnWppScrapperChatScrapStarted(chat wppscrapper.Chat) {
	fmt.Println("logger OnWppScrapperChatScrapStarted")
}

func (logger *EventLogger) OnWppScrapperChatScrapFinished(chat wppscrapper.Chat) {
	fmt.Println("logger OnWppScrapperChatScrapFinished")
}

func getwindowSize() (float32, float32) {
	return float32(win.GetSystemMetrics(win.SM_CXSCREEN)) * 0.5, float32(win.GetSystemMetrics(win.SM_CYSCREEN)) * 0.5

}

func makeList() fyne.CanvasObject {
	// images := []string{`All-Might-1.jpg`, `bakugou.png`}

	list := widget.NewList(func() int { return 1 },
		func() fyne.CanvasObject {
			btn := widget.NewButton("Coletar", func() {
				log.Println("tapped")

				if !scrapper.Initialized() {
					<-scrapper.WaitInitialization()
				}

				for k, v := range scrapper.GetChats() {
					// chatList = append(chatList, v.Name())
					fmt.Println("k:", k, "v:", v, "Status:", v.GetStatus())
				}

				fmt.Println("---------------\n\n\n\nSTART SCRAPPER\n\n\n\n----------------")
				scrapper.StartScrapper(true)
				// <-time.After(10000 * time.Second)

			})
			btn2 := widget.NewButton("Parar", func() {
				log.Println("tapped2")
				scrapper.StopScrapper()
				fmt.Println("---------------\n\n\n\nSTOP SCRAPPER\n\n\n\n----------------")

				for k, v := range scrapper.GetChats() {
					// chatList = append(chatList, v.Name())
					fmt.Println("k:", k, "v:", v, "Status:", v.GetStatus())
				}

			})

			// label := widget.NewLabel("view 1")
			return container.NewVBox(btn, btn2)
		}, func(lii widget.ListItemID, co fyne.CanvasObject) {
			/* func(i int, c fyne.CanvasObject) {
				cont := c.(*fyne.Container)
				label := cont.Objects[0].(*widget.Label)
				label.SetText(fmt.Sprintf("view %v", i+1))
			} */
		})

	// list.OnSelected = func(id int) {
	// 	image.File = images[id]
	// 	image.Refresh()
	// }

	return list
}

func makeUI() fyne.CanvasObject {

	return container.NewBorder(nil, nil, makeList(), nil, groupList())
}

func main() {

	screen_width, screen_height := getwindowSize()

	btn := widget.NewButton("Start", func() {
		log.Println("tapped")
		cont := container.NewBorder(nil, nil, nil, nil, loadImage(""))
		window.SetContent(cont)
		tst()
	})

	// window.SetContent(makeUI())

	window.SetContent(container.NewBorder(nil, nil, nil, nil, btn))

	window.Resize(fyne.NewSize(screen_width, screen_height))

	// window.SetContent(container.NewBorder(nil, nil, nil, nil, loadImage("")))

	window.ShowAndRun()

}

func groupList() fyne.CanvasObject {
	chatList := make([]listItem, 0)
	for _, v := range scrapper.GetChats() {
		item := listItem{grupo: v.Name(), status: -1}
		chatList = append(chatList, item)
		// fmt.Println("k:", k, "v:", v, "Status:", v.GetStatus())
	}

	evtLogger := &EventLogger{}

	evtHandler := scrapper.GetWppScrapperEventHandler()
	fmt.Println("evtHandler")
	fmt.Println(evtHandler)

	evtHandler.AddOnChatScrapStartedListener(evtLogger)
	evtHandler.AddOnChatScrapFinishedListener(evtLogger)

	list := widget.NewList(
		func() int {
			return len(chatList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("template"), widget.NewLabel("sts"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			labelGrupo := box.Objects[0].(*widget.Label)
			labelGrupo.SetText(chatList[i].grupo)
			labelstatus := box.Objects[1].(*widget.Label)
			labelstatus.SetText(strconv.Itoa(chatList[i].status))

		})

	return list

}

func tst() {

	qrc := make(chan string)
	go func() {
		qrCode, _ := qr.Encode(<-qrc, qr.M, qr.Auto)
		fmt.Println("QRfunc")
		// Scale the barcode to 200x200 pixels
		qrCode, _ = barcode.Scale(qrCode, 200, 200)

		file, _ := os.Create("qrcode.png")
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, qrCode)

		loadImage("qrcode.png")

	}()

	_, err := scrapper.ReAuth(qrc, "other")
	if err != nil {
		log.Fatalf("error scrapper.ReAuth in: %v\n", err)
	}

	if !scrapper.Initialized() {
		<-scrapper.WaitInitialization()
		window.SetContent(makeUI())
		// <-time.After(2 * time.Second)

	}

}

func verifyChatStatus() {
	if !scrapper.Initialized() {
		for {

		}
	}
}
