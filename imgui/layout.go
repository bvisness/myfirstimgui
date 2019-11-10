package imgui

import (
	"fmt"

	"github.com/bvisness/myfirstimgui/rectutil"
)

func dirToPlacerMode(d Direction) rectutil.PlacementMode {
	switch d {
	case Right, Down:
		return rectutil.UpperLeft
	case Up:
		return rectutil.LowerLeft
	case Left:
		return rectutil.UpperRight
	}

	panic(fmt.Errorf("invalid direction: %v", d))
}
