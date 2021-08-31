package sprite

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	ren "github.com/OpenDiablo2/AbyssEngine/renderer"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
	dcc "github.com/OpenDiablo2/dcc/pkg"
)

type Sprite struct {
	*node.Node

	mousePosProvider  common.MousePositionProvider
	Sequences         common.SequenceProvider
	blendModeProvider common.BlendModeProvider
	palette           string
	currentSequence   int
	CurrentFrame      int
	Visible           bool
	CellSizeX         int
	CellSizeY         int
	isPressed         bool
	isMouseOver       bool
	canPress          bool
	//textures          []rl.Texture2D
	textures          []ren.Texture
	lastFrameTime     float64
	playedCount       int
	playMode          playMode
	playLength        float64
	hasSubLoop        bool
	subStartingFrame  int
	subEndingFrame    int
	playLoop          bool
	bottomOrigin      bool
	blendMode         common.BlendMode
	paletteShift      int
	onMouseButtonDown func()
	onMouseButtonUp   func()
	onMouseOver       func()
	onMouseLeave      func()
}

func New(loaderProvider common.LoaderProvider, mousePosProvider common.MousePositionProvider,
	blendModeProvider common.BlendModeProvider,
	filePath, palette string) (*Sprite, error) {
	result := &Sprite{
		Node:              node.New(),
		mousePosProvider:  mousePosProvider,
		blendModeProvider: blendModeProvider,
		Visible:           true,
		currentSequence:   0,
		CurrentFrame:      0,
		CellSizeX:         1,
		CellSizeY:         1,
		//textures:          make([]rl.Texture2D, 0),
		textures:          make([]ren.Texture, 0),
		isPressed:        false,
		isMouseOver:      false,
		canPress:         true,
		playMode:         playModePause,
		playLength:       defaultPlayLength,
		playedCount:      0,
		lastFrameTime:    0,
		paletteShift:     0,
		bottomOrigin:     false,
		blendMode:        common.BlendModeNone,
		subStartingFrame: 0,
		subEndingFrame:   0,
		hasSubLoop:       false,
		playLoop:         true,
		palette:          palette,
	}

	result.RenderCallback = result.render
	result.UpdateCallback = result.update

	fileExt := strings.ToLower(path.Ext(filePath))

	fileStream, err := loaderProvider.Load(filePath)
	defer fileStream.Close()

	if err != nil {
		return nil, err
	}

	_, ok := common.PaletteTexture[palette]
	if !ok {
		return nil, errors.New("sprite loaded with non-existent palette")
	}

	switch fileExt {
	case ".dcc":
		bytes, err := ioutil.ReadAll(fileStream)

		if err != nil {
			return nil, err
		}

		dccRes, err := dcc.FromBytes(bytes)

		if err != nil {
			return nil, err
		}

		result.Sequences = &common.DCCSequenceProvider{Sequences: dccRes.Directions()}

	case ".dc6":
		bytes, err := ioutil.ReadAll(fileStream)

		if err != nil {
			return nil, err
		}

		dc6Res, err := dc6.FromBytes(bytes)

		if err != nil {
			return nil, err
		}

		result.Sequences = &common.DC6SequenceProvider{Sequences: dc6Res.Directions}

	default:
		return nil, errors.New("unsupported file format")
	}

	//result.textures = make([]rl.Texture2D, result.Sequences.FrameCount(result.CurrentSequence()))
	result.textures = make([]ren.Texture, result.Sequences.FrameCount(result.CurrentSequence()))
	return result, nil
}

func (s *Sprite) CurrentSequence() int {
	return s.currentSequence
}

func (s *Sprite) SetSequence(seqId int) {
	if seqId < 0 || seqId >= s.Sequences.SequenceCount() {
		return
	}

	//for texIdx := range s.textures {
	//	rl.UnloadTexture(s.textures[texIdx])
	//}

	s.currentSequence = seqId
	//s.textures = make([]rl.Texture2D, s.Sequences.FrameCount(s.CurrentSequence()))
}

func (s *Sprite) setPalette(palette string) {
	s.palette = palette
}

func (s *Sprite) Destroy() {
	s.ShouldRemove = true
	s.Active = false

	//for idx := range s.textures {
	//	rl.UnloadTexture(s.textures[idx])
	//}

}
