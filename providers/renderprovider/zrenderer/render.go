package zrenderprovider

import (
	"container/list"
	"image"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

const (
	TexUnitImage int = iota
	TexUnitPalette
)

type BlendMode int
const (
	BlendModeNone BlendMode = iota
	BlendModeAlpha
	BlendModeAdditive
	BlendModeMultiplied
	BlendModeAddColors
	BlendModeSubtractColors
)

func SetBlendMode(mode BlendMode) {
	switch mode {
	case BlendModeNone:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	case BlendModeAdditive:
		gl.BlendFunc(gl.ONE, gl.ONE)
	case BlendModeMultiplied:
		gl.BlendFunc(gl.DST_COLOR, gl.ONE_MINUS_SRC_ALPHA)
	default:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	}
}

var (
	verts = []float32{
		0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0,
	}
)

var (
	projection = mgl.Ortho(0.0, 800.0, 600.0, 0.0, -1.0, 1.0)
	fbw, fbh   int
)

type Renderer struct {
	drawList *list.List
	vbo      uint32
	fanVAO   uint32
	shader   Shader
	fbo      uint32
	palette  map[string]*Palette
}

func NewRenderer() *Renderer {
	ren := Renderer{
		palette:  make(map[string]*Palette),
	}

	err := gl.Init()
	if err != nil {
		panic(err)
	}

	gl.DebugMessageCallback(debugCb, nil)
	gl.Enable(gl.DEBUG_OUTPUT)

	gl.GenVertexArrays(1, &ren.fanVAO)
	gl.GenBuffers(1, &ren.vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, ren.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)
	gl.BindVertexArray(ren.fanVAO)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	gl.ClearColor(0, 0, 0, 255)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	ren.shader = InitShaders()
	ren.InitUniforms()
	return &ren
}

func (r *Renderer) InitUniforms() {
	sh := r.shader
	gl.UseProgram(sh.id)
	gl.UniformMatrix4fv(sh.UniformProjection, 1, false, &projection[0])
	gl.Uniform1i(sh.UniformImage, int32(TexUnitImage))
	gl.Uniform1i(sh.UniformPaletteTex, int32(TexUnitImage))
	gl.Uniform1i(sh.UniformPaletteOffset, 0)
	gl.Uniform1i(sh.UniformUsePalette, 0)
}

func (r *Renderer) LoadPalette(name string, pal *Palette) {
	r.palette[name] = pal
}

func (r *Renderer) BindRenderTexture(tex *zTexture) {
	localProj := mgl.Ortho(0, float32(tex.width), 0, float32(tex.height), 1.0, -1.0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, tex.fbo)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Viewport(0, 0, int32(tex.width), int32(tex.height))
	gl.UniformMatrix4fv(r.shader.UniformProjection, 1, false, &localProj[0])
}

func (r *Renderer) UnbindRenderTexture() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.UniformMatrix4fv(r.shader.UniformProjection, 1, false, &projection[0])
	gl.Viewport(0, 0, int32(fbw), int32(fbh))
}

func (r *Renderer) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (r *Renderer) setupPalette(palette string, paletteOffset int) error {
	if palette != "" {
		pal := r.palette[palette]
		if !pal.Init {
			paltex, err := NewTexture(gl.Ptr(pal.Data), 256, PaletteTransformsCount, PixelFmtRGBA8)
			if err != nil {
				return err
			}
			pal.Texture = paltex
			pal.Init = true
		}

		gl.Uniform1i(r.shader.UniformUsePalette, 1)
		gl.Uniform1i(r.shader.UniformPaletteTex, int32(TexUnitPalette))
		gl.Uniform1i(r.shader.UniformPaletteOffset, int32(paletteOffset))
		gl.ActiveTexture(gl.TEXTURE0 + uint32(TexUnitPalette))
		pal.Texture.Bind()
	} else {
		gl.Uniform1i(r.shader.UniformUsePalette, 0)
	}

	return nil
}

func (r *Renderer) DrawTextureEX(tex *zTexture, srcRect, destRect image.Rectangle, palette string, paletteOffset int) error {
	// TODO source
	model := mgl.Ident4().
		Mul4(mgl.Translate3D(float32(destRect.Min.X), float32(destRect.Min.Y), 0)).
		Mul4(mgl.Scale3D(float32(destRect.Dx()), float32(destRect.Dy()), 0.0))

	err := r.setupPalette(palette, paletteOffset)
	if err != nil {
		return err
	}

	gl.UniformMatrix4fv(int32(r.shader.UniformModel), 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0 + uint32(TexUnitImage))
	tex.Bind()
	gl.BindVertexArray(r.fanVAO)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return nil
}

func (r *Renderer) DrawTexture(tex *zTexture, x, y int, palette string, paletteOffset int) error {
	model := mgl.Ident4().
		Mul4(mgl.Translate3D(float32(x), float32(y), 0)).
		Mul4(mgl.Scale3D(float32(tex.Width()), float32(tex.Height()), 0.0))
	err := r.setupPalette(palette, paletteOffset)
	if err != nil {
		return err
	}

	gl.UniformMatrix4fv(int32(r.shader.UniformModel), 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0 + uint32(TexUnitImage))
	tex.Bind()
	gl.BindVertexArray(r.fanVAO)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	return nil
}

func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	fbw = width
	fbh = height
	projection = mgl.Ortho(0.0, float32(width), float32(height), 0.0, -1.0, 1.0)
}
