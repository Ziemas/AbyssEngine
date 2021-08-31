package renderer

import (
	"image"
	"image/draw"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Texture struct {
	ID            uint32
	Width, Height int
}

func DrawTextureP(texture Texture, posX, posY int, palette Texture) {
	model := mgl32.Ident4().
		Mul4(mgl32.Translate3D(float32(posX), float32(posY), 0)).
		Mul4(mgl32.Scale3D(float32(texture.Width), float32(texture.Height), 0.0))

	gl.UniformMatrix4fv(int32(UniformModelLoc), 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0)
	texture.Bind()

	gl.ActiveTexture(uint32(gl.TEXTURE1))
	palette.Bind()

	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)

}

func DrawTexture(texture Texture, posX, posY int) {

	model := mgl32.Ident4().
		Mul4(mgl32.Translate3D(float32(posX), float32(posY), 0)).
		Mul4(mgl32.Scale3D(float32(texture.Width), float32(texture.Height), 0.0))

	gl.UniformMatrix4fv(int32(UniformModelLoc), 1, false, &model[0])

	gl.ActiveTexture(gl.TEXTURE0)
	texture.Bind()
	gl.BindVertexArray(quadVAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)

}

func UnloadTexture(t Texture) {
	gl.DeleteTextures(1, &t.ID)
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.ID)
}

func NewTextureIndexed(bytes []byte, width, height int) Texture {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		gl.R8,
		int32(width),
		int32(height),
		0,
		gl.RED,
		gl.UNSIGNED_BYTE,
		gl.Ptr(bytes))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	tex := Texture{
		ID:     texture,
		Width:  width,
		Height: height,
	}

	return tex
}

func NewTextureRGBABytes(bytes []byte, width, height int) Texture {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(width),
		int32(height),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(bytes))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	tex := Texture{
		ID:     texture,
		Width:  width,
		Height: height,
	}

	return tex
}

func NewTextureRGBA(img image.Image) Texture {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(img.Bounds().Size().X),
		int32(img.Bounds().Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.BindTexture(gl.TEXTURE_2D, 0)

	tex := Texture{
		ID:     texture,
		Width:  img.Bounds().Size().X,
		Height: img.Bounds().Size().Y,
	}

	return tex
}
