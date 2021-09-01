package sprite

import (
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
)

func (s *Sprite) update(elapsed float64) {
	if s.textures[s.CurrentFrame] == nil {
		s.initializeTexture()
	}

	if s.OnMouseButtonUp != nil || s.OnMouseButtonDown != nil || s.OnMouseOver != nil || s.OnMouseLeave != nil {
		mx, my := s.mousePosProvider.GetMousePosition()
		posX, posY := s.GetPosition()

		mouseIsOver := mx >= posX && my >= posY && mx < (posX+int(s.textures[s.CurrentFrame].Width())) && my < (posY+int(s.textures[s.CurrentFrame].Height()))

		if s.renderProvider.IsMouseButtonPressed(renderprovider.MouseButtonLeft) {
			if !s.isPressed {
				if s.canPress && mouseIsOver {

					s.isPressed = true

					if s.OnMouseButtonDown != nil {
						s.OnMouseButtonDown()
					}
				} else {
					s.canPress = false
				}
			}

		} else {
			if s.isPressed {
				s.isPressed = false

				if mouseIsOver {
					if s.OnMouseButtonUp != nil {
						s.OnMouseButtonUp()
					}
				}
			}
			s.canPress = true
		}

		if mouseIsOver && !s.isMouseOver {
			s.isMouseOver = true
			if s.OnMouseOver != nil {
				s.OnMouseOver()
			}
		} else if !mouseIsOver && s.isMouseOver {
			s.isMouseOver = false
			if s.OnMouseLeave != nil {
				s.OnMouseLeave()
			}
		}
	}

	s.animate(elapsed)

	if s.textures[s.CurrentFrame] == nil {
		s.initializeTexture()
	}
}
