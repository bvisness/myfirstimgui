package imgui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/bvisness/myfirstimgui/imath"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type HalfRectMode int

const (
	UpperLeft HalfRectMode = iota
	LowerLeft
	UpperRight
	LowerRight
)

type UIImage struct {
	*image.RGBA
}

type UITexture struct {
	GLTexture uint32
	Img       UIImage
}

func NewUITexture(width, height int, c color.RGBA) *UITexture {
	img := UIImage{image.NewRGBA(image.Rect(0, 0, width, height))}
	if img.Stride != img.Rect.Size().X*4 {
		panic(fmt.Errorf("unsupported stride"))
	}
	img.Fill(c)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	bufferTextureData(texture, img.RGBA)

	return &UITexture{
		GLTexture: texture,
		Img:       img,
	}
}

func (t *UITexture) BufferImage() {
	bufferTextureData(t.GLTexture, t.Img.RGBA)
}

func (i UIImage) Fill(c color.RGBA) {
	for x := 0; x < i.Bounds().Max.X; x++ {
		for y := 0; y < i.Bounds().Max.Y; y++ {
			i.SetRGBA(x, y, c)
		}
	}
}

func (i UIImage) DrawRect(r image.Rectangle, c color.RGBA) {
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			i.SetRGBA(x, y, AlphaOver(i.RGBAAt(x, y), c))
		}
	}
}

func (i UIImage) DrawHalfRect(r image.Rectangle, c color.RGBA, mode HalfRectMode) {
	for x := r.Min.X; x <= r.Max.X; x++ {
		var yThreshold int
		switch mode {
		case UpperLeft, LowerRight:
			yThreshold = imath.LerpInt(r.Max.Y, r.Min.Y, r.Min.X, r.Max.X, x)
		default:
			yThreshold = imath.LerpInt(r.Min.Y, r.Max.Y, r.Min.X, r.Max.X, x)
		}

		var whenAbove bool
		switch mode {
		case UpperLeft, UpperRight:
			whenAbove = true
		default:
			whenAbove = false
		}

		for y := r.Min.Y; y <= r.Max.Y; y++ {
			if whenAbove && y <= yThreshold || !whenAbove && y >= yThreshold {
				i.SetRGBA(x, y, AlphaOver(i.RGBAAt(x, y), c))
			}
		}
	}
}

func AlphaOver(dst, src color.RGBA) color.RGBA {
	return color.RGBA{
		R: LerpUint8(dst.R, src.R, src.A),
		G: LerpUint8(dst.G, src.G, src.A),
		B: LerpUint8(dst.B, src.B, src.A),
		A: uint8(float32(dst.A)*(1-float32(src.A)/255)) + src.A,
	}
}

func LerpUint8(a, b, t uint8) uint8 {
	tf := float32(t) / 255
	return uint8(float32(a)*(1-tf) + float32(b)*tf) // TODO: sad times in lerpville
}
