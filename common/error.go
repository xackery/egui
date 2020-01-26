package common

import (
	"fmt"
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
)
