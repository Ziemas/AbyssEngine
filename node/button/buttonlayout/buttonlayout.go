package buttonlayout

type ButtonLayout struct {
	ResourceName     string
	PaletteName      string
	FontPath         string
	XSegments        int
	YSegments        int
	BaseFrame        int
	DisabledFrame    int
	DisabledColor    uint32
	TextOffsetX      int
	TextOffsetY      int
	FixedWidth       int
	FixedHeight      int
	LabelColor       uint32
	Toggleable       bool
	AllowFrameChange bool
	HasImage         bool
	Tooltip          int
	TooltipXOffset   int
	TooltipYOffset   int
}
