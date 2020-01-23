package egui

import (
	"fmt"
	"image"
	"time"

	// Used for decoding
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/text/language"
)

var (
	// ErrElementNameInvalid is returned when a element name has invalid characters or too short
	ErrElementNameInvalid = fmt.Errorf("element name invalid")
	// ErrElementAlreadyExists is returned when a element already exists
	ErrElementAlreadyExists = fmt.Errorf("element already exists")
	// ErrElementNotFound is returned when a element is not loaded into the UI
	ErrElementNotFound = fmt.Errorf("element not found")
	// ErrFontNameInvalid is returned when a font name has invalid characters or too short
	ErrFontNameInvalid = fmt.Errorf("font name invalid")
	// ErrFontAlreadyExists is returned when a font already exists
	ErrFontAlreadyExists = fmt.Errorf("font already exists")
	// ErrFontNotFound is returned when a font was not found
	ErrFontNotFound = fmt.Errorf("font not found")
	// ErrFontCannotRemoveDefault is returned when you attempt to delete a font currently set as default
	ErrFontCannotRemoveDefault = fmt.Errorf("font is default, cannot remove")
	// ErrImageNameInvalid is returned when a image name has invalid characters or too short
	ErrImageNameInvalid = fmt.Errorf("image name invalid")
	// ErrImageAlreadyExists is returned when a image already exists
	ErrImageAlreadyExists = fmt.Errorf("image already exists")
	// ErrImageNotFound is returned when a image was not found
	ErrImageNotFound = fmt.Errorf("image not found")
	// ErrSceneNameInvalid is returned when a scene name has invalid characters or too short
	ErrSceneNameInvalid = fmt.Errorf("scene name invalid")
	// ErrSceneAlreadyExists is returned when a Scene already exists
	ErrSceneAlreadyExists = fmt.Errorf("scene already exists")
	// ErrSceneNotFound is returned when a scene is not loaded into the UI
	ErrSceneNotFound = fmt.Errorf("scene not found")

	op = &ebiten.DrawImageOptions{}
)

// UI contains core game components
type UI struct {
	defaultFont      *Font
	scenes           map[string]*Scene
	currentScene     *Scene
	globalScene      *Scene
	screenResolution image.Point
	images           map[string]*Image
	fonts            map[string]*Font
	lastUpdate       time.Time
	tileScale        float64
	textScale        float64
	defaultLanguage  language.Tag
}

// NewUI instantiates a new User Interface
func NewUI(screenResolution image.Point, scale float64) (*UI, error) {
	u := &UI{
		scenes:           make(map[string]*Scene),
		images:           make(map[string]*Image),
		fonts:            make(map[string]*Font),
		screenResolution: screenResolution,
		tileScale:        1,
		textScale:        1,
		defaultLanguage:  language.AmericanEnglish,
	}
	gs, err := u.NewScene("global")
	if err != nil {
		return nil, err
	}
	u.globalScene = gs
	u.currentScene = gs
	u.lastUpdate = time.Now()
	u.defaultFont, err = u.NewFontTTF("goregular", goregular.TTF, nil, 'M')
	if err != nil {
		return nil, errors.Wrap(err, "goregular font")
	}
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

// AddImage adds an image image to ui
func (u *UI) AddImage(img *Image) error {
	if img.Name() == "" {
		return ErrImageNameInvalid
	}
	_, ok := u.images[img.Name()]
	if ok {
		return ErrImageAlreadyExists
	}
	u.images[img.Name()] = img
	return nil
}

// Image returns a named image
func (u *UI) Image(name string) (*Image, error) {
	img, ok := u.images[name]
	if !ok {
		return nil, ErrImageNotFound
	}
	return img, nil
}

// Update updates all UI elements
func (u *UI) Update(image *ebiten.Image) error {
	dt := time.Since(u.lastUpdate).Seconds()

	u.onUpdate(dt)

	//graphical elements
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	u.Draw(image)
	return nil
}

func (u *UI) onUpdate(dt float64) {
	if u.globalScene != nil {
		u.globalScene.Update(dt)
	}
	if u.currentScene != nil {
		u.currentScene.Update(dt)
	}
}

// Draw renders all UI elements
func (u *UI) Draw(screen *ebiten.Image) {
	if u.currentScene != nil {
		u.currentScene.Draw(screen)
	}
	if u.globalScene != nil {
		u.globalScene.Draw(screen)
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

// SetDefaultFont updates all elements to use a new default font
func (u *UI) SetDefaultFont(name string) error {
	font, ok := u.fonts[name]
	if !ok {
		return ErrFontNotFound
	}
	u.defaultFont = font
	return nil
}

// DefaultFont returns the current default font
func (u *UI) DefaultFont(name string) error {
	font, ok := u.fonts[name]
	if !ok {
		return ErrFontNotFound
	}
	u.defaultFont = font
	return nil
}

// AddFont adds an font font to ui
func (u *UI) AddFont(font *Font) error {
	if font == nil {
		return ErrFontNameInvalid
	}
	_, ok := u.fonts[font.Name]
	if ok {
		return ErrFontAlreadyExists
	}
	u.fonts[font.Name] = font
	return nil
}

// Font returns named font
func (u *UI) Font(name string) (*Font, error) {
	img, ok := u.fonts[name]
	if !ok {
		return nil, ErrFontNotFound
	}
	return img, nil
}

// RemoveFont is used to unload and remove a font
func (u *UI) RemoveFont(name string) error {
	if len(name) < 1 {
		return ErrFontNameInvalid
	}
	_, ok := u.fonts[name]
	if !ok {
		return ErrFontNotFound
	}
	if u.defaultFont.Name == name {
		return ErrFontCannotRemoveDefault
	}
	delete(u.fonts, name)
	return nil
}
