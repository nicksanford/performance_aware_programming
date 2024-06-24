package inst

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ToLines(data []byte) []string {
	s := string(data)
	body := strings.Split(s, "bits 16")[1]
	return strings.Split(strings.TrimSpace(body), "\n")
}

var regToFullReg = map[string]string{
	"ax": "ax",
	"al": "ax",
	"ah": "ax",
	"bx": "bx",
	"bl": "bx",
	"bh": "bx",
	"cx": "cx",
	"cl": "cx",
	"ch": "cx",
	"dx": "dx",
	"dl": "dx",
	"dh": "dx",
	"sp": "sp",
	"bp": "bp",
	"si": "si",
	"di": "di",
}

var regToRegSize = map[string]RegSize{
	"ax": RegSizeFull,
	"al": RegSizeLow,
	"ah": RegSizeHigh,
	"bx": RegSizeFull,
	"bl": RegSizeLow,
	"bh": RegSizeHigh,
	"cx": RegSizeFull,
	"cl": RegSizeLow,
	"ch": RegSizeHigh,
	"dx": RegSizeFull,
	"dl": RegSizeLow,
	"dh": RegSizeHigh,
	"sp": RegSizeFull,
	"bp": RegSizeFull,
	"si": RegSizeFull,
	"di": RegSizeFull,
}

var regToIndex = map[string]int{
	"ax": 0,
	"bx": 1,
	"cx": 2,
	"dx": 3,
	"sp": 4,
	"bp": 5,
	"si": 6,
	"di": 7,
}

var indexToReg = [8]string{
	"ax",
	"bx",
	"cx",
	"dx",
	"sp",
	"bp",
	"si",
	"di",
}

type CPU struct {
	regs  [8]uint16
	signF bool
	zeroF bool
}

type SimulationResult struct {
	instructions []Inst
	changes      []string
	cpu          CPU
}

type InstType int

const (
	InstTypeUnknown InstType = iota
	InstTypeMov
	InstTypeAdd
	InstTypeSub
	InstTypeCmp
)

type InstSubType int

const (
	InstSubTypeUnknown InstSubType = iota
	InstSubTypeImmToReg
	InstSubTypeRegToReg
	InstSubTypeRegToMem
	InstSubTypeMemToReg
	InstSubTypeImmToDirectMem
	InstSubTypeImmToMem
)

type RegSize int

const (
	RegSizeFull RegSize = iota
	RegSizeLow
	RegSizeHigh
)

type Inst struct {
	instType    InstType
	instSubType InstSubType
	s           string
	bytes       []byte

	targetRegSize RegSize
	targetRegIdx  int

	sourceRegIdx  int
	sourceRegSize RegSize

	targetMem byte

	sourceMem    byte
	displacement uint16

	imm     uint16
	immWord bool

	directMem uint16
}

type Disassembly []Inst

func (d Disassembly) String() string {
	ret := "bits 16\n\n"
	for _, i := range d {
		ret += i.s
	}
	return ret
}

func (it InstType) String() string {
	switch it {
	case InstTypeMov:
		return "mov"
	case InstTypeAdd:
		return "add"
	case InstTypeSub:
		return "sub"
	case InstTypeCmp:
		return "cmp"
	default:
		return fmt.Sprintf("InstTypeUnknown(%d)", it)
	}
}

func (ist InstSubType) String() string {
	switch ist {
	case InstSubTypeImmToReg:
		return "InstSubTypeImmediate"
	case InstSubTypeRegToReg:
		return "InstSubTypeRegToReg"
	case InstSubTypeRegToMem:
		return "InstSubTypeRegToMem"
	case InstSubTypeMemToReg:
		return "InstSubTypeMemToReg"
	default:
		return fmt.Sprintf("InstTypeUnknown(%d)", ist)
	}
}

func PrintSim(data []byte, w io.Writer) error {
	is, err := Dasm(data)
	if err != nil {
		return err
	}

	// tokens, err := Tokenize(lines)
	// if err != nil {
	// 	return err
	// }

	result, err := Simulate(is)
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(result.String())); err != nil {
		return err
	}
	return nil
}

