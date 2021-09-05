package zrenderprovider

import (
	"errors"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type zTexture struct {
	id            uint32
	format        PixelFormat
	width, height int
	fbo           uint32
}

func textureFormat(fmt PixelFormat) (glInternalFmt int, glFormat int, glType int, err error) {
	switch fmt {
	case PixelFmtGrayscale:
		return gl.R8, gl.RED, gl.UNSIGNED_BYTE, nil
	case PixelFmtRGBA8:
		return gl.RGBA8, gl.RGBA, gl.UNSIGNED_BYTE, nil
	}

	return 0, 0, 0, errors.New("Unsupported pixel format")
}

func textureSwizzle(fmt PixelFormat) []int32 {
	switch fmt {
	case PixelFmtGrayscale:
		return []int32{gl.RED, gl.RED, gl.RED, gl.ONE}
	default:
		return []int32{gl.RED, gl.GREEN, gl.BLUE, gl.ALPHA}
	}
}

func NewTexture(pixels unsafe.Pointer, width, height int, format PixelFormat) (*zTexture, error) {
	glifmt, glfmt, gltype, err := textureFormat(format)

	if err != nil {
		return nil, err
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	println(texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		int32(glifmt),
		int32(width),
		int32(height),
		0,
		uint32(glfmt),
		uint32(gltype),
		pixels)

	swizzle := textureSwizzle(format)
	gl.TexParameteriv(gl.TEXTURE_2D, gl.TEXTURE_SWIZZLE_RGBA, &swizzle[0])

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	tex := zTexture{
		id:     texture,
		format: format,
		width:  width,
		height: height,
	}

	return &tex, nil
}

func NewRenderTexture(width, height int) (*zTexture, error) {
	tex, err := NewTexture(nil, width, height, PixelFmtRGBA8)
	if err != nil {
		return nil, err
	}

	gl.GenFramebuffers(1, &tex.fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, tex.fbo)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, tex.id, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return tex, nil
}

func NewTextureFromImage(img *zImage) (*zTexture, error) {
	return NewTexture(gl.Ptr(img.data), img.width, img.height, img.fmt)
}

func (t zTexture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.id)
}

func (t zTexture) Width() int {
	return t.width
}

func (t zTexture) Height() int {
	return t.height
}

func (t zTexture) ID() int {
	return int(t.id)
}
