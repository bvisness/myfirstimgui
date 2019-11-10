package imgui

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"runtime"
	"unsafe"

	"github.com/bvisness/myfirstimgui/rectutil"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const width = 1000
const height = 800

var quadVerts = []float32{
	// vert positions double as texture coordinates (with some mapping)
	0, 0,
	0, 1,
	1, 0,
	1, 0,
	0, 1,
	1, 1,
}
var sizeOfVertNum = int(unsafe.Sizeof(quadVerts[0]))

func Main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	initOpenGL()

	v, err := ioutil.ReadFile("v.glsl")
	if err != nil {
		panic(err)
	}

	f, err := ioutil.ReadFile("f.glsl")
	if err != nil {
		panic(err)
	}

	program, err := newProgram(string(v), string(f))
	if err != nil {
		panic(err)
	}

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("color\x00"))

	vao := genQuadVAO()

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, int32(2*sizeOfVertNum), gl.PtrOffset(0))

	tex1 := NewUITexture(width, height, color.RGBA{255, 0, 0, 255})
	tex2 := NewUITexture(width, height, color.RGBA{255, 50, 0, 255})

	currentTex := tex1

	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	var hotPrevious *UIID = nil
	var activePrevious *UIID = nil
	var elementStatePrevious map[UIID]interface{} = make(map[UIID]interface{})

	mouseXPrevious := -1
	mouseYPrevious := -1
	mouseLeftStatePrevious := glfw.Release

	for !window.ShouldClose() {
		func() {
			mouseXF, mouseYF := window.GetCursorPos()
			mouseX := int(mouseXF)
			mouseY := int(mouseYF)
			mouseLeftState := window.GetMouseButton(glfw.MouseButton1)
			defer func() {
				mouseXPrevious = mouseX
				mouseYPrevious = mouseY
				mouseLeftStatePrevious = mouseLeftState
			}()

			ctx := UIContext{
				Hot:    hotPrevious,
				Active: activePrevious,

				Style: UIStyle{
					Spacing: 16,
				},
				Mouse: UIMouse{
					PosPrevious:         image.Pt(mouseXPrevious, mouseYPrevious),
					Pos:                 image.Pt(mouseX, mouseY),
					IsMouseDownPrevious: mouseLeftStatePrevious == glfw.Press,
					IsMouseDown:         mouseLeftState == glfw.Press,
				},

				Img: currentTex.Img,

				ElementState: elementStatePrevious,
			}
			defer func() {
				// TODO: This is untenable, probably.
				hotPrevious = ctx.Hot
				activePrevious = ctx.Active
				elementStatePrevious = ctx.ElementState
			}()

			currentTex.Img.Fill(color.RGBA{255, 255, 255, 0})
			doUI(&ctx)
			currentTex.BufferImage()

			defer func() {
				// Set up for next frame
				if currentTex == tex1 {
					currentTex = tex2
				} else {
					currentTex = tex1
				}
			}()

			// Draw stuff

			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			gl.UseProgram(program)
			gl.BindVertexArray(vao)
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, currentTex.GLTexture)

			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(quadVerts)))

			glfw.PollEvents()
			window.SwapBuffers()
		}()
	}
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "My first IMGUI", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func genQuadVAO() uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadVerts)*sizeOfVertNum, gl.Ptr(quadVerts), gl.STATIC_DRAW)

	return vao
}

func doUI(ui *UIContext) {
	//columnList := ui.ListLayout(image.Pt(20, 20), 20, true)
	//
	//for c := 0; c < 4; c++ {
	//	columnList.Item(func(pos image.Point) image.Point {
	//		rowList := ui.ListLayout(pos, 20, false)
	//
	//		for r := 0; r < 4; r++ {
	//			id := c*4 + r
	//			rowList.Item(func(pos image.Point) image.Point {
	//				t := float64(time.Now().UnixNano())/float64(time.Second) + float64(id)/2
	//				width := 60 + int(math.Cos(t)*30)
	//				height := 60 + int(math.Sin(t)*30)
	//
	//				result := ui.Button(fmt.Sprintf("b%v", id), "Hello, world!", pos, image.Pt(width, height), buttonColor(r, c))
	//				if result.Clicked {
	//					log.Printf("Button %v clicked!", id)
	//				}
	//
	//				return result.Size
	//			})
	//		}
	//
	//		return rowList.Size
	//	})
	//}

	if open, windowContent := ui.Window("test", image.Pt(100, 100), image.Pt(300, 300), false); open {
		list := ui.ListLayoutWithExcess(rectutil.GetLL(windowContent), windowContent.Size(), Up)

		list.Item(func(placer rectutil.Placer, width int) image.Point {
			buttons := ui.EvenlySpacedListFixedCross(2, placer.PlaceSize(image.Pt(width, 40)), Right)

			buttons.Item(func(rect image.Rectangle) {
				if result := ui.Button("bCancel", "Cancel", rect.Min, rect.Size(), color.RGBA{255, 100, 100, 255}); result.Clicked {
					log.Print("Cancel button clicked")
				}
			})
			buttons.Item(func(rect image.Rectangle) {
				if result := ui.Button("bSubmit", "Submit", rect.Min, rect.Size(), color.RGBA{100, 255, 100, 255}); result.Clicked {
					log.Print("Submit button clicked")
				}
			})

			return buttons.Size()
		})
		list.Excess(func(excess image.Rectangle) {
			ui.Button("b2", "asdf", excess.Min, excess.Size(), color.RGBA{100, 100, 255, 255})
		})
	}
}

func buttonColor(r, c int) color.RGBA {
	return color.RGBA{100 + uint8(r*25), 100 + uint8(c*25), 100, 255}
}
