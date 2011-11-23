package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gl"
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
var vboTriangle gl.Buffer

var vs gl.Shader
var fs gl.Shader
var program gl.Program
var attributeCoord2d gl.AttribLocation

func fileRead(name string) (string, os.Error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func createShader(name string, shaderType int) (gl.Shader, os.Error) {
	data, err := fileRead(name)
	if err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, os.NewError("No shader code.")
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

	// Similar to print_log in the C code example
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

	// Create GLSL program with loaded shaders
	program = gl.CreateProgram()
	program.AttachShader(vs)
	program.AttachShader(fs)
	program.Link()
	// TODO: Check for linking errors

	// Get the attribute location from the GLSL program (here from the vertex shader)
	attributeName := "coord2d"
	attributeCoord2d = program.GetAttribLocation(attributeName)
	if attributeCoord2d == -1 {
		fmt.Printf("Could not bind attribute %s\n", attributeName)
	}

	// Generate a buffer for the VertexBufferObject
	vboTriangle = gl.GenBuffer()
	vboTriangle.Bind(gl.ARRAY_BUFFER)
	// Submit the vertices of the triangle to the graphic card
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleVertices)*4, triangleVertices, gl.STATIC_DRAW)
	// Unset the active buffer
	gl.Buffer(0).Bind(gl.ARRAY_BUFFER)
}

func main() {
	var err os.Error
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

	initStatus := gl.Init() // Init glew
	if initStatus != 0 {
		fmt.Printf("Error-code: %d Init-Status: %d\n", gl.GetError(), initStatus)
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
	program.Delete()
	vboTriangle.Delete()
}

func display() {
	// Clear the background as white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Use the GLSL program
	program.Use()

	vboTriangle.Bind(gl.ARRAY_BUFFER)
	attributeCoord2d.EnableArray()

	// Describe our vertices array to OpenGL (it can't guess its format automatically)
	attributeCoord2d.AttribPointerOffset(2, gl.FLOAT, false, 0, 0)

	// Push each element in buffer_vertices to the vertex shader
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	attributeCoord2d.DisableArray()

	// Display the result
	glfw.SwapBuffers()
}
