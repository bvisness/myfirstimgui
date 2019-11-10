package imgui

import (
	"image"
	"image/color"
)

type Direction int

const (
	Up Direction = iota + 1
	Down
	Left
	Right
)

var (
	ColorBarBackground = color.RGBA{200, 200, 200, 50}
	ColorWidget        = color.RGBA{200, 200, 200, 100}
)

type UIID struct {
	Name string
}

func (id UIID) String() string {
	return id.Name
}

type UIMouse struct {
	PosPrevious image.Point
	Pos         image.Point

	IsMouseDownPrevious bool
	IsMouseDown         bool
}

type UIStyle struct {
	Spacing    int
	WidgetSize int
	// in future, some kind of overall scaling factor for interactive things?
}

type UIContext struct {
	Hot    *UIID
	Active *UIID

	Style UIStyle
	Mouse UIMouse

	Img UIImage

	ElementState map[UIID]interface{}
}

func (ui *UIContext) WithClip(clip image.Rectangle, f func()) {
	originalClipArea := ui.Img.ClipRect
	defer func() {
		ui.Img.ClipRect = originalClipArea
	}()

	ui.Img.ClipRect = clip
	f()
}

func (ui *UIContext) IsHot(obj UIID) bool {
	if ui.Hot == nil {
		return false
	}

	return *ui.Hot == obj
}

func (ui *UIContext) IsActive(obj UIID) bool {
	if ui.Active == nil {
		return false
	}

	return *ui.Active == obj
}

func (ui *UIContext) SetHot(obj UIID) {
	if ui.Active == nil {
		ui.Hot = &obj
	}
}

func (ui *UIContext) SetActive(obj UIID) {
	ui.Active = &obj
}

func (ui *UIContext) SetNoneActive() {
	ui.Active = nil
}

func (ui *UIContext) IsMouseDownThisFrame() bool {
	return !ui.Mouse.IsMouseDownPrevious && ui.Mouse.IsMouseDown
}

func (ui *UIContext) IsMouseUpThisFrame() bool {
	return ui.Mouse.IsMouseDownPrevious && !ui.Mouse.IsMouseDown
}

func (ui *UIContext) MouseDelta() image.Point {
	return image.Pt(ui.Mouse.Pos.X-ui.Mouse.PosPrevious.X, ui.Mouse.Pos.Y-ui.Mouse.PosPrevious.Y)
}
