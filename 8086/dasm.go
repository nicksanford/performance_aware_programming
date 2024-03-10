package main

import (
	"errors"
	"fmt"
)

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
