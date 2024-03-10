package main

import (
	"bytes"
	"fmt"
	"testing"
)

var asm1 []byte = []byte(`bits 16

mov cx, bx
`)

var bin1 []byte = []byte{0b10001001, 0b11011001}

var asm2 []byte = []byte(`bits 16

mov cx, bx
mov ch, ah
mov dx, bx
mov si, bx
mov bx, di
mov al, cl
mov ch, ch
mov bx, ax
mov bx, si
mov sp, di
mov bp, ax
`)

var bin2 []byte = []byte{
	0b10001001, 0b11011001, 0b10001000, 0b11100101,
	0b10001001, 0b11011010, 0b10001001, 0b11011110,
	0b10001001, 0b11111011, 0b10001000, 0b11001000,
	0b10001000, 0b11101101, 0b10001001, 0b11000011,
	0b10001001, 0b11110011, 0b10001001, 0b11111100,
	0b10001001, 0b11000101}

func TestAsm(t *testing.T) {
	type testCase struct {
		name            string
		inFile          []byte
		expectedOutFile []byte
	}

	tcs := []testCase{
		{
			name:            "single instruction",
			inFile:          bin1,
			expectedOutFile: asm1,
		},
		{
			name:            "multiple instructions",
			inFile:          bin2,
			expectedOutFile: asm2,
		},
	}

	for _, tc := range tcs {
		t.Run("tc.name", func(t *testing.T) {
			out, err := dasm(tc.inFile)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if !bytes.Equal(out, tc.expectedOutFile) {
				t.Fatalf("expected asm(%s) to return %s", tc.inFile, tc.expectedOutFile)
			}

			in, err := asm(out)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if !bytes.Equal(in, tc.inFile) {
				s := ""
				for _, b := range tc.inFile {
					s += fmt.Sprintf("%08b ", b)
				}

				s2 := ""
				for _, b := range in {
					s2 += fmt.Sprintf("%08b ", b)
				}

				t.Fatalf("expected dasm(%s) to return %s but instead returned %s", out, s, s2)
			}
		})
	}

}
