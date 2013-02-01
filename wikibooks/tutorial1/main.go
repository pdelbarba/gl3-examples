package main

import (
	"fmt"

	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
)

const (
	ScreenHeight = 480
	ScreenWidth  = 640
	WindowTitle  = "OpenGL 3 tutorial 1 - My First Triangle"
)

var (
	vsSource = "#version 120\n" +
		"attribute vec2 coord2d;\n" +
		"void main(void) {\n" +
		"  gl_Position = vec4(coord2d, 0.0, 1.0);\n" +
		"}"
	fsSource = "#version 120\n" +
		"void main(void) {\n" +
		"  gl_FragColor[0] = 0.0;\n" +
		"  gl_FragColor[1] = 0.0;\n" +
		"  gl_FragColor[2] = 1.0;\n" +
		"}"
)

var triangleVertices = []float32{
	0.0, 0.8,
	-0.8, -0.8,
	0.8, -0.8,
}

var vs gl.Uint
var fs gl.Uint
var program gl.Uint
var attributeCoord2d gl.Uint

func main() {
	fmt.Println("OpenGL Programming/Modern OpenGL Introduction")
	fmt.Println("Tutorial taken from http://en.wikibooks.org/wiki/OpenGL_Programming/Modern_OpenGL_Introduction")

	var err error
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

	// Init extension loading
	err = gl.Init()
	if err != nil {
		fmt.Printf("Init OpenGL extension loading failed with %s.\n", err)
	}

	initResources()

	// Render loop
	for glfw.WindowParam(glfw.Opened) == 1 {
		display()
	}

	// Free resources
	free()
}

func initResources() {
	var compileOk gl.Int
	// Vertex Shader
	vs = gl.CreateShader(gl.VERTEX_SHADER)
	vsSrc := gl.GLStringArray(vsSource)
	defer gl.GLStringArrayFree(vsSrc)
	gl.ShaderSource(vs, gl.Sizei(len(vsSrc)), &vsSrc[0], nil)
	gl.CompileShader(vs)
	gl.GetShaderiv(vs, gl.COMPILE_STATUS, &compileOk)
	if compileOk == 0 {
		errNum := gl.GetError()
		fmt.Printf("Error in vertex shader: %d\n", errNum)
	}

	// Fragment Shader
	fs = gl.CreateShader(gl.FRAGMENT_SHADER)
	fsSrc := gl.GLStringArray(fsSource)
	defer gl.GLStringArrayFree(fsSrc)
	gl.ShaderSource(fs, gl.Sizei(1), &fsSrc[0], nil)
	gl.CompileShader(fs)
	gl.GetShaderiv(fs, gl.COMPILE_STATUS, &compileOk)
	if compileOk == 0 {
		errNum := gl.GetError()
		fmt.Printf("Error in fragment shader: %d\n", errNum)
	}

	// GLSL program
	program = gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)
	gl.GetProgramiv(program, gl.LINK_STATUS, &compileOk)
	if compileOk == 0 {
		fmt.Printf("Error in program.\n")

	}

	// Get the attribute location from the GLSL program (here from the vertex shader)
	attributeName := gl.GLString("coord2d")
	defer gl.GLStringFree(attributeName)
	attributeTemp := gl.GetAttribLocation(program, attributeName)
	if attributeTemp == -1 {
		fmt.Printf("Could not bind attribute %s\n", gl.GoString(attributeName))
	}
	attributeCoord2d = gl.Uint(attributeTemp)
}

func free() {
	// Free OpenGL buffers
	gl.DeleteProgram(program)
}

func display() {
	// Clear the background as white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Use the GLSL program
	gl.UseProgram(program)

	gl.EnableVertexAttribArray(attributeCoord2d)

	// Describe our vertices array to OpenGL (it can't guess its format automatically)
	gl.VertexAttribPointer(attributeCoord2d, 2, gl.FLOAT, gl.FALSE, 0, gl.Pointer(&triangleVertices[0]))

	// Push each element in buffer_vertices to the vertex shader
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.DisableVertexAttribArray(attributeCoord2d)

	// Display the result
	glfw.SwapBuffers()
}
