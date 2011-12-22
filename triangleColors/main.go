package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gl"
	"github.com/jteeuwen/glfw"
)

const (
	ScreenHeight = 640
	ScreenWidth  = 640
	WindowTitle  = "OpenGL 3 example - TriangleColors"
)

const (
	VertexShaderType = iota
	FragmentShaderType
)

var triangle = []float32{
	0.0, 0.8,
	-0.8, -0.8,
	0.8, -0.8}

var colorTriangle = []float32{
	1, 1, 0,
	0, 0, 1,
	1, 0, 0}

var buffer gl.Buffer
var colorBuffer gl.Buffer

var vShader gl.Shader
var fShader gl.Shader
var program gl.Program

var attribLoc gl.AttribLocation
var attribLocColor gl.AttribLocation

func loadShader(name string) (string, os.Error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func createShader(name string, shaderType int) (gl.Shader, os.Error) {
	data, err := loadShader(name)
	if err != nil {
		return 0, err
	}
	var shader gl.Shader
	switch shaderType {
	case VertexShaderType:
		shader = gl.CreateShader(gl.VERTEX_SHADER)
	case FragmentShaderType:
		shader = gl.CreateShader(gl.FRAGMENT_SHADER)
	default:
		return 0, os.NewError("Unknown ShaderType.")
	}
	shader.Source(data)
	shader.Compile()
	infoLog := shader.GetInfoLog()
	if len(infoLog) != 0 {
		shader.Delete()
		return 0, err
	}
	errNum := gl.GetError()
	if errNum != 0 {
		return 0, os.NewError(fmt.Sprintf("Error code: %d", errNum))
	}
	return shader, nil
}

func initResources() {
	var err os.Error
	buffer = gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(triangle)*4, triangle, gl.STATIC_DRAW)
	gl.Buffer(0).Bind(gl.ARRAY_BUFFER)

	// Color
	colorBuffer = gl.GenBuffer()
	colorBuffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(colorTriangle)*4, colorTriangle, gl.STATIC_DRAW)
	gl.Buffer(0).Bind(gl.ARRAY_BUFFER)

	vShader, err = createShader("triangle.v.glsl", VertexShaderType)
	if err != nil {
		fmt.Printf("Shader: %s\n", err)
	}
	fShader, err = createShader("triangle.f.glsl", FragmentShaderType)
	if err != nil {
		fmt.Printf("Shader: %s\n", err)
	}

	program = gl.CreateProgram()
	program.AttachShader(vShader)
	program.AttachShader(fShader)
	program.Link()

	attribLoc = program.GetAttribLocation("coord2d")
	attribLocColor = program.GetAttribLocation("v_color")
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

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	fmt.Printf("Error-code: %d Init-Status: %d\n", gl.GetError(), initStatus)

	initResources()

	fmt.Printf("Error-code: %d Init-Status: %d\n", gl.GetError(), initStatus)

	for i := 0; i < 1000; i++ {
		display()
	}

	// Free resources
	free()
}

func free() {
	program.Delete()
	buffer.Delete()
}

func display() {
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	program.Use()
	attribLoc.EnableArray()
	attribLocColor.EnableArray()

	buffer.Bind(gl.ARRAY_BUFFER)
	attribLoc.AttribPointerOffset(2, gl.FLOAT, false, 0, 0)
	colorBuffer.Bind(gl.ARRAY_BUFFER)
	attribLocColor.AttribPointerOffset(3, gl.FLOAT, false, 0, 0)

	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.Buffer(0).Bind(gl.ARRAY_BUFFER)
	attribLoc.DisableArray()
	attribLocColor.DisableArray()

	glfw.SwapBuffers()
}
