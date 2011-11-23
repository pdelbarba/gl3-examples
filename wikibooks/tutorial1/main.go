package main

import (
	"fmt"
	"os"

	"gl"
	"github.com/jteeuwen/glfw"
)

const (
	ScreenHeight = 480
	ScreenWidth  = 640
	WindowTitle  = "OpenGL 3 tutorial 1 - My First Triangle"
	vsSource     = "#version 120\n" +
		"attribute vec2 coord2d; " +
		"void main(void) { " +
		"  gl_Position = vec4(coord2d, 0.0, 1.0); " +
		"}"
	fsSource = "#version 120\n" +
		"void main(void) {" +
		"  gl_FragColor[0] = 0.0; " +
		"  gl_FragColor[1] = 0.0; " +
		"  gl_FragColor[2] = 1.0; " +
		"}"
)

var triangleVertices = []float32{
	0.0, 0.8,
	-0.8, -0.8,
	0.8, -0.8,
}

var vs gl.Shader
var fs gl.Shader
var program gl.Program
var attributeCoord2d gl.AttribLocation

func main() {
	fmt.Println("OpenGL Programming/Modern OpenGL Introduction")
	fmt.Println("Tutorial taken from http://en.wikibooks.org/wiki/OpenGL_Programming/Modern_OpenGL_Introduction")

	var err os.Error
	err = glfw.Init()
	if err != nil {
		fmt.Printf("GLFW: %s\n", err)
		return
	}
	defer glfw.Terminate()

	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 3)
	glfw.OpenWindowHint(glfw.OpenGLProfile, 1)

	err = glfw.OpenWindow(ScreenWidth, ScreenHeight, 0, 0, 0, 0, 0, 0, glfw.Windowed)
	if err != nil {
		fmt.Printf("GLFW: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetWindowTitle(WindowTitle)

	// Init glew
	initStatus := gl.Init()
	if initStatus != 0 {
		fmt.Printf("Init glew failed with %d.\n", initStatus)
	}

	initResources()

	// Render loop
	for {
		display()
	}

	// Free resources
	free()
}

func initResources() {
	var status int
	// Vertex Shader
	vs = gl.CreateShader(gl.VERTEX_SHADER)
	vs.Source(vsSource)
	vs.Compile()
	status = vs.Get(gl.COMPILE_STATUS)
	if status == 0 {
		fmt.Printf("Error in vertex shader\n")
	}

	// Fragment Shader
	fs = gl.CreateShader(gl.FRAGMENT_SHADER)
	fs.Source(fsSource)
	fs.Compile()
	status = fs.Get(gl.COMPILE_STATUS)
	if status == 0 {
		fmt.Printf("Error in fragment shader\n")
	}

	// GLSL program
	program = gl.CreateProgram()
	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()
	status = program.Get(gl.LINK_STATUS)
	if status == 0 {
		fmt.Printf("glLinkProgram\n")
	}

	// Get the attribute location from the GLSL program (here from the vertex shader)
	attributeName := "coord2d"
	attributeCoord2d = program.GetAttribLocation(attributeName)
	if attributeCoord2d == -1 {
		fmt.Printf("Could not bind attribute %s\n", attributeName)
	}
}

func free() {
	// Free OpenGL buffers
	program.Delete()
}

func display() {
	// Clear the background as white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Use the GLSL program
	program.Use()

	attributeCoord2d.EnableArray()

	// Describe our vertices array to OpenGL (it can't guess its format automatically)
	attributeCoord2d.AttribPointer(2, false, 0, triangleVertices)

	// Push each element in buffer_vertices to the vertex shader
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	attributeCoord2d.DisableArray()

	// Display the result
	glfw.SwapBuffers()
}
