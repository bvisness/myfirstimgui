package imgui

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"runtime"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const width = 500
const height = 500

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
				Base: &UIBase{
					Hot:    hotPrevious,
					Active: activePrevious,
				},

				Mouse: UIMouse{
					PosPrevious:         image.Pt(mouseXPrevious, mouseYPrevious),
					Pos:                 image.Pt(mouseX, mouseY),
					IsMouseDownPrevious: mouseLeftStatePrevious == glfw.Press,
					IsMouseDown:         mouseLeftState == glfw.Press,
				},

				Img: currentTex.Img,
			}
			defer func() {
				hotPrevious = ctx.Base.Hot
				activePrevious = ctx.Base.Active
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

func doUI(ctx *UIContext) {
	doListItem := ctx.NewListLayouter(image.Pt(100, 100), 20)

	doListItem(func(ctx *UIContext) SizedResult {
		result := ctx.WithSize(image.Pt(200, 40)).DoButton("b1", "Hello, world!", color.RGBA{255, 0, 0, 255})
		if result.Clicked {
			log.Print("Button 1 clicked!")
		}

		return result
	})
	doListItem(func(ctx *UIContext) SizedResult {
		height := 80 + int(math.Sin(float64(time.Now().UnixNano())/float64(time.Second))*30)

		result := ctx.WithSize(image.Pt(200, height)).DoButton("b2", "Hello, world!", color.RGBA{0, 0, 255, 255})
		if result.Clicked {
			log.Print("Button 2 clicked!")
		}

		return result
	})
	doListItem(func(ctx *UIContext) SizedResult {
		result := ctx.WithSize(image.Pt(200, 80)).DoButton("b3", "Hello, world!", color.RGBA{0, 255, 0, 255})
		if result.Clicked {
			log.Print("Button 3 clicked!")
		}

		return result
	})
}
