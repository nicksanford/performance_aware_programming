package inst

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

func TestMovType(t *testing.T) {
	t.Run("MovTypeRegMemToFromReg", func(t *testing.T) {
		if movType(0b10001011) != MovTypeRegMemToFromReg {
			t.FailNow()
		}
		if movType(0b10001010) != MovTypeRegMemToFromReg {
			t.FailNow()
		}
		if movType(0b10001001) != MovTypeRegMemToFromReg {
			t.FailNow()
		}
		if movType(0b10001000) != MovTypeRegMemToFromReg {
			t.FailNow()
		}
	})

	t.Run("MovTypeImmToRegOrMem", func(t *testing.T) {
		if movType(0b11000111) != MovTypeImmToRegOrMem {
			t.FailNow()
		}
		if movType(0b11000110) != MovTypeImmToRegOrMem {
			t.FailNow()
		}
	})

	t.Run("MovTypeImmToReg", func(t *testing.T) {
		if movType(0b10110000) != MovTypeImmToReg {
			t.FailNow()
		}
		if movType(0b10111000) != MovTypeImmToReg {
			t.FailNow()
		}
		if movType(0b10111100) != MovTypeImmToReg {
			t.FailNow()
		}
		if movType(0b10111110) != MovTypeImmToReg {
			t.FailNow()
		}
		if movType(0b10111111) != MovTypeImmToReg {
			t.FailNow()
		}
	})

	t.Run("MovTypeMemToAcc", func(t *testing.T) {
		if movType(0b10100000) != MovTypeMemToAcc {
			t.FailNow()
		}
		if movType(0b10100001) != MovTypeMemToAcc {
			t.FailNow()
		}
	})

	t.Run("MovTypeAccToMem", func(t *testing.T) {
		if movType(0b10100010) != MovTypeAccToMem {
			t.FailNow()
		}
		if movType(0b10100011) != MovTypeAccToMem {
			t.FailNow()
		}
	})

	t.Run("MovTypeRegOrMemToSegReg", func(t *testing.T) {
		if movType(0b10001110) != MovTypeRegOrMemToSegReg {
			t.FailNow()
		}
	})

	t.Run("MovTypeSegRegToRegMemory", func(t *testing.T) {
		if movType(0b10001100) != MovTypeSegRegToRegMemory {
			t.FailNow()
		}
	})
}
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
			out, err := Dasm(tc.inFile)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if !bytes.Equal(out, tc.expectedOutFile) {
				t.Fatalf("expected asm(%s) to return %s", tc.inFile, tc.expectedOutFile)
			}

			in, err := Asm(out)
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

func TestDasmTemp(t *testing.T) {
	expected := `bits 16

mov si, bx
mov dh, al
`
	res, err := Dasm([]byte{0b10001001, 0b11011110, 0b10001000, 0b11000110})
	if err != nil {
		t.Logf("expected err to be nil, instead was %s", err.Error())
		t.FailNow()
	}
	if string(res) != expected {
		t.Logf("expected inst.Dasm([]byte{0b10001001, 0b11011110, 0b10001000, 0b11000110}) to equal:\n%s", expected)
		t.FailNow()
	}

}
