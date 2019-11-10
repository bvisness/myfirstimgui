package util

import (
	"image"
)

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func MinPoint(p1, p2 image.Point) image.Point {
	return image.Pt(MinInt(p1.X, p2.X), MinInt(p1.Y, p2.Y))
}

func MaxPoint(p1, p2 image.Point) image.Point {
	return image.Pt(MaxInt(p1.X, p2.X), MaxInt(p1.Y, p2.Y))
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
