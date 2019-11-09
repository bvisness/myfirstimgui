package imgui

import (
	"image"
	"image/color"

	"github.com/bvisness/myfirstimgui/util"
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

type UIContext struct {
	Hot    *UIID
	Active *UIID

	Mouse UIMouse

	Img UIImage
}

func (ctx *UIContext) IsHot(obj UIID) bool {
	if ctx.Hot == nil {
		return false
	}

	return *ctx.Hot == obj
}

func (ctx *UIContext) IsActive(obj UIID) bool {
	if ctx.Active == nil {
		return false
	}

	return *ctx.Active == obj
}

func (ctx *UIContext) SetHot(obj UIID) {
	if ctx.Active == nil {
		ctx.Hot = &obj
	}
}

func (ctx *UIContext) SetActive(obj UIID) {
	ctx.Active = &obj
}

func (ctx *UIContext) SetNoneActive() {
	ctx.Active = nil
}

func (ctx *UIContext) IsMouseDownThisFrame() bool {
	return !ctx.Mouse.IsMouseDownPrevious && ctx.Mouse.IsMouseDown
}

func (ctx *UIContext) IsMouseUpThisFrame() bool {
	return ctx.Mouse.IsMouseDownPrevious && !ctx.Mouse.IsMouseDown
}

type ButtonResult struct {
	Clicked bool
	Size    image.Point
}

func (ctx *UIContext) Button(id, text string, pos image.Point, size image.Point, c color.RGBA) ButtonResult {
	me := UIID{
		Name: id,
	}
	result := false

	r := image.Rect(
		pos.X,
		pos.Y,
		pos.X+size.X,
		pos.Y+size.Y,
	)

	if ctx.IsActive(me) {
		if ctx.IsMouseUpThisFrame() {
			if ctx.IsHot(me) {
				result = true
			}
			ctx.SetNoneActive()
		}
	} else if ctx.IsHot(me) {
		if ctx.IsMouseDownThisFrame() {
			ctx.SetActive(me)
		}
	}

	if util.PointInRect(ctx.Mouse.Pos, r) {
		ctx.SetHot(me)
	}

	if ctx.IsActive(me) {
		c = AlphaOver(c, color.RGBA{255, 255, 255, 50})
	} else if ctx.IsHot(me) {
		c = AlphaOver(c, color.RGBA{255, 255, 255, 100})
	}
	ctx.Img.DrawRect(r, c)

	return ButtonResult{
		Clicked: result,
		Size:    r.Size(),
	}
}

type ListLayouter struct {
	Size image.Point

	ctx        *UIContext
	itemPos    image.Point
	spacing    int
	horizontal bool
}

func (ctx *UIContext) NewListLayouter(startPos image.Point, spacing int, horizontal bool) *ListLayouter {
	return &ListLayouter{
		Size:       image.Pt(0, 0),
		ctx:        ctx,
		itemPos:    startPos,
		spacing:    spacing,
		horizontal: horizontal,
	}
}

func (l *ListLayouter) Item(f func(pos image.Point) image.Point) {
	resultSize := f(l.itemPos)

	if l.horizontal {
		l.itemPos = l.itemPos.Add(image.Pt(resultSize.X+l.spacing, 0))
	} else {
		l.itemPos = l.itemPos.Add(image.Pt(0, resultSize.Y+l.spacing))
	}

	if resultSize.X > l.Size.X {
		l.Size.X = resultSize.X
	}
	if resultSize.Y > l.Size.Y {
		l.Size.Y = resultSize.Y
	}
}
