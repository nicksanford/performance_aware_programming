package main

var movMask = byte(0b11111100)
var mov = byte(0b10001000)
var dMask = byte(0b00000010)
var wMask = byte(0b00000001)
var regMoveMask = byte(0b11000000)
var m = map[byte][2]string{
	byte(0b0000_0000): {"al", "ax"},
	byte(0b0000_0001): {"cl", "cx"},
	byte(0b0000_0010): {"dl", "dx"},
	byte(0b0000_0011): {"bl", "bx"},
	byte(0b0000_0100): {"ah", "sp"},
	byte(0b0000_0101): {"ch", "bp"},
	byte(0b0000_0110): {"dh", "si"},
	byte(0b0000_0111): {"bh", "di"},
}
var v = map[string]byte{
	"al": byte(0b0000_0000),
	"ax": byte(0b0000_0000),
	"cl": byte(0b0000_0001),
	"cx": byte(0b0000_0001),
	"dl": byte(0b0000_0010),
	"dx": byte(0b0000_0010),
	"bl": byte(0b0000_0011),
	"bx": byte(0b0000_0011),
	"ah": byte(0b0000_0100),
	"sp": byte(0b0000_0100),
	"ch": byte(0b0000_0101),
	"bp": byte(0b0000_0101),
	"dh": byte(0b0000_0110),
	"si": byte(0b0000_0110),
	"bh": byte(0b0000_0111),
	"di": byte(0b0000_0111),
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
