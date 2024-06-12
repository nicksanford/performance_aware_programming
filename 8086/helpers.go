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

func opType(b byte) OpType {
	if b&0b11110000 == 0b10110000 {
		return OpTypeMovImmToReg
	}

	if b&0b11111100 == 0b10001000 {
		return OpTypeMovRegMemToFromReg
	}

	switch b & 0b11111110 {
	case 0b11000110:
		return OpTypeMovImmToRegOrMem
	case 0b10100000:
		return OpTypeMovMemToAcc
	case 0b10100010:
		return OpTypeMovAccToMem
	}

	switch b {
	case 0b10001110:
		return OpTypeMovRegOrMemToSegReg
	case 0b10001100:
		return OpTypeMovSegRegToRegMemory
	}

	if b&0b11111100 == 0 {
		return OpTypeAddRegMemWithReg
	}

	if b&0b11111100 == 0b10000000 {
		return OpTypeImmToRegOrMem
	}

	if b&0b11111110 == 0b00000100 {
		return OpTypeAddImmToAcc
	}

	if b&0b11111100 == 0b00101000 {
		return OpTypeSubRegMemWithReg
	}

	if b&0b11111110 == 0b00101100 {
		return OpTypeSubImmToAcc
	}

	if b&0b11111100 == 0b00111000 {
		return OpTypeCmpRegMemWithReg
	}

	if b&0b11111110 == 0b00111100 {
		return OpTypeCmpImmToAcc
	}

	switch b {
	case 0b01110101:
		return OpTypeJneOrJnz
	case 0b01110100:
		return OpTypeJeOrJz
	case 0b01111100:
		return OpTypeJlOrJnge
	case 0b01111110:
		return OpTypeJleOrJng
	case 0b01110110:
		return OpTypeJbeOrJna
	case 0b01110010:
		return OpTypeJbOrJnae
	case 0b01111010:
		return OpTypeJpOrJpe
	case 0b01110000:
		return OpTypeJo
	case 0b01111000:
		return OpTypeJs
	case 0b01111101:
		return OpTypeJnlOrJge
	case 0b01111111:
		return OpTypeJnleOrJg
	case 0b01110011:
		return OpTypeJnbOrJae
	case 0b01110111:
		return OpTypeJnbeOrJa
	case 0b01111011:
		return OpTypeJnpOrJpo
	case 0b01110001:
		return OpTypeJno
	case 0b01111001:
		return OpTypeJns
	case 0b11100010:
		return OpTypeLoop
	case 0b11100001:
		return OpTypeLoopzOrLoope
	case 0b11100000:
		return OpTypeLoopnzOrLoopne
	case 0b11100011:
		return OpTypeJcxz
	}

	return OpTypeInvalid
}
