package imgui

import (
	"image"
	"math"

	"github.com/bvisness/myfirstimgui/rectutil"

	"github.com/bvisness/myfirstimgui/imath"
)

type scrollPaneState struct {
	T float32
}

func (ui *UIContext) ScrollPane(id string, area image.Rectangle, contentHeight int, f func(contentArea image.Rectangle)) {
	me := UIID{Name: "scroll-" + id}

	istate, stateExists := ui.ElementState[me]
	if !stateExists {
		istate = scrollPaneState{
			T: 0,
		}
	}
	state := istate.(scrollPaneState)
	defer func() {
		ui.ElementState[me] = state
	}()

	handleHeight := 100
	handleMaxTravel := area.Size().Y - handleHeight
	contentMaxTravel := imath.Max(0, contentHeight-area.Size().Y)

	clipArea := image.Rectangle{area.Min, area.Max.Sub(image.Pt(ui.Style.WidgetSize, 0))}
	scrollbarArea := image.Rect(
		area.Max.X-ui.Style.WidgetSize,
		area.Min.Y,
		area.Max.X,
		area.Max.Y,
	)
	handleArea := rectutil.SizeRect(
		image.Pt(scrollbarArea.Min.X, scrollbarArea.Min.Y+imath.LerpF(0, handleMaxTravel, state.T)),
		image.Pt(ui.Style.WidgetSize, handleHeight),
	)

	if rectutil.PointInRect(ui.Mouse.Pos, scrollbarArea) {
		ui.SetHot(me)
	}

	if ui.IsActive(me) {
		state.T = float32(math.Max(0, math.Min(1,
			float64(state.T)+float64(ui.MouseDelta().Y)/float64(handleMaxTravel)),
		))

		if ui.IsMouseUpThisFrame() {
			ui.SetNoneActive()
		}
	} else if ui.IsHot(me) {
		if ui.IsMouseDownThisFrame() && rectutil.PointInRect(ui.Mouse.Pos, handleArea) {
			ui.SetActive(me)
		}
	}

	//ui.Img.DrawRect(clipArea, color.RGBA{255, 0, 0, 255})
	ui.Img.DrawRect(scrollbarArea, ColorBarBackground)
	ui.Img.DrawRect(handleArea, ColorWidget)

	ui.WithClip(clipArea, func() {
		f(rectutil.SizeRect(
			image.Pt(
				area.Min.X,
				area.Min.Y-int(state.T*float32(contentMaxTravel)),
			),
			image.Pt(clipArea.Size().X, contentHeight),
		))
	})
}
