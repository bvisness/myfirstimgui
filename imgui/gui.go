package imgui

import (
	"image"
	"image/color"

	"github.com/bvisness/myfirstimgui/imath"

	"github.com/bvisness/myfirstimgui/rectutil"
)

type Direction int

const (
	Up Direction = iota + 1
	Down
	Left
	Right
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
	Spacing int
	// in future: widget size, or some kind of overall scaling factor for interactive things
}

type UIContext struct {
	Hot    *UIID
	Active *UIID

	Style UIStyle
	Mouse UIMouse

	Img UIImage

	ElementState map[UIID]interface{}
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

type ButtonResult struct {
	Clicked bool
	Size    image.Point
}

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

type windowState struct {
	Open         bool
	Pos          image.Point
	Size         image.Point
	ActiveWidget int
}

func (ui *UIContext) Window(id string, initialPos, initialSize image.Point, forceNewState bool) (bool, image.Rectangle) {
	widgetSize := 16

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
		windowRect = rectutil.SizeRect(state.Pos, image.Pt(state.Size.X, widgetSize))
	}
	titleBarRect := rectutil.SizeRect(image.Pt(state.Pos.X+widgetSize, state.Pos.Y), image.Pt(state.Size.X-widgetSize, widgetSize))
	toggleRect := rectutil.SizeRect(state.Pos, image.Pt(widgetSize, widgetSize))
	resizeRect := image.Rectangle{windowRect.Max.Sub(image.Pt(widgetSize, widgetSize)), windowRect.Max}

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
	ui.Img.DrawRect(titleBarRect, color.RGBA{200, 200, 200, 50})
	ui.Img.DrawRect(toggleRect, color.RGBA{200, 200, 200, 100})
	ui.Img.DrawRect(resizeRect, color.RGBA{200, 200, 200, 100})

	contentRect := image.Rect(
		windowRect.Min.X+ui.Style.Spacing,
		windowRect.Min.Y+2*ui.Style.Spacing,
		windowRect.Max.X-ui.Style.Spacing,
		windowRect.Max.Y-ui.Style.Spacing,
	)

	return state.Open, contentRect
}

type ListLayout struct {
	Size image.Point

	ui         *UIContext
	itemPos    image.Point
	horizontal bool
}

func (ui *UIContext) ListLayout(startPos image.Point, spacing int, horizontal bool) ListLayout {
	return ListLayout{
		Size:       image.Pt(0, 0),
		ui:         ui,
		itemPos:    startPos,
		horizontal: horizontal,
	}
}

func (l *ListLayout) Item(f func(pos image.Point) image.Point) {
	resultSize := f(l.itemPos)

	if l.horizontal {
		l.itemPos = l.itemPos.Add(image.Pt(resultSize.X+l.ui.Style.Spacing, 0))
	} else {
		l.itemPos = l.itemPos.Add(image.Pt(0, resultSize.Y+l.ui.Style.Spacing))
	}

	if resultSize.X > l.Size.X {
		l.Size.X = resultSize.X
	}
	if resultSize.Y > l.Size.Y {
		l.Size.Y = resultSize.Y
	}
}

// TODO: The Up and Left directions force you to do this annoying subtraction dance. There has to be a way to make that easier.
type ListLayoutWithExcess struct {
	ui        *UIContext
	dir       Direction
	startPos  image.Point
	itemPos   image.Point
	totalSize image.Point
}

func (ui *UIContext) ListLayoutWithExcess(pos, totalSize image.Point, dir Direction) ListLayoutWithExcess {
	return ListLayoutWithExcess{
		ui:        ui,
		dir:       dir,
		startPos:  pos,
		itemPos:   pos,
		totalSize: totalSize,
	}
}

func (l *ListLayoutWithExcess) Item(f func(pos image.Point, crossAxisLength int) image.Point) {
	var crossAxisLength int
	switch l.dir {
	case Up, Down:
		crossAxisLength = l.totalSize.X
	case Left, Right:
		crossAxisLength = l.totalSize.Y
	}

	resultSize := f(l.itemPos, crossAxisLength)

	switch l.dir {
	case Up:
		l.itemPos = l.itemPos.Sub(image.Pt(0, resultSize.Y+l.ui.Style.Spacing))
	case Down:
		l.itemPos = l.itemPos.Add(image.Pt(0, resultSize.Y+l.ui.Style.Spacing))
	case Left:
		l.itemPos = l.itemPos.Sub(image.Pt(resultSize.X+l.ui.Style.Spacing, 0))
	case Right:
		l.itemPos = l.itemPos.Add(image.Pt(resultSize.X+l.ui.Style.Spacing, 0))
	}
}

func (l *ListLayoutWithExcess) Excess(f func(pos, size image.Point)) {
	deltaPos := l.itemPos.Sub(l.startPos)

	var remainingSize image.Point

	switch l.dir {
	case Up, Down:
		remainingSize = image.Pt(l.totalSize.X, l.totalSize.Y-imath.AbsInt(deltaPos.Y))
	case Left, Right:
		remainingSize = image.Pt(l.totalSize.X-imath.AbsInt(deltaPos.X), l.totalSize.Y)
	}

	f(l.itemPos, remainingSize)
}
