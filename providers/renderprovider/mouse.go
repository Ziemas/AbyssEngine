package renderprovider

type MouseButton int

const (
	MouseButtonLeft MouseButton = 1 << iota
	MouseButtonMiddle
	MouseButtonRight
)
