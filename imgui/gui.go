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
		log.Printf("%v is now hot", obj)
		ctx.Base.Hot = &obj
	}
}

func (ctx *UIContext) SetActive(obj UIID) {
	log.Printf("%v is now active", obj)
	ctx.Base.Active = &obj
}

func (ctx *UIContext) SetNoneActive() {
	log.Printf("nothing is active")
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

func (ctx *UIContext) NewListLayouter(startPos image.Point, spacing int) func(func(ctx *UIContext) SizedResult) {
	pos := startPos

	return func(f func(ctx *UIContext) SizedResult) {
		result := f(ctx.WithPosition(pos))
		pos = pos.Add(image.Pt(0, result.DrawnSize().Y+spacing))
	}
}
