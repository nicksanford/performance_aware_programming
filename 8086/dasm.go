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
		} else {
			eac, err := memModeLookup(rm)
			if err != nil {
				return "", 0, err
			}

			dest, source = fmt.Sprintf("[%s]", eac), regS
			i += 2

		}
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
		displacement := int16(data[i+2]) | int16(data[i+3])<<8
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
	var prefix string
	switch modType(b2) {
	case ModTypeMemoryNoDisplacement:
		// check for if R/M == 110 & if so do the 16 bit displacement
		// DIRECT
		s := byteIs(b1, dMask, dMask)
		if rm == 0b110 {
			// panic(fmt.Sprintf("ModTypeMemoryNoDisplacement unimplemented %08b %08b %08b %08b", b1, b2, data[i+2], data[i+3]))
			iInc := 2
			dest = fmt.Sprintf("[%s]", fmt.Sprintf("%d", int16(data[i+iInc])|int16(data[i+iInc+1])<<8))
			iInc += 2
			var imm string
			if w {
				prefix = "word"
				if s {
					imm = fmt.Sprintf("%d", int8(data[i+iInc]))
					iInc += 1
				} else {
					// instruction operates on word (2 byte) data
					imm = fmt.Sprintf("%d", int16(data[i+iInc])|int16(data[i+iInc+1])<<8)
					iInc += 2
				}
			} else {
				prefix = "byte"
				imm = fmt.Sprintf("%d", int8(data[i+iInc]))
				iInc += 1
			}
			i += iInc
			source = imm
		} else {
			eac, err := memModeLookup(rm)
			if err != nil {
				return "", 0, err
			}
			// panic(fmt.Sprintf("%s memmodenodisplacement %08b %08b w: %t, reg: %s, eac: %s", inst, b1, b2, w, regS, eac))
			var imm string
			iInc := 2
			if s {
				prefix = "byte"
				// sign extend 8 bit immediate data to 16 bits
				// instruction operates on word (2 byte) data
				imm = fmt.Sprintf("%d", int8(data[i+iInc]))
				iInc += 1
			} else {
				if w {
					prefix = "word"
					// instruction operates on word (2 byte) data
					imm = fmt.Sprintf("%d", int16(data[i+iInc])|int16(data[i+iInc+1])<<8)
					iInc += 2
				} else {
					prefix = "byte"
					imm = fmt.Sprintf("%d", int8(data[i+iInc]))
					iInc += 1
				}
			}

			source = imm
			dest = fmt.Sprintf("[%s]", eac)
			i += iInc
		}
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
		s := byteIs(b1, dMask, dMask)
		eac, err := memMode16BitDisplacmentLookup(rm)
		if err != nil {
			return "", 0, err
		}

		displacement := int16(data[i+2]) | int16(data[i+3])<<8
		if displacement == 0 {
			dest = fmt.Sprintf("[%s]", eac)
		} else {
			dest = fmt.Sprintf("[%s + %d]", eac, displacement)
		}

		var imm string
		iInc := 4
		if w {
			prefix = "word"
			if s {
				imm = fmt.Sprintf("%d", int8(data[i+iInc]))
				iInc += 1
			} else {
				// instruction operates on word (2 byte) data
				imm = fmt.Sprintf("%d", int16(data[i+iInc])|int16(data[i+iInc+1])<<8)
				iInc += 2
			}
		} else {
			prefix = "byte"
			imm = fmt.Sprintf("%d", int8(data[i+iInc]))
			iInc += 1
		}
		// }
		source = imm

		i += iInc
	case ModTypeRegToReg:
		// TODO: Write mote tests for this case
		s := byteIs(b1, dMask, dMask)
		rmS, err := regLookup(rm, w)
		if err != nil {
			return "", 0, err
		}
		dest = rmS
		var imm string
		iInc := 2
		if s {
			imm = fmt.Sprintf("%d", int8(data[i+iInc]))
			iInc += 1
		} else {
			if w {
				imm = fmt.Sprintf("%d", int16(data[i+iInc])|int16(data[i+iInc+1])<<8)
				iInc += 2
			} else {
				imm = fmt.Sprintf("%d", int8(data[i+iInc]))
				iInc += 1
			}
		}

		i += iInc
		source = imm
	default:
		return "", 0, errors.New("mod field had unexpected value")
	}
	if prefix != "" {
		return fmt.Sprintf("%s %s %s, %s\n", inst, prefix, dest, source), i, nil
	}
	return fmt.Sprintf("%s %s, %s\n", inst, dest, source), i, nil
}

