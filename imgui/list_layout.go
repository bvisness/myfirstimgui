package imgui

import "image"

type ListLayout struct {
	Size image.Point

	ui         *UIContext
	itemPos    image.Point
	horizontal bool
}

func (ui *UIContext) ListLayout(startPos image.Point, spacing int, horizontal bool) ListLayout {
	return ListLayout{
		Size:       image.Pt(0, 0),
		ui:         ui,
		itemPos:    startPos,
		horizontal: horizontal,
	}
}

func (l *ListLayout) Item(f func(pos image.Point) image.Point) {
	resultSize := f(l.itemPos)

	if l.horizontal {
		l.itemPos = l.itemPos.Add(image.Pt(resultSize.X+l.ui.Style.Spacing, 0))
	} else {
		l.itemPos = l.itemPos.Add(image.Pt(0, resultSize.Y+l.ui.Style.Spacing))
	}

	if resultSize.X > l.Size.X {
		l.Size.X = resultSize.X
	}
	if resultSize.Y > l.Size.Y {
		l.Size.Y = resultSize.Y
	}
}
