package imgui

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"runtime"
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
				Hot:    hotPrevious,
				Active: activePrevious,

				MousePosPrevious:    image.Pt(mouseXPrevious, mouseYPrevious),
				MousePos:            image.Pt(mouseX, mouseY),
				IsMouseDownPrevious: mouseLeftStatePrevious == glfw.Press,
				IsMouseDown:         mouseLeftState == glfw.Press,

				Img: currentTex.Img,
			}
			defer func() {
				hotPrevious = ctx.Hot
				activePrevious = ctx.Active
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
	if (ctx.DoButton("b1", "Hello, world!", image.Rect(100, 100, 300, 180), color.RGBA{255, 0, 0, 255})) {
		log.Print("Button 1 clicked!")
	}

	if (ctx.DoButton("b2", "Hello, world!", image.Rect(200, 200, 400, 280), color.RGBA{0, 0, 255, 255})) {
		log.Print("Button 2 clicked!")
	}
}
