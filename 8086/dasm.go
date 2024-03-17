package inst

import (
	"errors"
	"fmt"
)

func Dasm(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	res := "bits 16\n\n"
	var movT MovType
	for i := 0; i < len(data); {
		movT = movType(data[i])
		switch movT {
		case MovTypeMemToAcc:
			panic("MovTypeMemToAcc unimplemented")
		case MovTypeAccToMem:
			panic("MovTypeAccToMem unimplemented")
		case MovTypeRegOrMemToSegReg:
			panic("MovTypeRegOrMemToSegReg unimplemented")
		case MovTypeSegRegToRegMemory:
			panic("MovTypeSegRegToRegMemory unimplemented")
		case MovTypeImmToReg:
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
		case MovTypeImmToRegOrMem:
			panic("MovTypeImmToRegOrMem unimplemented")
		case MovTypeRegMemToFromReg:
			b1 := data[i]
			b2 := data[i+1]

			d := byteIs(b1, dMask, dMask)
			w := byteIs(b1, wMask, wMask)

			reg := b2 << 2 >> 5
			rm := b2 << 5 >> 5

			regS, err := regLookup(reg, w)
			if err != nil {
				return nil, err
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
					return nil, err
				}
				dest, source = eac, regS
				if d {
					dest, source = regS, eac
				}
				i += 2
			case ModTypeMemory8BitDisplacement:
				eacFmt, err := memMode8BitDisplacmentLookup(rm)
				if err != nil {
					return nil, err
				}
				// TODO: Change the memMode8BitDisplacmentLookup function to return just the register & have addition templating be the responsibility of this function so that the 0 case is handled properly
				eac := fmt.Sprintf(eacFmt, int8(data[i+2]))
				dest, source = eac, regS
				if d {
					dest, source = regS, eac
				}
				i += 3

			case ModTypeMemory16BitDisplacement:
				eacFmt, err := memMode16BitDisplacmentLookup(rm)
				if err != nil {
					return nil, err
				}
				eac := fmt.Sprintf(eacFmt, int16(data[i+3])<<8|int16(data[i+2]))
				dest, source = eac, regS
				if d {
					dest, source = regS, eac
				}
				i += 4
			case ModTypeRegToReg:
				rmS, err := regLookup(rm, w)
				if err != nil {
					return nil, err
				}
				dest, source = rmS, regS
				if d {
					dest, source = regS, rmS
				}
				i += 2
			default:
				return nil, errors.New("mod field had unexpected value")
			}
			res += fmt.Sprintf("mov %s, %s\n", dest, source)
		}
	}
	return []byte(res), nil

	// for i := 0; i < len(data); i += 2 {
	// 	b1 := data[i]
	// 	b2 := data[i+1]

	// 	if !byteIs(b1, movMask, mov) {
	// 		return nil, errors.New("only mov commands supported")
	// 	}

	// }
}