func Tokenize(lines []string) ([]Inst, error) {
	insts := []Inst{}
	for _, l := range lines {
		spaceSep := strings.Split(l, " ")
		out := []string{}
		for _, x := range spaceSep {
			out = append(out, strings.Split(x, ",")[0])
		}
		if len(out) == 3 {
			targetRegister, targetRegOk := regToIndex[out[1]]
			i, err := strconv.ParseInt(out[2], 10, 16)
			sourceRegister, sourceRegOk := regToIndex[out[2]]
			switch {
			case out[0] == "mov" && targetRegOk && err == nil:
				insts = append(insts, Inst{
					instType:     InstTypeMovImmediate,
					s:            l,
					targetRegIdx: targetRegister,
					immWord:      uint16(i)})
			case out[0] == "mov" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					instType:     InstTypeMovRegToReg,
					s:            l,
					targetRegIdx: targetRegister,
					sourceRegIdx: sourceRegister})
			case out[0] == "add" && targetRegOk && err == nil:
				insts = append(insts, Inst{
					instType:     InstTypeAddImmediate,
					s:            l,
					targetRegIdx: targetRegister,
					immWord:      uint16(i)})
			case out[0] == "add" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					instType:     InstTypeAddRegToReg,
					s:            l,
					targetRegIdx: targetRegister,
					sourceRegIdx: sourceRegister})
			case out[0] == "sub" && targetRegOk && err == nil:
				insts = append(insts, Inst{
					instType:     InstTypeSubImmediate,
					s:            l,
					targetRegIdx: targetRegister,
					immWord:      uint16(i)})
			case out[0] == "sub" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					instType:     InstTypeSubRegToReg,
					s:            l,
					targetRegIdx: targetRegister,
					sourceRegIdx: sourceRegister})
			case out[0] == "cmp" && targetRegOk && err == nil:
				insts = append(insts, Inst{
					instType:     InstTypeCmpImmediate,
					s:            l,
					targetRegIdx: targetRegister,
					immWord:      uint16(i)})
			case out[0] == "cmp" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					instType:     InstTypeCmpRegToReg,
					s:            l,
					targetRegIdx: targetRegister,
					sourceRegIdx: sourceRegister})
			default:
				panic("at the disco")
			}
		} else {
			panic("at the club")
		}
	}
	return insts, nil
}

func Simulate(is []Inst) (SimulationResult, error) {
	var cpu CPU
	changes := []string{}
	for _, t := range is {
		switch t.instType {
		case InstTypeMovImmediate:
			before := cpu.regs[t.targetRegIdx]
			cpu.regs[t.targetRegIdx] = t.immWord
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x", indexToReg[t.targetRegIdx], before, cpu.regs[t.targetRegIdx]))
		case InstTypeMovRegToReg:
			before := cpu.regs[t.targetRegIdx]
			cpu.regs[t.targetRegIdx] = cpu.regs[t.sourceRegIdx]
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x", indexToReg[t.targetRegIdx], before, cpu.regs[t.targetRegIdx]))
		case InstTypeAddImmediate:
			before := cpu.regs[t.targetRegIdx]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegIdx] += t.immWord
			cpu.zeroF = cpu.regs[t.targetRegIdx] == 0
			cpu.signF = (cpu.regs[t.targetRegIdx]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegIdx], before, cpu.regs[t.targetRegIdx], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeAddRegToReg:
			before := cpu.regs[t.targetRegIdx]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegIdx] += cpu.regs[t.sourceRegIdx]
			cpu.zeroF = cpu.regs[t.targetRegIdx] == 0
			cpu.signF = (cpu.regs[t.targetRegIdx]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegIdx], before, cpu.regs[t.targetRegIdx], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeSubImmediate:
			before := cpu.regs[t.targetRegIdx]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegIdx] -= t.immWord
			cpu.zeroF = cpu.regs[t.targetRegIdx] == 0
			cpu.signF = (cpu.regs[t.targetRegIdx]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegIdx], before, cpu.regs[t.targetRegIdx], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeSubRegToReg:
			before := cpu.regs[t.targetRegIdx]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegIdx] -= cpu.regs[t.sourceRegIdx]
			cpu.zeroF = cpu.regs[t.targetRegIdx] == 0
			cpu.signF = (cpu.regs[t.targetRegIdx]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegIdx], before, cpu.regs[t.targetRegIdx], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeCmpImmediate:
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegIdx] -= t.immWord
			cpu.zeroF = cpu.regs[t.targetRegIdx] == 0
			cpu.signF = (cpu.regs[t.targetRegIdx]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("flags:%s->%s", cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeCmpRegToReg:
			cpuFlagsBefore := cpu.FlagsString()
			res := cpu.regs[t.targetRegIdx] - cpu.regs[t.sourceRegIdx]
			cpu.zeroF = res == 0
			cpu.signF = (res>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("flags:%s->%s", cpuFlagsBefore, cpu.FlagsString()))
		default:
			panic("at the discotek")
		}
	}
	return SimulationResult{
		instructions: is,
		changes:      changes,
		cpu:          cpu,
	}, nil
}

func (sr *SimulationResult) String() string {
	var ret string
	for i, inst := range sr.instructions {
		ret += fmt.Sprintf("%s ; %s\n", inst.s, sr.changes[i])
	}

	ret += "\nFinal registers:\n"
	for i := range sr.cpu.regs {
		ret += fmt.Sprintf("      %s: 0x%04x (%d)\n", indexToReg[i], sr.cpu.regs[i], sr.cpu.regs[i])
	}

	return ret
}

func (cpu CPU) FlagsString() string {
	switch {
	case cpu.signF && cpu.zeroF:
		return "SZ"
	case cpu.signF:
		return "S"
	case cpu.zeroF:
		return "Z"
	default:
		return ""
	}
}
