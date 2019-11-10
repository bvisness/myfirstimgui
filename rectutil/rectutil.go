package rectutil

import (
	"image"

	"github.com/bvisness/myfirstimgui/imath"
)

func MinPoint(p1, p2 image.Point) image.Point {
	return image.Pt(imath.MinInt(p1.X, p2.X), imath.MinInt(p1.Y, p2.Y))
}

func MaxPoint(p1, p2 image.Point) image.Point {
	return image.Pt(imath.MaxInt(p1.X, p2.X), imath.MaxInt(p1.Y, p2.Y))
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
