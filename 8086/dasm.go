package inst

import (
	"errors"
	"fmt"
)

func opTypeMovImmToRegOrMem(data []byte, i int) (string, int, error) {
	b1 := data[i]
	b2 := data[i+1]

	d := byteIs(b1, dMask, dMask)
	w := byteIs(b1, wMask, wMask)

	reg := b2 << 2 >> 5
	rm := b2 << 5 >> 5

	regS, err := regLookup(reg, w)
	if err != nil {
		return "", 0, err
	}
	var source string
	var dest string
	switch modType(b2) {
	case ModTypeMemoryNoDisplacement:
		// check for if R/M == 110 & if so do the 16 bit displacement
		// DIRECT
		if rm == 0b110 {
			i += 4
			panic("ModTypeMemoryNoDisplacement unimplemented")
		}
		eac, err := memModeLookup(rm)
		if err != nil {
			return "", 0, err
		}

		dest, source = fmt.Sprintf("[%s]", eac), regS
		i += 2
	case ModTypeMemory8BitDisplacement:
		eac, err := memMode8BitDisplacmentLookup(rm)
		if err != nil {
			return "", 0, err
		}

		displacement := int8(data[i+2])
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		i += 3

	case ModTypeMemory16BitDisplacement:
		eac, err := memMode16BitDisplacmentLookup(rm)
		if err != nil {
			return "", 0, err
		}
		displacement := int16(data[i+3])<<8 | int16(data[i+2])
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		i += 4
	case ModTypeRegToReg:
		rmS, err := regLookup(rm, w)
		if err != nil {
			return "", 0, err
		}
		dest, source = rmS, regS
		i += 2
	default:
		return "", 0, errors.New("mod field had unexpected value")
	}
	if d {
		dest, source = source, dest
	}
	return fmt.Sprintf("mov %s, %s\n", dest, source), i, nil
}

func opTypeImmToRegOrMem(data []byte, i int) (string, int, error) {
	b1 := data[i]
	b2 := data[i+1]
	// panic(fmt.Sprintf("%b %b", b1, b2))

	var inst string
	switch toAddSubCmp(b2 << 2 >> 5) {
	case Add:
		inst = "add"
	case Sub:
		inst = "sub"
	case Cmp:
		inst = "cmp"
	default:
		panic("invalid AddSubCmp type")
	}

	// s := byteIs(b1, dMask, dMask)
	w := byteIs(b1, wMask, wMask)

	reg := b2 << 2 >> 5
	// panic(fmt.Sprintf("%b", reg))
	rm := b2 << 5 >> 5

	regS, err := regLookup(reg, w)
	if err != nil {
		return "", 0, err
	}
	var source string
	var dest string
	switch modType(b2) {
	case ModTypeMemoryNoDisplacement:
		// check for if R/M == 110 & if so do the 16 bit displacement
		// DIRECT
		if rm == 0b110 {
			i += 4
			panic("ModTypeMemoryNoDisplacement unimplemented")
		}
		eac, err := memModeLookup(rm)
		if err != nil {
			return "", 0, err
		}

		dest, source = fmt.Sprintf("[%s]", eac), regS
		i += 2
	case ModTypeMemory8BitDisplacement:
		eac, err := memMode8BitDisplacmentLookup(rm)
		if err != nil {
			return "", 0, err
		}

		displacement := int8(data[i+2])
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		i += 3

	case ModTypeMemory16BitDisplacement:
		eac, err := memMode16BitDisplacmentLookup(rm)
		if err != nil {
			return "", 0, err
		}
		displacement := int16(data[i+3])<<8 | int16(data[i+2])
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		i += 4
	case ModTypeRegToReg:
		rmS, err := regLookup(rm, w)
		if err != nil {
			return "", 0, err
		}
		dest, source = rmS, regS
		i += 2
		// var iInc int
		var imm string
		if w {
			imm = fmt.Sprintf("%d", int16(data[i+2])<<8|int16(data[i+1]))
			// iInc = 2
		} else {
			imm = fmt.Sprintf("%d", int8(data[i+1]))
			// iInc = 1
		}
		panic(rmS + " " + imm)
	default:
		return "", 0, errors.New("mod field had unexpected value")
	}
	return fmt.Sprintf("%s %s, %s\n", inst, dest, source), i, nil
}

