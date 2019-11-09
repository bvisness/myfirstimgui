package imgui

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlphaOver(t *testing.T) {
	assert.Equal(t, color.RGBA{0, 0, 0, 255}, AlphaOver(color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 0}))
	assert.Equal(t, color.RGBA{0, 0, 0, 255}, AlphaOver(color.RGBA{255, 255, 255, 0}, color.RGBA{0, 0, 0, 255}))
	assert.Equal(t, color.RGBA{100, 50, 25, 254}, AlphaOver(color.RGBA{0, 0, 0, 255}, color.RGBA{200, 100, 50, 128}))
}
