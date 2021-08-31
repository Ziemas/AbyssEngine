package sprite

import (
	"github.com/OpenDiablo2/AbyssEngine/common"
	ren "github.com/OpenDiablo2/AbyssEngine/renderer"
)

func (s *Sprite) render() {
	if s.textures[s.CurrentFrame].ID == 0 || !s.Visible || !s.Active {
		return
	}

	tex := common.PaletteTexture[s.palette]
	if !tex.Init {
		//img := rl.NewImage(tex.Data, 256, int32(common.PaletteTransformsCount), 1, rl.UncompressedR8g8b8a8)
		//tex.Texture = rl.LoadTextureFromImage(img)
		tex.Texture = ren.NewTextureRGBABytes(tex.Data, 256, common.PaletteTransformsCount)

		tex.Init = true
	}

	posX, posY := s.GetPosition()

	posX += s.Sequences.FrameOffsetX(s.CurrentSequence(), s.CurrentFrame)

	if s.bottomOrigin {
		posY -= s.Sequences.FrameHeight(s.CurrentSequence(), s.CurrentFrame, s.CellSizeX, s.CellSizeY)
	}

	posY += s.Sequences.FrameOffsetY(s.CurrentSequence(), s.CurrentFrame)

	//rl.SetShaderValueTexture(common.PaletteShader, common.PaletteShaderLoc, tex.Texture)
	//rl.SetShaderValue(common.PaletteShader, common.PaletteShaderOffsetLoc, []float32{float32(s.paletteShift) / float32(common.PaletteTransformsCount-1)}, rl.ShaderUniformFloat)
	//ren.SetShaderValuei(common.PaletteShaderLoc, int32(tex.Texture.ID))
	ren.SetShaderValueF(common.PaletteShaderOffsetLoc, float32(s.paletteShift)/float32(common.PaletteTransformsCount-1))
	s.blendModeProvider.SetBlendMode(s.blendMode)
	//rl.DrawTexture(s.textures[s.CurrentFrame], int32(posX), int32(posY), rl.White)
	ren.DrawTextureP(s.textures[s.CurrentFrame], posX, posY, tex.Texture, common.PaletteShaderLoc)

}

func (s *Sprite) initializeTexture() {
	width := s.Sequences.FrameWidth(s.CurrentSequence(), s.CurrentFrame, s.CellSizeX)
	height := s.Sequences.FrameHeight(s.CurrentSequence(), s.CurrentFrame, s.CellSizeX, s.CellSizeY)

	pixels := make([]byte, width*height)

	targetStartX := 0
	targetStartY := 0

	for cellOffsetY := 0; cellOffsetY < s.CellSizeY; cellOffsetY++ {
		for cellOffsetX := 0; cellOffsetX < s.CellSizeX; cellOffsetX++ {
			cellIndex := s.CurrentFrame + (cellOffsetX + (cellOffsetY * s.CellSizeX))

			frameWidth := s.Sequences.FrameWidth(s.CurrentSequence(), cellIndex, 1)
			frameHeight := s.Sequences.FrameHeight(s.CurrentSequence(), cellIndex, 1, 1)

			for y := 0; y < frameHeight; y++ {
				idx := targetStartX + ((targetStartY + y) * width)
				for x := 0; x < frameWidth; x++ {
					c := s.Sequences.GetColorIndexAt(s.CurrentSequence(), cellIndex, x, y)

					pixels[idx] = c
					idx++
				}
			}

			targetStartX += frameWidth
		}

		targetStartX = 0
		targetStartY += s.Sequences.FrameHeight(s.CurrentSequence(), cellOffsetY, 1, 1)
	}

	//img := rl.NewImage(pixels, int32(width), int32(height), 1, rl.UncompressedGrayscale)

	//s.textures[s.CurrentFrame] = rl.LoadTextureFromImage(img)
	s.textures[s.CurrentFrame] = ren.NewTextureIndexed(pixels, width, height)
}
