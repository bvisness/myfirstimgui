package imgui

import (
	"image"

	"github.com/bvisness/myfirstimgui/imath"
	"github.com/bvisness/myfirstimgui/rectutil"
)

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

func (l *ListLayoutWithExcess) Item(f func(rectPlacer rectutil.Placer, crossAxisLength int) image.Point) {
	var crossAxisLength int
	switch l.dir {
	case Up, Down:
		crossAxisLength = l.totalSize.X
	case Left, Right:
		crossAxisLength = l.totalSize.Y
	}

	resultSize := f(rectutil.NewPlacer(dirToPlacerMode(l.dir), l.itemPos), crossAxisLength)

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

func (l *ListLayoutWithExcess) Excess(f func(r image.Rectangle)) {
	deltaPos := l.itemPos.Sub(l.startPos)

	var remainingSize image.Point

	switch l.dir {
	case Up, Down:
		remainingSize = image.Pt(l.totalSize.X, l.totalSize.Y-imath.Abs(deltaPos.Y))
	case Left, Right:
		remainingSize = image.Pt(l.totalSize.X-imath.Abs(deltaPos.X), l.totalSize.Y)
	}

	f(rectutil.PlaceSize(remainingSize, l.itemPos, dirToPlacerMode(l.dir)))
}
