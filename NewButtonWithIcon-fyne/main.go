package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Button Widget")

	// content := widget.NewButton("click me", func() {
	// 	log.Println("tapped")
	// })

	content := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		log.Println("tapped home")
	})

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

/*

sudo apt-get install  -y gcc libgl1-mesa-dev xorg-dev
go install fyne.io/fyne/v2/cmd/fyne@latest

go get fyne.io/fyne/v2@latest
go mod tidy -x
go run .
go run -x .
go build -x
./app
*/