func f(data []byte, i int) (string, int, error) {
	b1 := data[i]
	b2 := data[i+1]
	w := b1&0b00000001 == 0b00000001
	var inst string
	switch toAddSubCmp(data[i] << 2 >> 4) {
	case Add:
		inst = "add"
	case Sub:
		inst = "sub"
	case Cmp:
		inst = "cmp"
	default:
		panic("invalid AddSubCmp type")
	}

	reg := b2 << 2 >> 5
	rm := b2 << 5 >> 5
	d := byteIs(b1, dMask, dMask)

	regS, err := regLookup(reg, w)
	if err != nil {
		return "", 0, err
	}
	var source string
	var dest string
	switch modType(b2) {
	case ModTypeMemoryNoDisplacement:
		// check for if R/M == 110 & if so do the 16 bit displacement
		// DIRECT
		if rm == 0b110 {
			i += 4
			panic("ModTypeMemoryNoDisplacement unimplemented")
		}
		eac, err := memModeLookup(rm)
		if err != nil {
			return "", 0, err
		}

		dest, source = fmt.Sprintf("[%s]", eac), regS
		i += 2
	case ModTypeMemory8BitDisplacement:
		eac, err := memMode8BitDisplacmentLookup(rm)
		if err != nil {
			return "", 0, err
		}

		displacement := int8(data[i+2])
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		i += 3

	case ModTypeMemory16BitDisplacement:
		eac, err := memMode16BitDisplacmentLookup(rm)
		if err != nil {
			return "", 0, err
		}
		displacement := int16(data[i+3])<<8 | int16(data[i+2])
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		i += 4
	case ModTypeRegToReg:
		rmS, err := regLookup(rm, w)
		if err != nil {
			return "", 0, err
		}
		dest, source = rmS, regS
		i += 2
	default:
		return "", 0, errors.New("mod field had unexpected value")
	}
	if d {
		dest, source = source, dest
	}
	return fmt.Sprintf("%s %s, %s\n", inst, dest, source), i, nil
}

func Dasm(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	res := "bits 16\n\n"
	var movT OpType
	for i := 0; i < len(data); {
		movT = opType(data[i])
		switch movT {
		case OpTypeMovMemToAcc:
			panic("MovTypeMemToAcc unimplemented")
		case OpTypeMovAccToMem:
			panic("MovTypeAccToMem unimplemented")
		case OpTypeMovRegOrMemToSegReg:
			panic("MovTypeRegOrMemToSegReg unimplemented")
		case OpTypeMovSegRegToRegMemory:
			panic("MovTypeSegRegToRegMemory unimplemented")
		case OpTypeMovImmToReg:
			w := data[i]<<4>>7 == 0b1
			reg := data[i] << 5 >> 5
			regS, err := regLookup(reg, w)
			if err != nil {
				return nil, err
			}
			var iInc int
			var imm string
			if w {
				imm = fmt.Sprintf("%d", int16(data[i+2])<<8|int16(data[i+1]))
				iInc = 2
			} else {
				imm = fmt.Sprintf("%d", int8(data[i+1]))
				iInc = 1
			}
			res += fmt.Sprintf("mov %s, %s\n", regS, imm)
			i += iInc + 1
		case OpTypeMovImmToRegOrMem:
			panic("OpTypeMovImmToRegOrMem unimplemented")
		case OpTypeMovRegMemToFromReg:
			r, tmpI, err := opTypeMovImmToRegOrMem(data, i)
			if err != nil {
				return nil, err
			}
			i = tmpI
			res += r

		case OpTypeImmToRegOrMem:
			r, tmpI, err := opTypeImmToRegOrMem(data, i)
			if err != nil {
				return nil, err
			}
			i = tmpI
			res += r

			// w := data[i]<<4>>7 == 0b1
			// reg := data[i] << 5 >> 5
			// regS, err := regLookup(reg, w)
			// if err != nil {
			// 	return nil, err
			// }
			// var iInc int
			// var imm string
			// if w {
			// 	imm = fmt.Sprintf("%d", int16(data[i+2])<<8|int16(data[i+1]))
			// 	iInc = 2
			// } else {
			// 	imm = fmt.Sprintf("%d", int8(data[i+1]))
			// 	iInc = 1
			// }
			// var inst string
			// switch toAddSubCmp(data[i+1] << 2 >> 5) {
			// case Add:
			// 	inst = "add"
			// case Sub:
			// 	inst = "sub"
			// case Cmp:
			// 	inst = "cmp"
			// default:
			// 	panic("invalid AddSubCmp type")
			// }
			// res += fmt.Sprintf("%s %s, %s\n", inst, regS, imm)
			// i += iInc + 1
		case OpTypeAddRegMemWithReg:
			r, tmpI, err := f(data, i)
			if err != nil {
				return nil, err
			}
			i = tmpI
			res += r
		case OpTypeAddImmToAcc:
			panic("OpTypeAddImmToAcc unimplemented")
		case OpTypeSubRegMemWithReg:
			panic("OpTypeSubRegMemWithReg unimplemented")
		case OpTypeSubImmToAcc:
			panic("OpTypeSubImmToAcc unimplemented")
		case OpTypeCmpRegMemWithReg:
			panic("OpTypeCmpRegMemWithReg unimplemented")
		case OpTypeCmpImmToAcc:
			panic("OpTypeCmpImmToAcc unimplemented")
		default:
			return nil, fmt.Errorf("unexpected opcode %d", movT)
		}
		println(res)
	}
	return []byte(res), nil
}

type AddSubCmp uint8

const (
	Unknown AddSubCmp = iota
	Add
	Sub
	Cmp
)

func toAddSubCmp(b byte) AddSubCmp {
	switch b {
	case 0b000:
		return Add
	case 0b101:
		return Sub
	case 0b111:
		return Cmp
	default:
		return Unknown
	}
}
