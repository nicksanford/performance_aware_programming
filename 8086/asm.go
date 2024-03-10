package main

import (
	"errors"
	"strings"
)

func asm(data []byte) ([]byte, error) {
	rawS := string(data)
	if !strings.Contains(rawS, "bits 16") {
		return nil, errors.New("only 16 bit 8086 instructions")
	}
	instructions := strings.Split(strings.Split(rawS, "bits 16")[1], "\n")
	out := []byte{}
	for _, is := range instructions {
		i := strings.TrimSpace(is)
		if i == "" {
			continue
		}
		if !strings.HasPrefix(i, "mov") {
			return nil, errors.New("only mov instructions supported")
		}
		rawRegs := strings.Split(strings.TrimPrefix(i, "mov"), ",")
		regs := []string{}
		for _, reg := range rawRegs {
			regs = append(regs, strings.ToLower(strings.TrimSpace(reg)))
		}
		if len(regs) != 2 {
			return nil, errors.New("malformed assembly")
		}
		reg1, err := regToBit(regs[0])
		if err != nil {
			return nil, err
		}
		reg2, err := regToBit(regs[1])
		if err != nil {
			return nil, err
		}
		w, err := wHi(regs[0])
		if err != nil {
			return nil, err
		}
		w2, err := wHi(regs[1])
		if err != nil {
			return nil, err
		}
		if w != w2 {
			return nil, errors.New("register 1 & 2 disagree on whether w bit is hi")
		}
		var wBit byte
		if w {
			wBit = wMask
		}
		// Never set d bit, no reason to
		var dBit byte

		out = append(out,
			[]byte{mov | dBit | wBit,
				regMoveMask | reg2<<3 | reg1,
			}...)
	}
	return out, nil
}
