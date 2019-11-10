package rectutil

import (
	"fmt"
	"image"

	"github.com/bvisness/myfirstimgui/imath"
)

type PlacementMode int

const (
	UpperLeft PlacementMode = iota + 1
	LowerLeft
	UpperRight
)

func MinPoint(p1, p2 image.Point) image.Point {
	return image.Pt(imath.Min(p1.X, p2.X), imath.Min(p1.Y, p2.Y))
}

func MaxPoint(p1, p2 image.Point) image.Point {
	return image.Pt(imath.Max(p1.X, p2.X), imath.Max(p1.Y, p2.Y))
}

func SizeRect(pos, size image.Point) image.Rectangle {
	return image.Rect(
		pos.X,
		pos.Y,
		pos.X+size.X,
		pos.Y+size.Y,
	)
}

func PointInRect(p image.Point, r image.Rectangle) bool {
	return p.X >= r.Min.X && p.X <= r.Max.X && p.Y >= r.Min.Y && p.Y <= r.Max.Y
}

func GetUL(r image.Rectangle) image.Point {
	return r.Min
}

func GetLL(r image.Rectangle) image.Point {
	return image.Pt(r.Min.X, r.Max.Y)
}

func GetUR(r image.Rectangle) image.Point {
	return image.Pt(r.Max.X, r.Min.Y)
}

func GetLR(r image.Rectangle) image.Point {
	return r.Max
}

func PlaceSizeUL(s, p image.Point) image.Rectangle {
	return SizeRect(p, s)
}

func PlaceSizeLL(s, p image.Point) image.Rectangle {
	return image.Rect(
		p.X,
		p.Y-s.Y,
		p.X+s.X,
		p.Y,
	)
}

func PlaceSizeUR(s, p image.Point) image.Rectangle {
	return image.Rect(
		p.X-s.X,
		p.Y,
		p.X,
		p.Y+s.Y,
	)
}

func PlaceSize(s, p image.Point, mode PlacementMode) image.Rectangle {
	switch mode {
	case UpperLeft:
		return PlaceSizeUL(s, p)
	case LowerLeft:
		return PlaceSizeLL(s, p)
	case UpperRight:
		return PlaceSizeUR(s, p)
	}

	panic(fmt.Errorf("invalid placement mode: %v", mode))
}

func PlaceRectUL(r image.Rectangle, p image.Point) image.Rectangle {
	return PlaceSizeUL(r.Size(), p)
}

func PlaceRectLL(r image.Rectangle, p image.Point) image.Rectangle {
	return PlaceSizeLL(r.Size(), p)
}

func PlaceRectUR(r image.Rectangle, p image.Point) image.Rectangle {
	return PlaceSizeUR(r.Size(), p)
}

func PlaceRect(r image.Rectangle, p image.Point, mode PlacementMode) image.Rectangle {
	switch mode {
	case UpperLeft:
		return PlaceRectUL(r, p)
	case LowerLeft:
		return PlaceRectLL(r, p)
	case UpperRight:
		return PlaceRectUR(r, p)
	}

	panic(fmt.Errorf("invalid placement mode: %v", mode))
}

// Places a rect in the right place, according to the current layout mode. Preserves the rect's size.
type Placer struct {
	Mode PlacementMode
	Pos  image.Point
}

func NewPlacer(m PlacementMode, p image.Point) Placer {
	return Placer{
		Mode: m,
		Pos:  p,
	}
}

func (p Placer) PlaceSize(s image.Point) image.Rectangle {
	return PlaceSize(s, p.Pos, p.Mode)
}

func (p Placer) PlaceRect(r image.Rectangle) image.Rectangle {
	return PlaceRect(r, p.Pos, p.Mode)
}