func f(data []byte, i int) (string, int, error) {
	b1 := data[i]
	b2 := data[i+1]
	w := b1&0b00000001 == 0b00000001
	var inst string
	switch toAddSubCmp(b1 << 2 >> 5) {
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
		displacement := int16(data[i+2]) | int16(data[i+3])<<8
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
		fmt.Printf("i: %d, %08b\n", i, data[i])
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
				imm = fmt.Sprintf("%d", int16(data[i+1])|int16(data[i+2])<<8)
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
			w := data[i]&0b00000001 == 0b00000001
			imm := fmt.Sprintf("%d", int8(data[i+1]))
			iInc := 2
			target := "al"
			if w {
				imm = fmt.Sprintf("%d", int16(data[i+1])|int16(data[i+2])<<8)
				iInc += 1
				target = "ax"
			}
			res += fmt.Sprintf("add %s, %s\n", target, imm)
			i += iInc
		case OpTypeSubRegMemWithReg:
			// panic(fmt.Sprintf("%08b %08b", data[i], data[i+1]))
			r, tmpI, err := f(data, i)
			if err != nil {
				return nil, err
			}
			i = tmpI
			res += r
		case OpTypeSubImmToAcc:
			w := data[i]&0b00000001 == 0b00000001
			imm := fmt.Sprintf("%d", int8(data[i+1]))
			iInc := 2
			target := "al"
			if w {
				imm = fmt.Sprintf("%d", int16(data[i+1])|int16(data[i+2])<<8)
				iInc += 1
				target = "ax"
			}
			res += fmt.Sprintf("sub %s, %s\n", target, imm)
			i += iInc
		case OpTypeCmpRegMemWithReg:
			r, tmpI, err := f(data, i)
			if err != nil {
				return nil, err
			}
			i = tmpI
			res += r
		case OpTypeCmpImmToAcc:
			w := data[i]&0b00000001 == 0b00000001
			imm := fmt.Sprintf("%d", int8(data[i+1]))
			iInc := 2
			target := "al"
			if w {
				imm = fmt.Sprintf("%d", int16(data[i+1])|int16(data[i+2])<<8)
				iInc += 1
				target = "ax"
			}
			res += fmt.Sprintf("cmp %s, %s\n", target, imm)
			i += iInc
		case OpTypeJneOrJnz:
			res += fmt.Sprintf("jne %d\n", data[i+1])
			i += 2
		case OpTypeJeOrJz:
			res += fmt.Sprintf("je %d\n", data[i+1])
			i += 2
		case OpTypeJlOrJnge:
			res += fmt.Sprintf("jl %d\n", data[i+1])
			i += 2
		case OpTypeJleOrJng:
			res += fmt.Sprintf("jle %d\n", data[i+1])
			i += 2
		case OpTypeJbeOrJna:
			res += fmt.Sprintf("jbe %d\n", data[i+1])
			i += 2
		case OpTypeJbOrJnae:
			res += fmt.Sprintf("jb %d\n", data[i+1])
			i += 2
		case OpTypeJpOrJpe:
			res += fmt.Sprintf("jp %d\n", data[i+1])
			i += 2
		case OpTypeJo:
			res += fmt.Sprintf("jo %d\n", data[i+1])
			i += 2
		case OpTypeJs:
			res += fmt.Sprintf("js %d\n", data[i+1])
			i += 2
		case OpTypeJnlOrJge:
			res += fmt.Sprintf("jnl %d\n", data[i+1])
			i += 2
		case OpTypeJnleOrJg:
			res += fmt.Sprintf("jnle %d\n", data[i+1])
			i += 2
		case OpTypeJnbOrJae:
			res += fmt.Sprintf("jnb %d\n", data[i+1])
			i += 2
		case OpTypeJnbeOrJa:
			res += fmt.Sprintf("jnbe %d\n", data[i+1])
			i += 2
		case OpTypeJnpOrJpo:
			res += fmt.Sprintf("jnp %d\n", data[i+1])
			i += 2
		case OpTypeJno:
			res += fmt.Sprintf("jno %d\n", data[i+1])
			i += 2
		case OpTypeJns:
			res += fmt.Sprintf("jns %d\n", data[i+1])
			i += 2
		case OpTypeLoop:
			res += fmt.Sprintf("loop %d\n", data[i+1])
			i += 2
		case OpTypeLoopzOrLoope:
			res += fmt.Sprintf("loopz %d\n", data[i+1])
			i += 2
		case OpTypeLoopnzOrLoopne:
			res += fmt.Sprintf("loopnz %d\n", data[i+1])
			i += 2
		case OpTypeJcxz:
			res += fmt.Sprintf("jcxz %d\n", data[i+1])
			i += 2

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
