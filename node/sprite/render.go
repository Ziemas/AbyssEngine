package sprite

import (
	"errors"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type blendMode int

const (
	blendModeNone blendMode = iota
	blendModeAlpha
	blendModeAdditive
	blendModeMultiplied
	blendModeAddColors
	blendModeSubtractColors
)

var blendModeLookup = []rl.BlendMode{
	-1,
	rl.BlendAlpha,
	rl.BlendAdditive,
	rl.BlendMultiplied,
	rl.BlendAddColors,
	rl.BlendSubtractColors,
}

func blendModeToString(mode blendMode) string {
	switch mode {
	case blendModeNone:
		return ""
	case blendModeAlpha:
		return "alpha"
	case blendModeAdditive:
		return "add"
	case blendModeMultiplied:
		return "multiply"
	case blendModeAddColors:
		return "addcolors"
	case blendModeSubtractColors:
		return "subcolors"
	default:
		return ""
	}
}

func stringToBlendMode(mode string) (blendMode, error) {
	switch strings.ToLower(mode) {
	case "":
		return blendModeNone, nil
	case "alpha":
		return blendModeAlpha, nil
	case "add":
		return blendModeAdditive, nil
	case "multiply":
		return blendModeMultiplied, nil
	case "addcolors":
		return blendModeAddColors, nil
	case "subcolors":
		return blendModeSubtractColors, nil
	default:
		return -1, errors.New("invalid blend mode")
	}
}

func (s *Sprite) render() {
	if s.textures[s.CurrentFrame].ID == 0 || !s.Visible || !s.Active {
		return
	}

	tex := common.PaletteTexture[s.palette]
	if !tex.Init {
		img := rl.NewImage(tex.Data, 256, 1, 1, rl.UncompressedR8g8b8a8)
		tex.Texture = rl.LoadTextureFromImage(img)

		tex.Init = true
	}

	posX, posY := s.GetPosition()

	posX += s.Sequences.GetFrameOffsetX(s.CurrentSequence(), s.CurrentFrame)

	if s.CellSizeX == 1 && s.CellSizeY == 1 {
		posY -= s.Sequences.FrameHeight(s.CurrentSequence(), s.CurrentFrame)
	}

	posY += s.Sequences.GetFrameOffsetY(s.CurrentSequence(), s.CurrentFrame)

	rl.BeginShaderMode(common.PaletteShader)
	rl.SetShaderValueTexture(common.PaletteShader, common.PaletteShaderLoc, tex.Texture)

	if blendModeLookup[s.blendMode] != -1 {
		rl.BeginBlendMode(blendModeLookup[s.blendMode])
	}

	rl.DrawTexture(s.textures[s.CurrentFrame], int32(posX), int32(posY), rl.White)

	if blendModeLookup[s.blendMode] != -1 {
		rl.EndBlendMode()
	}

	rl.EndShaderMode()
}

func (s *Sprite) initializeTexture() {
	width := 0
	height := 0

	for i := 0; i < s.CellSizeX; i++ {
		width += s.Sequences.FrameWidth(s.CurrentSequence(), s.CurrentFrame+i)
	}

	for i := 0; i < s.CellSizeY; i++ {
		height += s.Sequences.FrameHeight(s.CurrentSequence(), s.CurrentFrame+(i*s.CellSizeX))
	}

	pixels := make([]byte, width*height)

	targetStartX := 0
	targetStartY := 0

	for cellOffsetY := 0; cellOffsetY < s.CellSizeY; cellOffsetY++ {
		for cellOffsetX := 0; cellOffsetX < s.CellSizeX; cellOffsetX++ {
			cellIndex := s.CurrentFrame + (cellOffsetX + (cellOffsetY * s.CellSizeX))

			frameWidth := s.Sequences.FrameWidth(s.CurrentSequence(), cellIndex)
			frameHeight := s.Sequences.FrameHeight(s.CurrentSequence(), cellIndex)

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
		targetStartY += s.Sequences.FrameHeight(s.CurrentSequence(), cellOffsetY*s.CellSizeX)
	}

	img := rl.NewImage(pixels, int32(width), int32(height), 1, rl.UncompressedGrayscale)

	s.textures[s.CurrentFrame] = rl.LoadTextureFromImage(img)
}
