package imgui

import (
	"image"
	"image/color"
	"log"

	"github.com/bvisness/myfirstimgui/util"
)

type UIID struct {
	Name string
}

func (id UIID) String() string {
	return id.Name
}

type SizedResult interface {
	DrawnSize() image.Point
}

type UIBase struct {
	Hot    *UIID
	Active *UIID
}

type UIMouse struct {
	PosPrevious image.Point
	Pos         image.Point

	IsMouseDownPrevious bool
	IsMouseDown         bool
}

type UIPosition struct {
	image.Point
}

type UISize struct {
	image.Point
}

type UIContext struct {
	// this will never be nil
	Base *UIBase

	Mouse UIMouse

	Pos  *UIPosition
	Size *UISize

	Img UIImage
}

func (ctx *UIContext) IsHot(obj UIID) bool {
	if ctx.Base.Hot == nil {
		return false
	}

	return *ctx.Base.Hot == obj
}

func (ctx *UIContext) IsActive(obj UIID) bool {
	if ctx.Base.Active == nil {
		return false
	}

	return *ctx.Base.Active == obj
}

func (ctx *UIContext) SetHot(obj UIID) {
	if ctx.Base.Active == nil {
		ctx.Base.Hot = &obj
	}
}

func (ctx *UIContext) SetActive(obj UIID) {
	ctx.Base.Active = &obj
}

func (ctx *UIContext) SetNoneActive() {
	ctx.Base.Active = nil
}

func (ctx *UIContext) IsMouseDownThisFrame() bool {
	return !ctx.Mouse.IsMouseDownPrevious && ctx.Mouse.IsMouseDown
}

func (ctx *UIContext) IsMouseUpThisFrame() bool {
	return ctx.Mouse.IsMouseDownPrevious && !ctx.Mouse.IsMouseDown
}

func (ctx *UIContext) WithPosition(pos image.Point) *UIContext {
	newCtx := *ctx
	newCtx.Pos = &UIPosition{pos}

	return &newCtx
}

func (ctx *UIContext) WithSize(size image.Point) *UIContext {
	newCtx := *ctx
	newCtx.Size = &UISize{size}

	return &newCtx
}

type ButtonResult struct {
	Clicked bool
	Size    image.Point
}

func (br ButtonResult) DrawnSize() image.Point {
	return br.Size
}

func (ctx *UIContext) DoButton(id, text string, c color.RGBA) ButtonResult {
	if ctx.Pos == nil || ctx.Size == nil {
		log.Printf("ERROR: Button '%v' didn't have enough info to be displayed", id)
		return ButtonResult{}
	}

	me := UIID{
		Name: id,
	}
	result := false

	r := image.Rect(
		ctx.Pos.X,
		ctx.Pos.Y,
		ctx.Pos.X+ctx.Size.X,
		ctx.Pos.Y+ctx.Size.Y,
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

func (ctx *UIContext) NewListLayouter(spacing int, horizontal bool) *ListLayouter {
	if ctx.Pos == nil {
		log.Printf("ERROR: List layouter needs a starting position")
		return nil
	}

	return &ListLayouter{
		Size:       image.Pt(0, 0),
		ctx:        ctx,
		itemPos:    (*ctx.Pos).Point,
		spacing:    spacing,
		horizontal: horizontal,
	}
}

func (l *ListLayouter) Item(f func(ctx *UIContext) image.Point) {
	resultSize := f(l.ctx.WithPosition(l.itemPos))

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
