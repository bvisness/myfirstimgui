package util

import "image"

func PointInRect(p image.Point, r image.Rectangle) bool {
	return p.X >= r.Min.X && p.X <= r.Max.X && p.Y >= r.Min.Y && p.Y <= r.Max.Y
}
