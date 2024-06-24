package inst

import (
	"errors"
	"fmt"
)

func opTypeMovImmToRegOrMem(data []byte, i int) (Inst, error) {
	var iInc int
	b1 := data[i]
	b2 := data[i+1]

	d := byteIs(b1, dMask, dMask)
	w := byteIs(b1, wMask, wMask)

	reg := b2 << 2 >> 5
	rm := b2 << 5 >> 5

	regS, err := regLookup(reg, w)
	if err != nil {
		return Inst{}, err
	}
	var instSubType InstSubType
	var source, dest string
	var targetMem, sourceMem byte
	var targetReg, sourceReg int
	var targetRegSize, sourceRegSize RegSize
	switch modType(b2) {
	case ModTypeMemoryNoDisplacement:
		// check for if R/M == 110 & if so do the 16 bit displacement
		// DIRECT
		if rm == 0b110 {
			panic("ModTypeMemoryNoDisplacement unimplemented")
		} else {
			instSubType = InstSubTypeRegToMem
			eac, err := memModeLookup(rm)
			if err != nil {
				return Inst{}, err
			}

			dest, source = fmt.Sprintf("[%s]", eac), regS
			targetMem = rm
			sourceReg = regToIndex[regToFullReg[regS]]
			sourceRegSize = regToRegSize[regS]
			iInc = 2
		}
	case ModTypeMemory8BitDisplacement:
		instSubType = InstSubTypeRegToMem
		eac, err := memModeLookup(rm)
		if err != nil {
			return Inst{}, err
		}

		displacement := int8(data[i+2])
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		targetMem = rm
		sourceReg = regToIndex[regToFullReg[regS]]
		sourceRegSize = regToRegSize[regS]
		iInc = 3

	case ModTypeMemory16BitDisplacement:
		instSubType = InstSubTypeRegToMem
		eac, err := memModeLookup(rm)
		if err != nil {
			return Inst{}, err
		}
		displacement := int16(data[i+2]) | int16(data[i+3])<<8
		if displacement == 0 {
			dest, source = fmt.Sprintf("[%s]", eac), regS
		} else {
			dest, source = fmt.Sprintf("[%s + %d]", eac, displacement), regS
		}
		targetMem = rm
		sourceReg = regToIndex[regToFullReg[regS]]
		sourceRegSize = regToRegSize[regS]
		iInc = 4
	case ModTypeReg:
		instSubType = InstSubTypeRegToReg
		rmS, err := regLookup(rm, w)
		if err != nil {
			return Inst{}, err
		}
		dest, source = rmS, regS
		sourceReg = regToIndex[regToFullReg[regS]]
		sourceRegSize = regToRegSize[regS]
		targetReg = regToIndex[regToFullReg[rmS]]
		targetRegSize = regToRegSize[rmS]
		iInc = 2
	default:
		return Inst{}, errors.New("mod field had unexpected value")
	}

	if d {
		if instSubType == InstSubTypeRegToMem {
			instSubType = InstSubTypeMemToReg
			sourceMem = targetMem
			targetMem = 0

			targetReg = sourceReg
			sourceReg = 0

			targetRegSize = sourceRegSize
			sourceRegSize = RegSizeFull
		}
		dest, source = source, dest
	}
	return Inst{
		instType:      InstTypeMov,
		instSubType:   instSubType,
		s:             fmt.Sprintf("mov %s, %s", dest, source),
		bytes:         data[i : i+iInc],
		targetMem:     targetMem,
		sourceMem:     sourceMem,
		targetRegIdx:  targetReg,
		sourceRegIdx:  sourceReg,
		targetRegSize: targetRegSize,
		sourceRegSize: sourceRegSize,
	}, nil
}

