package main

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/xackery/egui"
	"github.com/xackery/egui/common"
	"golang.org/x/image/colornames"
)

var (
	lblHello         egui.Interfacer
	screenResolution = image.Point{X: 320, Y: 240}
)

func main() {

	ui, err := egui.NewUI(screenResolution)
	if err != nil {
		fmt.Println("failed to start ui:", err.Error())
		return
	}
	lblHello, err = ui.NewLabel("lblHello", "Hello, World!", common.Rect(100, 100, 100, 20), colornames.Yellow)
	if err != nil {
		fmt.Println("failed to create lblHello", err.Error())
		return
	}
	rand.Seed(time.Now().UnixNano())

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
