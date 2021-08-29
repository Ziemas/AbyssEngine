package common

type SequenceProvider interface {
	SequenceCount() int
	FrameCount(sequenceId int) int
	FrameWidth(sequenceId, frameId, frameSizeX int) int
	FrameHeight(sequenceId, frameId, frameSizeX, frameSizeY int) int
	GetColorIndexAt(sequenceId, frameId, x, y int) uint8
	FrameOffsetX(sequenceId, frameId int) int
	FrameOffsetY(sequenceId, frameId int) int
}