func opTypeImmToRegOrMem(data []byte, i int) (Inst, error) {
	b1 := data[i]
	b2 := data[i+1]

	var instType InstType
	switch toAddSubCmp(b2 << 2 >> 5) {
	case AddSubCmpAdd:
		instType = InstTypeAdd
	case AddSubCmpSub:
		instType = InstTypeSub
	case AddSubCmpCmp:
		instType = InstTypeCmp
	default:
		panic("invalid AddSubCmp type")
	}

	w := byteIs(b1, wMask, wMask)

	rm := b2 << 5 >> 5

	switch modType(b2) {
	case ModTypeMemoryNoDisplacement:
		// check for if R/M == 110 & if so do the 16 bit displacement
		// DIRECT
		s := byteIs(b1, dMask, dMask)
		if rm == 0b110 {
			iInc := 2
			directMem := uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
			iInc += 2
			var imm uint16
			wordOrByte := "byte"
			if w {
				wordOrByte = "word"
				if s {
					// instruction operates on word (2 byte) data
					imm = uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
					iInc += 2
				} else {
					imm = uint16(data[i+iInc])
					iInc += 1
				}
			} else {
				imm = uint16(data[i+iInc])
				iInc += 1
			}
			return Inst{
				instType:    instType,
				instSubType: InstSubTypeImmToDirectMem,
				s: fmt.Sprintf(
					"%s %s %s, %s",
					instType.String(),
					wordOrByte,
					fmt.Sprintf("[%s]", fmt.Sprintf("%d", directMem)),
					fmt.Sprintf("%d", imm),
				),
				bytes:     data[i : i+iInc],
				imm:       imm,
				immWord:   wordOrByte == "word",
				directMem: directMem,
			}, nil
		} else {
			eac, err := memModeLookup(rm)
			if err != nil {
				return Inst{}, err
			}
			// var imm string
			// iInc := 2
			// if s {
			// 	prefix = "byte"
			// 	// sign extend 8 bit immediate data to 16 bits
			// 	// instruction operates on word (2 byte) data
			// 	imm = fmt.Sprintf("%d", int8(data[i+iInc]))
			// 	iInc += 1
			// } else {
			// 	if w {
			// 		prefix = "word"
			// 		// instruction operates on word (2 byte) data
			// 		imm = fmt.Sprintf("%d", uint16(data[i+iInc])|uint16(data[i+iInc+1])<<8)
			// 		iInc += 2
			// 	} else {
			// 		prefix = "byte"
			// 		imm = fmt.Sprintf("%d", int8(data[i+iInc]))
			// 		iInc += 1
			// 	}
			// }
			var imm uint16
			wordOrByte := "byte"
			iInc := 2
			if w {
				wordOrByte = "word"
				if s {
					// instruction operates on word (2 byte) data
					imm = uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
					iInc += 2
				} else {
					imm = uint16(data[i+iInc])
					iInc += 1
				}
			} else {
				imm = uint16(data[i+iInc])
				iInc += 1
			}

			return Inst{
				instType:    instType,
				instSubType: InstSubTypeImmToMem,
				s: fmt.Sprintf(
					"%s %s %s, %s",
					instType.String(),
					wordOrByte,
					fmt.Sprintf("[%s]", eac),
					fmt.Sprintf("%d", imm),
				),
				bytes:     data[i : i+iInc],
				targetMem: rm,
				imm:       imm,
				immWord:   wordOrByte == "word",
			}, nil
		}
	case ModTypeMemory8BitDisplacement:
		s := byteIs(b1, dMask, dMask)
		eac, err := memModeLookup(rm)
		if err != nil {
			return Inst{}, err
		}

		displacement := uint16(data[i+2])
		var imm uint16
		wordOrByte := "byte"
		iInc := 3
		if w {
			wordOrByte = "word"
			if s {
				// instruction operates on word (2 byte) data
				imm = uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
				iInc += 2
			} else {
				imm = uint16(data[i+iInc])
				iInc += 1
			}
		} else {
			imm = uint16(data[i+iInc])
			iInc += 1
		}
		return Inst{
			instType:    instType,
			instSubType: InstSubTypeImmToMem,
			s: fmt.Sprintf(
				"%s %s %s, %s",
				instType.String(),
				wordOrByte,
				fmt.Sprintf("[%s + %d]", eac, displacement),
				fmt.Sprintf("%d", imm),
			),
			bytes:        data[i : i+iInc],
			targetMem:    rm,
			imm:          imm,
			immWord:      wordOrByte == "word",
			displacement: displacement,
		}, nil

	case ModTypeMemory16BitDisplacement:
		s := byteIs(b1, dMask, dMask)
		eac, err := memModeLookup(rm)
		if err != nil {
			return Inst{}, err
		}

		displacement := uint16(data[i+2]) | uint16(data[i+3])<<8

		var imm uint16
		wordOrByte := "byte"
		iInc := 4
		if w {
			wordOrByte = "word"
			if s {
				// instruction operates on word (2 byte) data
				imm = uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
				iInc += 2
			} else {
				imm = uint16(data[i+iInc])
				iInc += 1
			}
		} else {
			imm = uint16(data[i+iInc])
			iInc += 1
		}
		return Inst{
			instType:    instType,
			instSubType: InstSubTypeImmToMem,
			s: fmt.Sprintf(
				"%s %s %s, %s",
				instType.String(),
				wordOrByte,
				fmt.Sprintf("[%s + %d]", eac, displacement),
				fmt.Sprintf("%d", imm),
			),
			bytes:        data[i : i+iInc],
			targetMem:    rm,
			imm:          imm,
			immWord:      wordOrByte == "word",
			displacement: displacement,
		}, nil
	case ModTypeReg:
		// TODO: Write mote tests for this case
		s := byteIs(b1, dMask, dMask)
		rmS, err := regLookup(rm, w)
		if err != nil {
			return Inst{}, err
		}
		var imm uint16
		wordOrByte := "byte"
		iInc := 2
		if w {
			wordOrByte = "word"
			if s {
				// instruction operates on word (2 byte) data
				imm = uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
				iInc += 2
			} else {
				imm = uint16(data[i+iInc])
				iInc += 1
			}
		} else {
			imm = uint16(data[i+iInc])
			iInc += 1
		}
		return Inst{
			instType:    instType,
			instSubType: InstSubTypeImmToReg,
			s: fmt.Sprintf(
				"%s %s %s, %s",
				instType.String(),
				wordOrByte,
				rmS,
				fmt.Sprintf("%d", imm),
			),
			bytes:         data[i : i+iInc],
			targetRegIdx:  regToIndex[regToFullReg[rmS]],
			targetRegSize: regToRegSize[rmS],
			imm:           imm,
			immWord:       wordOrByte == "word",
		}, nil

	default:
		return Inst{}, errors.New("mod field had unexpected value")
	}
}

