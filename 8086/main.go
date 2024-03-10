package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

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

func main() {
	l := log.New(os.Stderr, "", 0)
	help := flag.Bool("h", false, "help")
	assemble := flag.Bool("a", false, "assemble, defaults to false i.e. disassemble")
	flag.Parse()
	if *help {
		l.Printf("usage: %s <file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(flag.Args()) != 1 {
		l.Fatalf("usage: %s <file>\n", os.Args[0])
	}

	data, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		l.Fatalf(err.Error())
	}

	if *assemble {
		data, err = asm(data)
		if err != nil {
			l.Fatal(err.Error())
		}
	} else {
		data, err = dasm(data)
		if err != nil {
			l.Fatal(err.Error())
		}
	}

	if _, err := os.Stdout.Write(data); err != nil {
		l.Fatal(err.Error())
	}
}

func byteIs(b byte, mask byte, test byte) bool {
	return b&mask == test
}

func regLookup(reg byte, w bool) (string, error) {
	ss, ok := m[reg]
	if !ok {
		return "", errors.New("invlid byte")
	}

	l := 0
	if w {
		l = 1
	}
	return ss[l], nil

}
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
		var wBit byte = 0b00000000
		if w {
			wBit = 0b00000001
		}
		// Never set d bit, no reason to
		var dBit byte = 0b00000000

		out = append(out,
			[]byte{mov | dBit | wBit,
				0b11000000 | reg2<<3 | reg1,
			}...)
	}
	return out, nil
}

func dasm(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	if len(data)%2 != 0 {
		return nil, errors.New("malformed input stream. Must have an even number of bytes")
	}

	res := "bits 16\n\n"
	for i := 0; i < len(data); i += 2 {
		b1 := data[i]
		b2 := data[i+1]

		if !byteIs(b1, movMask, mov) {
			return nil, errors.New("only mov commands supported")
		}

		d := byteIs(b1, dMask, dMask)
		w := byteIs(b1, wMask, wMask)

		if !byteIs(b2, regMoveMask, regMoveMask) {
			return nil, errors.New("only reg to reg mov commands supported")
		}

		reg := b2 << 2 >> 5
		rm := b2 << 5 >> 5

		regS, err := regLookup(reg, w)
		if err != nil {
			return nil, err
		}

		rmS, err := regLookup(rm, w)
		if err != nil {
			return nil, err
		}
		dest, source := rmS, regS
		if d {
			dest, source = regS, rmS
		}
		res += fmt.Sprintf("mov %s, %s\n", dest, source)
	}
	return []byte(res), nil
}
