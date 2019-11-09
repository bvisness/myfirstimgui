package imgui

import (
	"image"
	"image/color"

	"github.com/bvisness/myfirstimgui/util"
)

type UIID struct {
	Name string
}

type UIContext struct {
	Hot    *UIID
	Active *UIID

	MousePosPrevious    image.Point
	MousePos            image.Point
	IsMouseDownPrevious bool
	IsMouseDown         bool

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
	return !ctx.IsMouseDownPrevious && ctx.IsMouseDown
}

func (ctx *UIContext) IsMouseUpThisFrame() bool {
	return ctx.IsMouseDownPrevious && !ctx.IsMouseDown
}

func (ctx *UIContext) DoButton(id, text string, r image.Rectangle, c color.RGBA) bool {
	me := UIID{
		Name: id,
	}
	result := false

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

	if util.PointInRect(ctx.MousePos, r) {
		ctx.SetHot(me)
	}

	if ctx.IsActive(me) {
		c = AlphaOver(c, color.RGBA{255, 255, 255, 50})
	} else if ctx.IsHot(me) {
		c = AlphaOver(c, color.RGBA{255, 255, 255, 100})
	}
	ctx.Img.DrawRect(r, c)

	return result
}
