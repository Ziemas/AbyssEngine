package media

import _ "embed"

//go:embed fonts/diablo_h.otf
var FontDiabloHeavy []byte

//go:embed images/bootlogo.png
var BootLogo []byte

//go:embed images/abyssicon.png
var AbyssIcon []byte

//go:embed shaders/palette.fs
var PaletteFragmentShader string

//go:embed shaders/standard.vs
var StandardVertexShader string
