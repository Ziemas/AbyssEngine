package zrenderprovider

type RGBA8 struct {
	r, g, b, a uint8
}

type Grayscale struct {
	r uint8
}

type zImage struct {
	fmt PixelFormat
	data []byte
	width, height int
}

func NewImage(data []byte, width, height int, fmt PixelFormat) *zImage {
	img := zImage{
		fmt:    fmt,
		data:   data,
		width:  width,
		height: height,
	}

	return &img
}

func (i zImage) Width() int {
	return i.width
}

func (i zImage) Height() int {
	return i.height
}
