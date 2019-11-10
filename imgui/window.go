package imgui

import (
	"image"
	"image/color"

	"github.com/bvisness/myfirstimgui/rectutil"
)

type windowState struct {
	Open         bool
	Pos          image.Point
	Size         image.Point
	ActiveWidget int
}

func (ui *UIContext) Window(id string, initialPos, initialSize image.Point, forceNewState bool) (bool, image.Rectangle) {
	widgetToggle := 1
	widgetResize := 2
	widgetTitleBar := 3

	me := UIID{Name: id}

	istate, stateExists := ui.ElementState[me]
	if !stateExists || forceNewState {
		istate = windowState{
			Open: true,
			Pos:  initialPos,
			Size: initialSize,
		}
	}
	state := istate.(windowState)
	defer func() {
		ui.ElementState[me] = state
	}()

	windowRect := rectutil.SizeRect(state.Pos, state.Size)
	if !state.Open {
		windowRect = rectutil.SizeRect(state.Pos, image.Pt(state.Size.X, ui.Style.WidgetSize))
	}
	titleBarRect := rectutil.SizeRect(image.Pt(state.Pos.X+ui.Style.WidgetSize, state.Pos.Y), image.Pt(state.Size.X-ui.Style.WidgetSize, ui.Style.WidgetSize))
	toggleRect := rectutil.SizeRect(state.Pos, image.Pt(ui.Style.WidgetSize, ui.Style.WidgetSize))
	resizeRect := image.Rectangle{windowRect.Max.Sub(image.Pt(ui.Style.WidgetSize, ui.Style.WidgetSize)), windowRect.Max}

	if ui.IsActive(me) {
		switch state.ActiveWidget {
		case widgetTitleBar:
			state.Pos = state.Pos.Add(ui.MouseDelta())
		case widgetResize:
			sizeDelta := ui.MouseDelta()
			if !state.Open {
				sizeDelta.Y = 0
			}

			state.Size = rectutil.MaxPoint(state.Size.Add(sizeDelta), image.Pt(60, 40))
		case widgetToggle:
			if ui.IsMouseUpThisFrame() && ui.IsHot(me) {
				state.Open = !state.Open
			}
		}

		if ui.IsMouseUpThisFrame() {
			ui.SetNoneActive()
		}
	} else if ui.IsHot(me) {
		if ui.IsMouseDownThisFrame() {
			ui.SetActive(me)
			if rectutil.PointInRect(ui.Mouse.Pos, toggleRect) {
				state.ActiveWidget = widgetToggle
			} else if rectutil.PointInRect(ui.Mouse.Pos, resizeRect) {
				state.ActiveWidget = widgetResize
			} else if rectutil.PointInRect(ui.Mouse.Pos, titleBarRect) {
				state.ActiveWidget = widgetTitleBar
			} else {
				state.ActiveWidget = 0
			}
		}
	}

	if rectutil.PointInRect(ui.Mouse.Pos, windowRect) {
		ui.SetHot(me)
	}

	ui.Img.DrawRect(windowRect, color.RGBA{0, 0, 0, 200})
	ui.Img.DrawRect(titleBarRect, ColorBarBackground)
	ui.Img.DrawRect(toggleRect, ColorWidget)
	ui.Img.DrawHalfRect(resizeRect, ColorWidget, LowerRight)

	contentRect := image.Rect(
		windowRect.Min.X+ui.Style.Spacing,
		windowRect.Min.Y+2*ui.Style.Spacing,
		windowRect.Max.X-ui.Style.Spacing,
		windowRect.Max.Y-ui.Style.Spacing,
	)

	return state.Open, contentRect
}
