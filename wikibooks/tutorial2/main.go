package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
)

const (
	ScreenHeight = 480
	ScreenWidth  = 640
	WindowTitle  = "OpenGL 3 tutorial 2 - Managing shaders"
)

const (
	VertexShaderType = iota
	FragmentShaderType
)

var triangleVertices = []float32{
	0.0, 0.8,
	-0.8, -0.8,
	0.8, -0.8,
}
var vboTriangle gl.Uint

var vs gl.Uint
var fs gl.Uint
var program gl.Uint
var attributeCoord2d gl.Uint

func fileRead(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func createShader(name string, shaderType int) (gl.Uint, error) {
	data, err := fileRead(name)
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, errors.New("No shader code.")
	}
	var shader gl.Uint
	switch shaderType {
	case VertexShaderType:
		shader = gl.CreateShader(gl.VERTEX_SHADER)
	case FragmentShaderType:
		shader = gl.CreateShader(gl.FRAGMENT_SHADER)
	default:
		return 0, errors.New("Unknown ShaderType.")
	}
	src := gl.GLStringArray(string(data))
	defer gl.GLStringArrayFree(src)
	gl.ShaderSource(shader, gl.Sizei(1), &src[0], nil)
	gl.CompileShader(shader)

	// Similar to print_log in the C code example
	var length gl.Int
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
	if length > 1 {
		glString := gl.GLStringAlloc(gl.Sizei(length))
		defer gl.GLStringFree(glString)
		gl.GetShaderInfoLog(shader, gl.Sizei(length), nil, glString)
		return 0, errors.New(fmt.Sprintf("Shader log: %s", gl.GoString(glString)))
	}
	return shader, nil
}

func initResources() {
	var err error
	// Load shaders
	vs, err = createShader("triangle.v.glsl", VertexShaderType)
	if err != nil {
		fmt.Printf("Shader: %s\n", err)
		return
	}
	fs, err = createShader("triangle.f.glsl", FragmentShaderType)
	if err != nil {
		fmt.Printf("Shader: %s\n", err)
		return
	}
	var compileOk gl.Int
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

	// Generate a buffer for the VertexBufferObject
	gl.GenBuffers(1, &vboTriangle)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboTriangle)
	// Submit the vertices of the triangle to the graphic card
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(len(triangleVertices)*4), gl.Pointer(&triangleVertices[0]), gl.STATIC_DRAW)
	// Unbind the active buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func main() {
	var err error
	err = glfw.Init()
	if err != nil {
		fmt.Printf("GLFW: %s\n", err)
		return
	}
	defer glfw.Terminate()

	// You could probably change the required versions down
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 3)
	glfw.OpenWindowHint(glfw.OpenGLProfile, 1)

	// Open Window with 8 bit Alpha
	err = glfw.OpenWindow(ScreenWidth, ScreenHeight, 0, 0, 0, 8, 0, 0, glfw.Windowed)
	if err != nil {
		fmt.Printf("GLFW: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetWindowTitle(WindowTitle)

	major, minor, rev := glfw.GLVersion()
	if major < 3 {
		fmt.Printf("Error your graphic card does not support OpenGL 3.3\n Your GL-Version is: %d, %d, %d\n", major, minor, rev)
		fmt.Println("You can try to lower the settings in glfw.OpenWindowHint(glfw.OpenGLVersionMajor/Minor.")
	}

	// Init extension loading
	err = gl.Init()
	if err != nil {
		fmt.Printf("Init OpenGL extension loading failed with %s.\n", err)
	}

	// Enable transparency in OpenGL
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	initResources()

	for {
		display()
	}

	// Free resources
	free()
}

func free() {
	gl.DeleteProgram(program)
	gl.DeleteBuffers(1, &vboTriangle)
}

func display() {
	// Clear the background as white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Use the GLSL program
	gl.UseProgram(program)

	gl.BindBuffer(gl.ARRAY_BUFFER, vboTriangle)
	gl.EnableVertexAttribArray(attributeCoord2d)

	// Describe our vertices array to OpenGL (it can't guess its format automatically)
	gl.VertexAttribPointer(attributeCoord2d, 2, gl.FLOAT, gl.FALSE, 0, gl.Pointer(nil))

	// Push each element in buffer_vertices to the vertex shader
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.DisableVertexAttribArray(attributeCoord2d)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // Unbind

	// Display the result
	glfw.SwapBuffers()
}
