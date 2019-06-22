package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/xackery/egui"
	"github.com/xackery/egui/aseprite"
	"github.com/xackery/egui/common"
	"golang.org/x/image/colornames"
)

var (
	lblHello         egui.Interfacer
	screenResolution = image.Point{X: 320, Y: 240}
)

func main() {

	ui, err := egui.NewUI(screenResolution, 1)
	if err != nil {
		fmt.Println("failed to start ui:", err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())

	f, err := os.Open("ui.png")
	if err != nil {
		fmt.Println("failed to open ui.png", err.Error())
		return
	}

	img, err := ui.NewImage("ui", f, ebiten.FilterDefault)
	if err != nil {
		fmt.Println("failed to add newImage", err.Error())
		f.Close()
		return
	}
	f.Close()

	f, err = os.Open("ui.aseprite-data")
	if err != nil {
		fmt.Println("failed", err)
		return
	}
	defer f.Close()
	r := aseprite.NewReader(f)
	slices, err := r.ReadAll()
	if err != nil {
		fmt.Println("failed read", err)
		return
	}
	for _, slice := range slices {
		err = img.AddSlice(slice)
		if err != nil {
			fmt.Println("failed to add slice", slice.Name, err)
			return
		}
	}

	btnHello, err := ui.NewButton("btnTest", "global", "Hello!", common.Rect(50, 50, 100, 80), color.White, "ui", "btnPress", "btnUnpress")
	if err != nil {
		fmt.Println("failed to create btnTest", err.Error())
		return
	}
	btnHello.SetOnPressFunction(func() {
		fmt.Println("pressed", btnHello.Name())
	})

	lblHello, err = ui.NewLabel("lblHello", "Hello, World!", common.Rect(100, 100, 100, 20), colornames.Yellow)
	if err != nil {
		fmt.Println("failed to create lblHello", err.Error())
		return
	}
	randomBounce()
	err = ebiten.Run(ui.Update, screenResolution.X, screenResolution.Y, 2, "Complete Example")
	if err != nil {
		fmt.Println("failed to run")
		return
	}
}

func randomBounce() {
	x := float64(rand.Intn(screenResolution.X - int(lblHello.Shape().Max.X)))
	y := float64(rand.Intn(screenResolution.Y - int(lblHello.Shape().Max.Y)))
	lblHello.LerpPosition(common.Vect(x, y), 3*time.Second, false, randomBounce)
}
