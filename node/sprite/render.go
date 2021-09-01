package sprite

import (
	"bytes"
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
)

func (s *Sprite) render() {
	if s.textures[s.CurrentFrame] == nil || !s.Visible || !s.Active {
		return
	}

	posX, posY := s.GetPosition()

	posX += s.Sequences.FrameOffsetX(s.CurrentSequence(), s.CurrentFrame)

	if s.bottomOrigin {
		posY -= s.Sequences.FrameHeight(s.CurrentSequence(), s.CurrentFrame, s.CellSizeX, s.CellSizeY)
	}

	posY += s.Sequences.FrameOffsetY(s.CurrentSequence(), s.CurrentFrame)

	s.renderProvider.BeginBlendMode(s.blendMode)
	_ = s.renderProvider.DrawTexture(s.textures[s.CurrentFrame], posX, posY, s.palette, s.paletteShift)
	s.renderProvider.EndBlendMode()

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

	img, _ := s.renderProvider.NewImage(bytes.NewReader(pixels), width, height, renderprovider.ImageColorModeGrayscale)

	s.textures[s.CurrentFrame], _ = s.renderProvider.LoadTextureFromImage(img)
}