func opRegMemWithReg(data []byte, i int) (Inst, error) {
	b1 := data[i]
	b2 := data[i+1]
	w := b1&0b00000001 == 0b00000001
	var instType InstType
	switch toAddSubCmp(b1 << 2 >> 5) {
	case AddSubCmpAdd:
		instType = InstTypeAdd
	case AddSubCmpSub:
		instType = InstTypeSub
	case AddSubCmpCmp:
		instType = InstTypeCmp
	default:
		panic("invalid AddSubCmp type")
	}

	reg := b2 << 2 >> 5
	rm := b2 << 5 >> 5
	d := byteIs(b1, dMask, dMask)

	regS, err := regLookup(reg, w)
	if err != nil {
		return Inst{}, err
	}
	sourceRegIdx := regToIndex[regToFullReg[regS]]
	sourceRegSize := regToRegSize[regS]
	iInc := 2
	var instSubType InstSubType
	var displacement uint16
	var targetMem, sourceMem byte
	var targetRegSize RegSize
	var targetRegIdx int
	var source, dest string
	switch modType(b2) {
	case ModTypeMemoryNoDisplacement:
		// check for if R/M == 110 & if so do the 16 bit displacement
		// DIRECT
		if rm == 0b110 {
			panic("ModTypeMemoryNoDisplacement unimplemented")
		}
		instSubType = InstSubTypeRegToMem
		eac, err := memModeLookup(rm)
		if err != nil {
			return Inst{}, err
		}
		targetMem = rm

		dest, source = fmt.Sprintf("[%s]", eac), regS

	case ModTypeMemory8BitDisplacement:
		instSubType = InstSubTypeRegToMem
		eac, err := memModeLookup(rm)
		if err != nil {
			return Inst{}, err
		}

		targetMem = rm
		displacement = uint16(data[i+iInc])
		dest, source = fmt.Sprint("[%s %d]", eac, displacement), regS
		iInc += 1

	case ModTypeMemory16BitDisplacement:
		instSubType = InstSubTypeRegToMem
		eac, err := memModeLookup(rm)
		if err != nil {
			return Inst{}, err
		}
		targetMem = rm
		displacement = uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
		dest, source = fmt.Sprint("[%s %d]", eac, displacement), regS
		iInc += 2
	case ModTypeReg:
		instSubType = InstSubTypeRegToReg
		rmS, err := regLookup(rm, w)
		if err != nil {
			return Inst{}, err
		}

		targetRegIdx = regToIndex[regToFullReg[rmS]]
		targetRegSize = regToRegSize[rmS]
		dest, source = rmS, regS
	default:
		return Inst{}, errors.New("mod field had unexpected value")
	}

	if d && instSubType == InstSubTypeRegToReg {
		sourceRegIdx, sourceRegSize = targetRegIdx, targetRegSize
	}

	if d && instSubType == InstSubTypeRegToMem {
		instSubType = InstSubTypeMemToReg
		sourceMem = targetMem
		targetMem = 0

		targetRegIdx = sourceRegIdx
		targetRegSize = sourceRegSize
		sourceRegIdx, sourceRegSize = 0, RegSizeFull
	}

	if d {
		dest, source = source, dest
	}
	return Inst{
		instType:    instType,
		instSubType: instSubType,
		s: fmt.Sprintf("%s %s, %s",
			instType.String(),
			dest,
			source),
		bytes:         data[i : i+iInc],
		displacement:  displacement,
		targetMem:     targetMem,
		sourceMem:     sourceMem,
		targetRegIdx:  targetRegIdx,
		targetRegSize: targetRegSize,
		sourceRegIdx:  sourceRegIdx,
		sourceRegSize: sourceRegSize,
	}, nil
}

