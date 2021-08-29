package common

import dc6 "github.com/OpenDiablo2/dc6/pkg"

type DC6SequenceProvider struct {
	Sequences []*dc6.Direction
}

func (d *DC6SequenceProvider) SequenceCount() int {
	return len(d.Sequences)
}

func (d *DC6SequenceProvider) FrameCount(sequenceId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return len(d.Sequences[sequenceId].Frames)
}

func (d *DC6SequenceProvider) FrameWidth(sequenceId, frameId, frameSizeX int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames) {
		return 0
	}

	width := 0
	for i := 0; i < frameSizeX; i++ {
		width += int(d.Sequences[sequenceId].Frames[frameId+i].Width)
	}

	return width
}

func (d *DC6SequenceProvider) FrameHeight(sequenceId, frameId, frameSizeX, frameSizeY int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames) {
		return 0
	}

	height := 0
	for i := 0; i < frameSizeY; i++ {
		height += int(d.Sequences[sequenceId].Frames[frameId+(i*frameSizeX)].Height)
	}

	return height
}

func (d *DC6SequenceProvider) GetColorIndexAt(sequenceId, frameId, x, y int) uint8 {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	if frameId < 0 || frameId >= len(d.Sequences[sequenceId].Frames) {
		return 0
	}

	return d.Sequences[sequenceId].Frames[frameId].ColorIndexAt(x, y)
}

func (d *DC6SequenceProvider) FrameOffsetX(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return int(d.Sequences[sequenceId].Frames[frameId].OffsetX)
}

func (d *DC6SequenceProvider) FrameOffsetY(sequenceId, frameId int) int {
	if sequenceId < 0 || sequenceId >= len(d.Sequences) {
		return 0
	}

	return int(d.Sequences[sequenceId].Frames[frameId].OffsetY)
}
