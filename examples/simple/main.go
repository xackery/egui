package main

import (
	"fmt"
	"image"
	"os"

	"github.com/pkg/errors"
	"github.com/xackery/egui"
)

func main() {
	err := run()
	if err != nil {
		fmt.Println("failed to run", err.Error())
		os.Exit(1)
	}
}

func run() error {
	ui, err := egui.NewUI(image.Point{X: 640, Y: 480}, 1)
	if err != nil {
		return errors.Wrap(err, "start ui")
	}
	fmt.Println(ui)
	return nil
}