func Dasm(data []byte) (Disassembly, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}
	res := []Inst{}
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
			var imm uint16
			var immWord bool
			iInc := 1
			if w {
				immWord = true
				imm = uint16(data[i+iInc]) | uint16(data[i+iInc+1])<<8
				iInc += 2
			} else {
				imm = uint16(data[i+iInc])
				iInc += 1
			}
			res = append(res, Inst{
				instType:     InstTypeMov,
				instSubType:  InstSubTypeImmToReg,
				s:            fmt.Sprintf("mov %s, %d", regS, imm),
				bytes:        data[i : i+iInc],
				targetRegIdx: regToIndex[regS],
				imm:          imm,
				immWord:      immWord,
			})
			i += iInc
		case OpTypeMovImmToRegOrMem:
			panic("OpTypeMovImmToRegOrMem unimplemented")
		case OpTypeMovRegMemToFromReg:
			inst, err := opTypeMovImmToRegOrMem(data, i)
			if err != nil {
				return nil, err
			}
			res = append(res, inst)
			i = len(inst.bytes)

		case OpTypeImmToRegOrMem:
			inst, err := opTypeImmToRegOrMem(data, i)
			if err != nil {
				return nil, err
			}
			res = append(res, inst)
			i = len(inst.bytes)

		case OpTypeAddRegMemWithReg,
			OpTypeSubRegMemWithReg,
			OpTypeCmpRegMemWithReg:
			inst, err := opRegMemWithReg(data, i)
			if err != nil {
				return nil, err
			}
			res = append(res, inst)
			i = len(inst.bytes)
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
	}
	return []byte(res), nil
}

type AddSubCmp uint8

const (
	AddSubCmpUnknown AddSubCmp = iota
	AddSubCmpAdd
	AddSubCmpSub
	AddSubCmpCmp
)

func toAddSubCmp(b byte) AddSubCmp {
	switch b {
	case 0b000:
		return AddSubCmpAdd
	case 0b101:
		return AddSubCmpSub
	case 0b111:
		return AddSubCmpCmp
	default:
		return AddSubCmpUnknown
	}
}
