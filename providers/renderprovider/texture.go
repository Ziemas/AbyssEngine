package renderprovider

type Texture interface {
	Width() int
	Height() int
	ID() int
}
