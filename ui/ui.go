package ui

import (
	"fmt"
	"image"

	// Used for decoding
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

var (
	// ErrElementNameInvalid is returned when a element name has invalid characters or too short
	ErrElementNameInvalid = fmt.Errorf("element name invalid")
	// ErrElementAlreadyExists is returned when a element already exists
	ErrElementAlreadyExists = fmt.Errorf("element already exists")
	// ErrElementNotFound is returned when a element is not loaded into the UI
	ErrElementNotFound = fmt.Errorf("element not found")
	// ErrResourceNameInvalid is returned when a resource name has invalid characters or too short
	ErrResourceNameInvalid = fmt.Errorf("resource name invalid")
	// ErrResourceAlreadyExists is returned when a resource already exists
	ErrResourceAlreadyExists = fmt.Errorf("resource already exists")
	// ErrResourceNotFound is returned when a resource was not found
	ErrResourceNotFound = fmt.Errorf("resource not found")
	// ErrSceneNameInvalid is returned when a scene name has invalid characters or too short
	ErrSceneNameInvalid = fmt.Errorf("scene name invalid")
	// ErrSceneAlreadyExists is returned when a Scene already exists
	ErrSceneAlreadyExists = fmt.Errorf("scene already exists")
	// ErrSceneNotFound is returned when a scene is not loaded into the UI
	ErrSceneNotFound = fmt.Errorf("scene not found")
)

// UI contains core game components
type UI struct {
	font             font.Face
	fontMHeight      int
	scenes           map[string]*Scene
	currentScene     *Scene
	globalScene      *Scene
	currentMap       string
	screenResolution image.Point
	images           map[string]*ebiten.Image
}

// NewUI instantiates a new User Interface
func NewUI(font font.Face, fontMHeight int, screenResolution image.Point) (*UI, error) {
	u := &UI{
		font:             font,
		fontMHeight:      fontMHeight,
		scenes:           make(map[string]*Scene),
		images:           make(map[string]*ebiten.Image),
		screenResolution: screenResolution,
	}
	gs := NewScene()
	u.AddScene("global", gs)
	u.globalScene = gs
	return u, nil
}

// SetResolution changes the resolution of the UI
func (u *UI) SetResolution(resolution image.Point) {
	u.screenResolution = resolution
	for _, scene := range u.scenes {
		scene.onResolutionChange(resolution)
	}
}

// Resolution returns the resolution
func (u *UI) Resolution() image.Point {
	return u.screenResolution
}

// AddResource adds an image resource to ui
func (u *UI) AddResource(name string, img *ebiten.Image) error {
	_, ok := u.images[name]
	if !ok {
		return ErrResourceAlreadyExists
	}
	u.images[name] = img
	return nil
}

// Resource returns named resource
func (u *UI) Resource(name string) (*ebiten.Image, error) {
	img, ok := u.images[name]
	if !ok {
		return nil, ErrResourceNotFound
	}
	return img, nil
}

// Update updates all UI elements
func (u *UI) Update(dt float64) {
	if u.globalScene != nil {
		u.globalScene.update(dt)
	}
	if u.currentScene != nil {
		u.currentScene.update(dt)
	}
}

// Draw renders all UI elements
func (u *UI) Draw(screen *ebiten.Image) {
	if u.currentScene != nil {
		u.currentScene.draw(screen)
	}
	if u.globalScene != nil {
		u.globalScene.draw(screen)
	}
}

// SetCurrentScene sets current scene
func (u *UI) SetCurrentScene(name string) error {
	s, ok := u.scenes[name]
	if !ok {
		return ErrSceneNotFound
	}
	u.currentScene = s
	return nil
}

// CurrentScene returns the current scene
func (u *UI) CurrentScene() *Scene {
	return u.currentScene
}

// AddScene appends a new scene to the UI
func (u *UI) AddScene(name string, scene *Scene) error {
	_, ok := u.scenes[name]
	if ok {
		return ErrSceneAlreadyExists
	}
	scene.onResolutionChange(u.screenResolution)
	u.scenes[name] = scene
	return nil
}

// Scene gets a scene
func (u *UI) Scene(name string) (*Scene, error) {
	scene, ok := u.scenes[name]
	if !ok {
		return nil, ErrSceneNotFound
	}
	return scene, nil
}
