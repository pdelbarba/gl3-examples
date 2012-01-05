package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"runtime"
	"time"
	"os"
	
	"math3d"

	"gl"
	"github.com/jteeuwen/glfw"
)

const (
	ScreenHeight = 480
	ScreenWidth  = 640
	WindowTitle  = "OpenGL 3 tutorial 4 - transformation matrices"
)

const (
	VertexShaderType = iota
	FragmentShaderType
)

var triangleAttributes = []float32{
	-0.5, -0.5, 0.0,
	1.0, 1.0, 0.0,
	0.5, -0.5, 0.0,
	0.0, 0.0, 1.0,
	0.0, 0.5, 0.0,
	1.0, 0.0, 0.0,
}

var vboTriangle gl.Buffer

var vs gl.Shader
var fs gl.Shader
var program gl.Program

var attributeCoord3d gl.AttribLocation
var attributeColor gl.AttribLocation
var uniformMTransform gl.UniformLocation

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
	infoLog := program.GetInfoLog()
	if len(infoLog) != 0 {
		fmt.Printf("Program: %s\n", infoLog)
	}

	// Generate a buffer for the VertexBufferObject
	vboTriangle = gl.GenBuffer()
	vboTriangle.Bind(gl.ARRAY_BUFFER)
	// Submit the vertices of the triangle to the graphic card
	gl.BufferData(gl.ARRAY_BUFFER, len(triangleAttributes)*4, triangleAttributes, gl.STATIC_DRAW)
	// Unset the active buffer
	gl.Buffer(0).Bind(gl.ARRAY_BUFFER)

	// Get the attribute location from the GLSL program (here from the vertex shader)
	attributeName := "coord3d"
	attributeCoord3d = program.GetAttribLocation(attributeName)
	if attributeCoord3d == -1 {
		fmt.Printf("Could not bind attribute %s\n", attributeName)
	}

	attributeName = "v_color"
	attributeColor = program.GetAttribLocation(attributeName)
	if attributeColor == -1 {
		fmt.Printf("Could not bind attribute %s\n", attributeName)
	}

	uniformName := "m_transform"
	uniformMTransform = program.GetUniformLocation(uniformName)
	if uniformMTransform == -1 {
		fmt.Printf("Could not bind attribute %s\n", uniformName)
	}
}

const second = 1e9

func main() {
	// We need to lock the goroutine to one thread due time.Ticker
	runtime.LockOSThread()

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
	// We are limiting the calls to display() (frames per second) to 60. This prevents the 100% cpu usage.
	ticker := time.NewTicker(int64(second) / 60) // max 60 fps
	for {
		<-ticker.C
		move := float32(math.Sin(glfw.Time()))
		angle := float32(glfw.Time())
		matrix = math3d.MakeTranslationMatrix(move, 0.0, 0.0)
		matrix = matrix.Multiply(math3d.MakeZRotationMatrix(angle)).Transposed()
		display()
	}

	// Free resources
	free()

	runtime.UnlockOSThread()
}

func free() {
	program.Delete()
	vboTriangle.Delete()
}

var matrix = math3d.MakeIdentity()

func display() {
	// Clear the background as white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Use the GLSL program
	program.Use()

	uniformMTransform.UniformMatrix4fv(1, false, matrix)

	vboTriangle.Bind(gl.ARRAY_BUFFER)

	attributeCoord3d.EnableArray()
	// Describe our vertices array to OpenGL (it can't guess its format automatically)
	attributeCoord3d.AttribPointerOffset(3, gl.FLOAT, false, 6*4, 0)

	attributeColor.EnableArray()
	attributeColor.AttribPointerOffset(3, gl.FLOAT, false, 6*4, 3*4)

	// Push each element in buffer_vertices to the vertex shader
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	attributeCoord3d.DisableArray()
	attributeColor.DisableArray()

	// Display the result
	glfw.SwapBuffers()
}
