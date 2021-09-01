package sprite

import (
	"errors"
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
	"io/ioutil"
	"path"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
	dcc "github.com/OpenDiablo2/dcc/pkg"
)

type Sprite struct {
	*node.Node

	mousePosProvider  common.MousePositionProvider
	Sequences         common.SequenceProvider
	renderProvider    renderprovider.RenderProvider
	palette           string
	currentSequence   int
	CurrentFrame      int
	Visible           bool
	CellSizeX         int
	CellSizeY         int
	isPressed         bool
	isMouseOver       bool
	canPress          bool
	textures          []renderprovider.Texture
	lastFrameTime     float64
	playedCount       int
	playMode          playMode
	playLength        float64
	hasSubLoop        bool
	subStartingFrame  int
	subEndingFrame    int
	playLoop          bool
	bottomOrigin      bool
	blendMode         renderprovider.BlendMode
	paletteShift      int
	onMouseButtonDown func()
	onMouseButtonUp   func()
	onMouseOver       func()
	onMouseLeave      func()
}

func New(loaderProvider common.LoaderProvider, mousePosProvider common.MousePositionProvider,
	renderProvider renderprovider.RenderProvider, filePath, palette string) (*Sprite, error) {
	result := &Sprite{
		Node:             node.New(),
		mousePosProvider: mousePosProvider,
		renderProvider:   renderProvider,
		Visible:          true,
		currentSequence:  0,
		CurrentFrame:     0,
		CellSizeX:        1,
		CellSizeY:        1,
		textures:         make([]renderprovider.Texture, 0),
		isPressed:        false,
		isMouseOver:      false,
		canPress:         true,
		playMode:         playModePause,
		playLength:       defaultPlayLength,
		playedCount:      0,
		lastFrameTime:    0,
		paletteShift:     0,
		bottomOrigin:     false,
		blendMode:        renderprovider.BlendModeNone,
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
	defer func() { _ = fileStream.Close() }()

	if err != nil {
		return nil, err
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

	result.textures = make([]renderprovider.Texture, result.Sequences.FrameCount(result.CurrentSequence()))
	return result, nil
}

func (s *Sprite) CurrentSequence() int {
	return s.currentSequence
}

func (s *Sprite) SetSequence(seqId int) {
	if seqId < 0 || seqId >= s.Sequences.SequenceCount() {
		return
	}

	for texIdx := range s.textures {
		_ = s.renderProvider.FreeTexture(s.textures[texIdx])
	}

	s.currentSequence = seqId
	s.textures = make([]renderprovider.Texture, s.Sequences.FrameCount(s.CurrentSequence()))
}

func (s *Sprite) setPalette(palette string) {
	s.palette = palette
}

func (s *Sprite) Destroy() {
	s.ShouldRemove = true
	s.Active = false

	for idx := range s.textures {
		_ = s.renderProvider.FreeTexture(s.textures[idx])
	}

}
