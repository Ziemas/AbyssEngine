package sprite

type playMode int

const (
	playModePause playMode = iota
	playModeForward
	playModeBackward
)
const defaultPlayLength = 1.0

func (s *Sprite) PlayForward() {
	s.playMode = playModeForward
	s.lastFrameTime = 0
}

func (s *Sprite) animate(elapsed float64) {
	if s.playMode == playModePause {
		return
	}

	frameCount := s.Sequences.FrameCount(s.CurrentSequence())
	frameLength := s.playLength / float64(frameCount)
	s.lastFrameTime += elapsed
	framesAdvanced := int(s.lastFrameTime / frameLength)
	s.lastFrameTime -= float64(framesAdvanced) * frameLength

	for i := 0; i < framesAdvanced; i++ {
		s.advanceFrame()
	}

}

func (s *Sprite) advanceFrame() {
	startIndex := 0
	endIndex := s.Sequences.FrameCount(s.CurrentSequence())

	if s.hasSubLoop && s.playedCount > 0 {
		startIndex = s.subStartingFrame
		endIndex = s.subEndingFrame
	}

	switch s.playMode {
	case playModeForward:
		s.CurrentFrame++
		if s.CurrentFrame >= endIndex {
			s.playedCount++
			if s.playLoop {
				s.CurrentFrame = startIndex
			} else {
				s.CurrentFrame = endIndex - 1
				break
			}
		}
	case playModeBackward:
		s.CurrentFrame--
		if s.CurrentFrame < startIndex {
			s.playedCount++
			if s.playLoop {
				s.CurrentFrame = endIndex - 1
			} else {
				s.CurrentFrame = startIndex
				break
			}
		}
	}
}
