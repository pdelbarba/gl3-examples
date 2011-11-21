package main

import (
	"fmt"
	"os"

	"gl"
	"github.com/jteeuwen/glfw"
)

const (
	ScreenHeight = 640
	ScreenWidth  = 640
	WindowTitle  = "OpenGL 3 example"
)

var triangle = []float32{
	0.0, 0.8,
	-0.8, -0.8,
	0.8, -0.8,
}

var vertexShader = "#version 120\nattribute vec2 coord2d;\nvoid main(void) {\n gl_Position = vec4(coord2d, 0.0, 1.0);\n}"
var fragmentShader = "#version 120\nvoid main(void) {\n gl_FragColor[0] = gl_FragCoord.x/640.0;\n gl_FragColor[1] = gl_FragCoord.y/480.0;\n gl_FragColor[2] = 0.5;\n}"

var vs gl.Shader
var fs gl.Shader
var program gl.Program
var attributeCoord2d gl.AttribLocation

func initResources() {
	var infoLog string

	vs = gl.CreateShader(gl.VERTEX_SHADER)
	vs.Source(vertexShader)
	vs.Compile()
	infoLog = vs.GetInfoLog()
	if len(infoLog) != 0 {
		fmt.Println(infoLog)
	}

	fs = gl.CreateShader(gl.FRAGMENT_SHADER)
	fs.Source(fragmentShader)
	fs.Compile()
	infoLog = fs.GetInfoLog()
	if len(infoLog) != 0 {
		fmt.Println(infoLog)
	}

	program = gl.CreateProgram()
	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()
	infoLog = program.GetInfoLog()
	if len(infoLog) != 0 {
		fmt.Println(infoLog)
	}

	attribute_name := "coord2d"
	attributeCoord2d = program.GetAttribLocation(attribute_name)
	if attributeCoord2d == -1 {
		fmt.Printf("Could not bind attribute %s\n", attribute_name)
	}
}

func main() {
	var err os.Error
	err = glfw.Init()
	if err != nil {
		fmt.Printf("GLFW: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.WindowNoResize, 1)
	glfw.OpenWindowHint(glfw.OpenGLDebugContext, 1)

	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 3)
	glfw.OpenWindowHint(glfw.OpenGLProfile, 1)
	glfw.OpenWindowHint(glfw.OpenGLForwardCompat, 1)

	err = glfw.OpenWindow(ScreenWidth, ScreenHeight, 0, 0, 0, 0, 0, 0, glfw.Windowed)
	if err != nil {
		fmt.Printf("GLFW: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1)
	glfw.SetWindowTitle(WindowTitle)

	major, minor, rev := glfw.GLVersion()
	fmt.Printf("GL-Version: %d, %d, %d\n", major, minor, rev)

	initStatus := gl.Init() // Init glew

	fmt.Printf("Error-code: %d Init-Status: %d\n", gl.GetError(), initStatus)

	initResources()

	for i := 0; i < 1000; i++ {
		display()
	}

	// Free resources
	program.Delete()
}

func display() {
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	program.Use()
	attributeCoord2d.EnableArray()

	attributeCoord2d.AttribPointer(2, false, 0, triangle)

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	attributeCoord2d.DisableArray()

	glfw.SwapBuffers()
}
