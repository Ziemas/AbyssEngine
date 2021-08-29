package common

import (
	dcc "github.com/OpenDiablo2/dcc/pkg"
)

type DCCSequenceProvider struct {
	Sequences []*dcc.Direction
}

func (d *DCCSequenceProvider) FrameOffsetX(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return d.Sequences[sequenceId].Frame(frameId).XOffset
}

func (d *DCCSequenceProvider) FrameOffsetY(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return d.Sequences[sequenceId].Frame(frameId).YOffset
}

func (d *DCCSequenceProvider) SequenceCount() int {
	return len(d.Sequences)
}

func (d *DCCSequenceProvider) FrameCount(sequenceId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return len(d.Sequences[sequenceId].Frames())
}

func (d *DCCSequenceProvider) FrameWidth(sequenceId, frameId, frameSizeX int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames()) {
		return 0
	}

	width := 0
	for i := 0; i < frameSizeX; i++ {
		width += int(d.Sequences[sequenceId].Frames()[frameId+i].Width)
	}

	return width
}

func (d *DCCSequenceProvider) FrameHeight(sequenceId, frameId, frameSizeX, frameSizeY int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames()) {
		return 0
	}

	height := 0
	for i := 0; i < frameSizeY; i++ {
		height += d.Sequences[sequenceId].Frames()[frameId+(i*frameSizeX)].Height
	}

	return height
}

func (d *DCCSequenceProvider) GetColorIndexAt(sequenceId, frameId, x, y int) uint8 {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames()) {
		return 0
	}

	return d.Sequences[sequenceId].Frame(frameId).ColorIndexAt(x, y)
}
