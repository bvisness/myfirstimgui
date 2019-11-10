package imgui

import (
	"image"
	"image/color"

	"github.com/bvisness/myfirstimgui/rectutil"
)

type ButtonResult struct {
	Clicked bool
	Size    image.Point
}

// TODO: Make this take a rect
func (ui *UIContext) Button(id, text string, pos, size image.Point, c color.RGBA) ButtonResult {
	me := UIID{
		Name: id,
	}
	result := false

	r := rectutil.SizeRect(pos, size)

	if rectutil.PointInRect(ui.Mouse.Pos, r) {
		ui.SetHot(me)
	}

	if ui.IsActive(me) {
		if ui.IsMouseUpThisFrame() {
			if ui.IsHot(me) {
				result = true
			}
			ui.SetNoneActive()
		}
	} else if ui.IsHot(me) {
		if ui.IsMouseDownThisFrame() {
			ui.SetActive(me)
		}
	}

	if ui.IsActive(me) {
		c = AlphaOver(c, color.RGBA{255, 255, 255, 50})
	} else if ui.IsHot(me) {
		c = AlphaOver(c, color.RGBA{255, 255, 255, 100})
	}
	ui.Img.DrawRect(r, c)

	return ButtonResult{
		Clicked: result,
		Size:    r.Size(),
	}
}
