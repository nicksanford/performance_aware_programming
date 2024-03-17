package inst

import "errors"

func wHi(reg string) (bool, error) {
	tf, ok := wH[reg]
	if !ok {
		return false, errors.New("unvalid register name")
	}
	return tf, nil
}

func regToBit(reg string) (byte, error) {
	b, ok := v[reg]
	if !ok {
		return 0, errors.New("invalid register name")
	}
	return b, nil
}

func byteIs(b byte, mask byte, test byte) bool {
	return b&mask == test
}

// MOD:11
func regLookup(reg byte, w bool) (string, error) {
	ss, ok := regModeNoDisplacementMap[reg]
	if !ok {
		return "", errors.New("invlid byte")
	}

	l := 0
	if w {
		l = 1
	}
	return ss[l], nil
}

// MOD:00
func memModeLookup(rm byte) (string, error) {
	ss, ok := memModeNoDisplacmentMap[rm]
	if !ok {
		return "", errors.New("invlid byte")
	}

	return ss, nil
}

// MOD:01
func memMode8BitDisplacmentLookup(rm byte) (string, error) {
	ss, ok := memMode8BitDisplacmentMap[rm]
	if !ok {
		return "", errors.New("invlid byte")
	}

	return ss, nil
}

// MOD:10
func memMode16BitDisplacmentLookup(rm byte) (string, error) {
	ss, ok := memMode16BitDisplacmentMap[rm]
	if !ok {
		return "", errors.New("invlid byte")
	}

	return ss, nil
}

func modType(b byte) ModType {
	masked := b & 0b11000000
	switch masked {
	case 0b00000000:
		return ModTypeMemoryNoDisplacement
	case 0b01000000:
		return ModTypeMemory8BitDisplacement
	case 0b10000000:
		return ModTypeMemory16BitDisplacement
	case 0b11000000:
		return ModTypeRegToReg
	default:
		return ModTypeInvalid

	}
}

func movType(b byte) MovType {
	if b&0b11110000 == 0b10110000 {
		return MovTypeImmToReg
	}

	if b&0b11111100 == 0b10001000 {
		return MovTypeRegMemToFromReg
	}

	switch b & 0b11111110 {
	case 0b11000110:
		return MovTypeImmToRegOrMem
	case 0b10100000:
		return MovTypeMemToAcc
	case 0b10100010:
		return MovTypeAccToMem
	}

	switch b {
	case 0b10001110:
		return MovTypeRegOrMemToSegReg
	case 0b10001100:
		return MovTypeSegRegToRegMemory
	}

	return MovTypeInvalid
}
