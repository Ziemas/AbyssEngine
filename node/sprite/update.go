package sprite

import ren "github.com/OpenDiablo2/AbyssEngine/renderer"

func (s *Sprite) update(elapsed float64) {
	if s.onMouseButtonUp != nil || s.onMouseButtonDown != nil || s.onMouseOver != nil || s.onMouseLeave != nil {
		mx, my := s.mousePosProvider.GetMousePosition()
		posX, posY := s.GetPosition()
		mouseIsOver := mx >= posX && my >= posY && mx < (posX+int(s.textures[s.CurrentFrame].Width)) && my < (posY+int(s.textures[s.CurrentFrame].Height))

		if ren.IsMouseButtonDown(ren.MouseLeftButton) {
			if !s.isPressed {
				if s.canPress && mouseIsOver {

					s.isPressed = true

					if s.onMouseButtonDown != nil {
						s.onMouseButtonDown()
					}
				} else {
					s.canPress = false
				}
			}

		} else {
			if s.isPressed {
				s.isPressed = false

				if mouseIsOver {
					if s.onMouseButtonUp != nil {
						s.onMouseButtonUp()
					}
				}
			}
			s.canPress = true
		}

		if mouseIsOver && !s.isMouseOver {
			s.isMouseOver = true
			if s.onMouseOver != nil {
				s.onMouseOver()
			}
		} else if !mouseIsOver && s.isMouseOver {
			s.isMouseOver = false
			if s.onMouseLeave != nil {
				s.onMouseLeave()
			}
		}
	}

	s.animate(elapsed)

	if s.textures[s.CurrentFrame].ID == 0 {
		s.initializeTexture()
	}
}
