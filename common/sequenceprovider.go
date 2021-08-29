package common

type SequenceProvider interface {
	SequenceCount() int
	FrameCount(sequenceId int) int
	FrameWidth(sequenceId, frameId int) int
	FrameHeight(sequenceId, frameId int) int
	GetColorIndexAt(sequenceId, frameId, x, y int) uint8
	GetFrameOffsetX(sequenceId, frameId int) int
	GetFrameOffsetY(sequenceId, frameId int) int
}
