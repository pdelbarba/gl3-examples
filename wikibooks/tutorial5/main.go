package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
	"os"

	"math3d"

	"gl"
	"github.com/jteeuwen/glfw"
)

const (
	WindowTitle = "OpenGL 3 tutorial 5 - Adding 3rd dimension: a cube, plus a camera"
)

var ScreenHeight = 600
var ScreenWidth = 800

const (
	VertexShaderType = iota
	FragmentShaderType
)

var cubeVertices = []float32{
	// front
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	// back
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
}

var cubeColors = []float32{
	// front colors
	1.0, 0.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 1.0,
	1.0, 1.0, 1.0,
	// back colors
	1.0, 0.0, 0.0,
	0.0, 1.0, 0.0,
	0.0, 0.0, 1.0,
	1.0, 1.0, 1.0,
}

var cubeElements = []gl.GLushort{
	// front
	0, 1, 2,
	2, 3, 0,
	// top
	1, 5, 6,
	6, 2, 1,
	// back
	7, 6, 5,
	5, 4, 7,
	// bottom
	4, 0, 3,
	3, 7, 4,
	// left
	4, 5, 1,
	1, 0, 4,
	// right
	3, 2, 6,
	6, 7, 3,
}

var vboCubeVertices gl.Buffer
var vboCubeColors gl.Buffer
var iboCubeElements gl.Buffer

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
	vs, err = createShader("cube.v.glsl", VertexShaderType)
	if err != nil {
		fmt.Printf("Shader: %s\n", err)
		return
	}
	fs, err = createShader("cube.f.glsl", FragmentShaderType)
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

	vboCubeColors = gl.GenBuffer()
	vboCubeColors.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeColors)*4, cubeColors, gl.STATIC_DRAW)
	gl.Buffer(0).Bind(gl.ARRAY_BUFFER)

	vboCubeVertices = gl.GenBuffer()
	vboCubeVertices.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, cubeVertices, gl.STATIC_DRAW)
	gl.Buffer(0).Bind(gl.ARRAY_BUFFER)

	// Generate a buffer for the IndexBufferObject
	iboCubeElements = gl.GenBuffer()
	iboCubeElements.Bind(gl.ELEMENT_ARRAY_BUFFER)
	// Submit the indexes to the graphic card
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cubeElements)*2, cubeElements, gl.STATIC_DRAW)
	// Unset the active buffer
	gl.Buffer(0).Bind(gl.ELEMENT_ARRAY_BUFFER)

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

	uniformName := "mvp"
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
	err = glfw.OpenWindow(ScreenWidth, ScreenHeight, 0, 0, 0, 8, 8, 0, glfw.Windowed)
	if err != nil {
		fmt.Printf("GLFW: %s\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetWindowTitle(WindowTitle)
	glfw.SetWindowSizeCallback(onResize)

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
	gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	initResources()

	// We are limiting the calls to display() (frames per second) to 60. This prevents the 100% cpu usage.
	ticker := time.NewTicker(int64(second) / 60) // max 60 fps
	for {
		<-ticker.C
		angle := float32(glfw.Time())
		anim := math3d.MakeYRotationMatrix(angle)
		model := math3d.MakeTranslationMatrix(0, 0, -4)
		view := math3d.MakeLookAtMatrix(math3d.Vector3{0, 2, 0}, math3d.Vector3{0, 0, -4}, math3d.Vector3{0, 1, 0})
		projection := math3d.MakePerspectiveMatrix(45, float32(ScreenWidth)/float32(ScreenHeight), 0.1, 10.0)
		matrix = math3d.MakeIdentity().Multiply(projection).Multiply(view).Multiply(model).Multiply(anim)
		program.Use()
		uniformMTransform.UniformMatrix4fv(1, false, matrix.Transposed())
		display()
	}

	// Free resources
	free()

	runtime.UnlockOSThread()
}

func onResize(w, h int) {
	ScreenWidth = w
	ScreenHeight = h
	gl.Viewport(0, 0, ScreenWidth, ScreenHeight)
}

func free() {
	program.Delete()
	vboCubeColors.Delete()
	vboCubeVertices.Delete()
	iboCubeElements.Delete()
}
// 
var matrix = math3d.MakeIdentity()

func display() {
	// Clear the background as white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Use the GLSL program
	program.Use()

	attributeCoord3d.EnableArray()
	vboCubeVertices.Bind(gl.ARRAY_BUFFER)
	attributeCoord3d.AttribPointerOffset(3, gl.FLOAT, false, 0, 0)

	attributeColor.EnableArray()
	vboCubeColors.Bind(gl.ARRAY_BUFFER)
	attributeColor.AttribPointerOffset(3, gl.FLOAT, false, 0, 0)

	iboCubeElements.Bind(gl.ELEMENT_ARRAY_BUFFER)
	//size := []int32{0}
	//gl.GetBufferParameteriv(gl.ELEMENT_ARRAY_BUFFER, gl.BUFFER_SIZE, size)
	gl.DrawElementsOffset(gl.TRIANGLES, len(cubeElements), gl.UNSIGNED_SHORT, 0)

	attributeCoord3d.DisableArray()
	attributeColor.DisableArray()

	// Display the result
	glfw.SwapBuffers()
}
