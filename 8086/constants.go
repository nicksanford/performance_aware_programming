package inst

var mov = byte(0b10001000)
var dMask = byte(0b00000010)
var wMask = byte(0b0000001)
var regMoveMask = byte(0b11000000)

// MOD: 00
var memModeNoDisplacmentMap = map[byte]string{
	byte(0b000): "bx + si",
	byte(0b001): "bx + di",
	byte(0b010): "bp + si",
	byte(0b011): "bp + di",
	byte(0b100): "si",
	byte(0b101): "di",
	// DIRECT ADDRESS
	// byte(0b110): ,
	byte(0b111): "bx",
}

// MOD: 01
var memMode8BitDisplacmentMap = map[byte]string{
	byte(0b000): "bx + si",
	byte(0b001): "bx + di",
	byte(0b010): "bp + si",
	byte(0b011): "bp + di",
	byte(0b100): "si",
	byte(0b101): "di",
	byte(0b110): "bp",
	byte(0b111): "bx",
}

// MOD: 10
var memMode16BitDisplacmentMap = map[byte]string{
	byte(0b000): "bx + si",
	byte(0b001): "bx + di",
	byte(0b010): "bp + si",
	byte(0b011): "bp + di",
	byte(0b100): "si",
	byte(0b101): "di",
	byte(0b110): "bp",
	byte(0b111): "bx",
}

// MOD: 11
var regModeNoDisplacementMap = map[byte][2]string{
	byte(0b000): {"al", "ax"},
	byte(0b001): {"cl", "cx"},
	byte(0b010): {"dl", "dx"},
	byte(0b011): {"bl", "bx"},
	byte(0b100): {"ah", "sp"},
	byte(0b101): {"ch", "bp"},
	byte(0b110): {"dh", "si"},
	byte(0b111): {"bh", "di"},
}

var v = map[string]byte{
	"al": byte(0b000),
	"ax": byte(0b000),
	"cl": byte(0b001),
	"cx": byte(0b001),
	"dl": byte(0b010),
	"dx": byte(0b010),
	"bl": byte(0b011),
	"bx": byte(0b011),
	"ah": byte(0b100),
	"sp": byte(0b100),
	"ch": byte(0b101),
	"bp": byte(0b101),
	"dh": byte(0b110),
	"si": byte(0b110),
	"bh": byte(0b111),
	"di": byte(0b111),
}

var wH = map[string]bool{
	"al": false,
	"ax": true,
	"cl": false,
	"cx": true,
	"dl": false,
	"dx": true,
	"bl": false,
	"bx": true,
	"ah": false,
	"sp": true,
	"ch": false,
	"bp": true,
	"dh": false,
	"si": true,
	"bh": false,
	"di": true,
}

type ModType uint8

const (
	ModTypeInvalid = iota
	ModTypeMemoryNoDisplacement
	ModTypeMemory8BitDisplacement
	ModTypeMemory16BitDisplacement
	ModTypeRegToReg
)

type OpType uint8

const (
	OpTypeInvalid = iota
	OpTypeMovRegMemToFromReg
	OpTypeMovImmToRegOrMem
	OpTypeMovImmToReg
	OpTypeMovMemToAcc
	OpTypeMovAccToMem
	OpTypeMovRegOrMemToSegReg
	OpTypeMovSegRegToRegMemory
	OpTypeAddRegMemWithReg
	OpTypeImmToRegOrMem
	OpTypeAddImmToAcc
	OpTypeSubRegMemWithReg
	OpTypeSubImmToAcc
	OpTypeCmpRegMemWithReg
	OpTypeCmpImmToAcc
)
