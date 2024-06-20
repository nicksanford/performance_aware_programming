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

type Tokens struct {
	instructions []Inst
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
	InstTypeMovImmediate
	InstTypeMovRegToReg
	InstTypeAddImmediate
	InstTypeAddRegToReg
	InstTypeSubImmediate
	InstTypeSubRegToReg
	InstTypeCmpImmediate
	InstTypeCmpRegToReg
)

type Inst struct {
	t               InstType
	s               string
	targetRegister  int
	sourceRegister  int
	sourceImmediate uint16
}

func (it InstType) String() string {
	switch it {
	case InstTypeMovImmediate:
		return "InstTypeMov"
	default:
		return fmt.Sprintf("InstTypeUnknown(%d)", it)
	}
}

func PrintSim(data []byte, w io.Writer) error {
	data, err := Dasm(data)
	if err != nil {
		return err
	}

	lines := ToLines(data)

	tokens, err := Tokenize(lines)
	if err != nil {
		return err
	}

	result, err := Simulate(tokens)
	if err != nil {
		return err
	}

	if _, err := w.Write([]byte(result.String())); err != nil {
		return err
	}
	return nil
}

// TODO: Implement
func Tokenize(lines []string) (Tokens, error) {
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
					t:               InstTypeMovImmediate,
					s:               l,
					targetRegister:  targetRegister,
					sourceImmediate: uint16(i)})
			case out[0] == "mov" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					t:              InstTypeMovRegToReg,
					s:              l,
					targetRegister: targetRegister,
					sourceRegister: sourceRegister})
			case out[0] == "add" && targetRegOk && err == nil:
				insts = append(insts, Inst{
					t:               InstTypeAddImmediate,
					s:               l,
					targetRegister:  targetRegister,
					sourceImmediate: uint16(i)})
			case out[0] == "add" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					t:              InstTypeAddRegToReg,
					s:              l,
					targetRegister: targetRegister,
					sourceRegister: sourceRegister})
			case out[0] == "sub" && targetRegOk && err == nil:
				insts = append(insts, Inst{
					t:               InstTypeSubImmediate,
					s:               l,
					targetRegister:  targetRegister,
					sourceImmediate: uint16(i)})
			case out[0] == "sub" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					t:              InstTypeSubRegToReg,
					s:              l,
					targetRegister: targetRegister,
					sourceRegister: sourceRegister})
			case out[0] == "cmp" && targetRegOk && err == nil:
				insts = append(insts, Inst{
					t:               InstTypeCmpImmediate,
					s:               l,
					targetRegister:  targetRegister,
					sourceImmediate: uint16(i)})
			case out[0] == "cmp" && targetRegOk && sourceRegOk:
				insts = append(insts, Inst{
					t:              InstTypeCmpRegToReg,
					s:              l,
					targetRegister: targetRegister,
					sourceRegister: sourceRegister})
			default:
				panic("at the disco")
			}
		} else {
			panic("at the club")
		}
	}
	return Tokens{instructions: insts}, nil
}

// TODO: Implement
func Simulate(ts Tokens) (SimulationResult, error) {
	var cpu CPU
	changes := []string{}
	for _, t := range ts.instructions {
		switch t.t {
		case InstTypeMovImmediate:
			before := cpu.regs[t.targetRegister]
			cpu.regs[t.targetRegister] = t.sourceImmediate
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x", indexToReg[t.targetRegister], before, cpu.regs[t.targetRegister]))
		case InstTypeMovRegToReg:
			before := cpu.regs[t.targetRegister]
			cpu.regs[t.targetRegister] = cpu.regs[t.sourceRegister]
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x", indexToReg[t.targetRegister], before, cpu.regs[t.targetRegister]))
		case InstTypeAddImmediate:
			before := cpu.regs[t.targetRegister]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegister] += t.sourceImmediate
			cpu.zeroF = cpu.regs[t.targetRegister] == 0
			cpu.signF = (cpu.regs[t.targetRegister]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegister], before, cpu.regs[t.targetRegister], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeAddRegToReg:
			before := cpu.regs[t.targetRegister]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegister] += cpu.regs[t.sourceRegister]
			cpu.zeroF = cpu.regs[t.targetRegister] == 0
			cpu.signF = (cpu.regs[t.targetRegister]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegister], before, cpu.regs[t.targetRegister], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeSubImmediate:
			before := cpu.regs[t.targetRegister]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegister] -= t.sourceImmediate
			cpu.zeroF = cpu.regs[t.targetRegister] == 0
			cpu.signF = (cpu.regs[t.targetRegister]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegister], before, cpu.regs[t.targetRegister], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeSubRegToReg:
			before := cpu.regs[t.targetRegister]
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegister] -= cpu.regs[t.sourceRegister]
			cpu.zeroF = cpu.regs[t.targetRegister] == 0
			cpu.signF = (cpu.regs[t.targetRegister]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("%s:%#x->%#x flags:%s->%s", indexToReg[t.targetRegister], before, cpu.regs[t.targetRegister], cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeCmpImmediate:
			cpuFlagsBefore := cpu.FlagsString()
			cpu.regs[t.targetRegister] -= t.sourceImmediate
			cpu.zeroF = cpu.regs[t.targetRegister] == 0
			cpu.signF = (cpu.regs[t.targetRegister]>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("flags:%s->%s", cpuFlagsBefore, cpu.FlagsString()))
		case InstTypeCmpRegToReg:
			cpuFlagsBefore := cpu.FlagsString()
			res := cpu.regs[t.targetRegister] - cpu.regs[t.sourceRegister]
			cpu.zeroF = res == 0
			cpu.signF = (res>>15)&0b1 == 0b1
			changes = append(changes, fmt.Sprintf("flags:%s->%s", cpuFlagsBefore, cpu.FlagsString()))
		default:
			panic("at the discotek")
		}
	}
	return SimulationResult{
		instructions: ts.instructions,
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
