package main

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
