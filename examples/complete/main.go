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
	"github.com/xackery/egui/element/label"
	"golang.org/x/image/colornames"
)

var (
	lblHello         *label.Element
	screenResolution = image.Point{X: 320, Y: 240}
)

func main() {
	ui, err := egui.NewUI(screenResolution, 1)
	if err != nil {
		fmt.Println("failed to start ui:", err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())

	f, err := os.Open("ui-base-01.png")
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

	f, err = os.Open("ui.json")
	if err != nil {
		fmt.Println("failed", err)
		return
	}
	defer f.Close()
	r := aseprite.NewReader(f)
	slices, err := r.ReadSlices()
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

	f, err = os.Open("reaper_blade_1.png")
	if err != nil {
		fmt.Println("failed", err)
		return
	}
	defer f.Close()
	_, err = ui.NewImage("reaper", f, ebiten.FilterDefault)
	if err != nil {
		fmt.Println("failed", err)
		return
	}

	reaper, err := ui.NewSprite("reaper", "global", 0, 0, "reaper")
	if err != nil {
		fmt.Println("failed to create sprReaper", err.Error())
		return
	}

	btnChange, err := ui.NewButton("btnTest", "global", "Change Direction", 50, 50, 150, 30, color.White, "btnPress", "btnUnpress")
	if err != nil {
		fmt.Println("failed to create btnTest", err.Error())
		return
	}
	lastDirection := 0
	btnChange.SetOnPressFunction(func() {
		directions := []string{"down", "left", "up", "right"}
		lastDirection++
		if len(directions) < lastDirection+1 {
			lastDirection = 0
		}

		reaper.SetAnimationName(directions[lastDirection])
	})

	lblHello, err = ui.NewLabel("lblHello", "global", "Hello", 100, 100, colornames.Yellow)
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

	x := float64(rand.Intn(screenResolution.X - int(lblHello.Width())))
	y := float64(rand.Intn(screenResolution.Y - int(lblHello.Height())))
	lblHello.LerpPosition(common.Vect(x, y), 3*time.Second, false, randomBounce)
}
