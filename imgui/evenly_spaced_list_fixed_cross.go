package imgui

import (
	"image"

	"github.com/bvisness/myfirstimgui/rectutil"
)

type EvenlySpacedListFixedCross struct {
	ui        *UIContext
	n         int
	itemIndex int
	dir       Direction
	totalArea image.Rectangle
}

func (ui *UIContext) EvenlySpacedListFixedCross(n int, totalArea image.Rectangle, dir Direction) EvenlySpacedListFixedCross {
	return EvenlySpacedListFixedCross{
		ui:        ui,
		n:         n,
		itemIndex: 0,
		dir:       dir,
		totalArea: totalArea,
	}
}

func (l *EvenlySpacedListFixedCross) Item(f func(rect image.Rectangle)) {
	var mainSize int
	var crossSize int
	switch l.dir {
	case Up, Down:
		mainSize = l.totalArea.Size().Y
		crossSize = l.totalArea.Size().X
	case Left, Right:
		mainSize = l.totalArea.Size().X
		crossSize = l.totalArea.Size().Y
	}

	sizeLessDividers := mainSize - l.ui.Style.Spacing*(l.n-1)
	itemMainSize := sizeLessDividers / l.n
	itemOffset := (itemMainSize + l.ui.Style.Spacing) * l.itemIndex

	var resultRect image.Rectangle
	switch l.dir {
	case Up:
		resultRect = rectutil.PlaceSizeLL(image.Pt(crossSize, itemMainSize), rectutil.GetLL(l.totalArea).Sub(image.Pt(0, itemOffset)))
	case Down:
		resultRect = rectutil.PlaceSizeUL(image.Pt(crossSize, itemMainSize), rectutil.GetUL(l.totalArea).Add(image.Pt(0, itemOffset)))
	case Left:
		resultRect = rectutil.PlaceSizeUR(image.Pt(itemMainSize, crossSize), rectutil.GetUR(l.totalArea).Sub(image.Pt(itemOffset, 0)))
	case Right:
		resultRect = rectutil.PlaceSizeUL(image.Pt(itemMainSize, crossSize), rectutil.GetUL(l.totalArea).Add(image.Pt(itemOffset, 0)))
	}

	f(resultRect)

	l.itemIndex++
}

func (l *EvenlySpacedListFixedCross) Size() image.Point {
	return l.totalArea.Size()
}
