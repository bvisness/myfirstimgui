package imgui

import "image"

type ListLayout struct {
	Size image.Point

	ui         *UIContext
	spacing    int
	itemPos    image.Point
	horizontal bool
}

func (ui *UIContext) ListLayout(startPos image.Point, horizontal bool, styleOverride *UIStyle) ListLayout {
	spacing := ui.Style.Spacing
	if styleOverride != nil {
		spacing = styleOverride.Spacing
	}

	return ListLayout{
		Size:       image.Pt(0, 0),
		ui:         ui,
		spacing:    spacing,
		itemPos:    startPos,
		horizontal: horizontal,
	}
}

func (l *ListLayout) Item(f func(pos image.Point) image.Point) {
	resultSize := f(l.itemPos)

	if l.horizontal {
		l.itemPos = l.itemPos.Add(image.Pt(resultSize.X+l.spacing, 0))
	} else {
		l.itemPos = l.itemPos.Add(image.Pt(0, resultSize.Y+l.spacing))
	}

	if resultSize.X > l.Size.X {
		l.Size.X = resultSize.X
	}
	if resultSize.Y > l.Size.Y {
		l.Size.Y = resultSize.Y
	}
}
