package inst

import (
	"bytes"
	"testing"
)

var (
	asm1 []byte = []byte(`bits 16

mov cx, bx
`)
	bin1 []byte = []byte{0b10001001, 0b11011001}
	asm2 []byte = []byte(`bits 16

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
	bin2 []byte = []byte{
		0b10001001, 0b11011001, 0b10001000, 0b11100101,
		0b10001001, 0b11011010, 0b10001001, 0b11011110,
		0b10001001, 0b11111011, 0b10001000, 0b11001000,
		0b10001000, 0b11101101, 0b10001001, 0b11000011,
		0b10001001, 0b11110011, 0b10001001, 0b11111100,
		0b10001001, 0b11000101}

	asm3 []byte = []byte(`bits 16

mov si, bx
mov dh, al
mov cl, 12
mov ch, -12
mov cx, 12
mov cx, -12
mov dx, 3948
mov dx, -3948
mov al, [bx + si]
mov bx, [bp + di]
mov dx, [bp]
mov ah, [bx + si + 4]
mov al, [bx + si + 4999]
mov [bx + di], cx
mov [bp + si], cl
mov [bp], ch
`)
	bin3 []byte = []byte{
		0b10001001, 0b11011110, 0b10001000, 0b11000110, 0b10110001, 0b00001100,
		0b10110101, 0b11110100, 0b10111001, 0b00001100, 0b00000000, 0b10111001,
		0b11110100, 0b11111111, 0b10111010, 0b01101100, 0b00001111, 0b10111010,
		0b10010100, 0b11110000, 0b10001010, 0b00000000, 0b10001011, 0b00011011,
		0b10001011, 0b01010110, 0b00000000, 0b10001010, 0b01100000, 0b00000100,
		0b10001010, 0b10000000, 0b10000111, 0b00010011, 0b10001001, 0b00001001,
		0b10001000, 0b00001010, 0b10001000, 0b01101110, 0b00000000}
)

func TestOpType(t *testing.T) {
	t.Run("mov", func(t *testing.T) {
		t.Run("MovTypeRegMemToFromReg", func(t *testing.T) {
			if opType(0b10001011) != OpTypeMovRegMemToFromReg {
				t.FailNow()
			}
			if opType(0b10001010) != OpTypeMovRegMemToFromReg {
				t.FailNow()
			}
			if opType(0b10001001) != OpTypeMovRegMemToFromReg {
				t.FailNow()
			}
			if opType(0b10001000) != OpTypeMovRegMemToFromReg {
				t.FailNow()
			}
		})

		t.Run("MovTypeImmToRegOrMem", func(t *testing.T) {
			if opType(0b11000111) != OpTypeMovImmToRegOrMem {
				t.FailNow()
			}
			if opType(0b11000110) != OpTypeMovImmToRegOrMem {
				t.FailNow()
			}
		})

		t.Run("MovTypeImmToReg", func(t *testing.T) {
			if opType(0b10110000) != OpTypeMovImmToReg {
				t.FailNow()
			}
			if opType(0b10111000) != OpTypeMovImmToReg {
				t.FailNow()
			}
			if opType(0b10111100) != OpTypeMovImmToReg {
				t.FailNow()
			}
			if opType(0b10111110) != OpTypeMovImmToReg {
				t.FailNow()
			}
			if opType(0b10111111) != OpTypeMovImmToReg {
				t.FailNow()
			}
		})

		t.Run("MovTypeMemToAcc", func(t *testing.T) {
			if opType(0b10100000) != OpTypeMovMemToAcc {
				t.FailNow()
			}
			if opType(0b10100001) != OpTypeMovMemToAcc {
				t.FailNow()
			}
		})

		t.Run("MovTypeAccToMem", func(t *testing.T) {
			if opType(0b10100010) != OpTypeMovAccToMem {
				t.FailNow()
			}
			if opType(0b10100011) != OpTypeMovAccToMem {
				t.FailNow()
			}
		})

		t.Run("MovTypeRegOrMemToSegReg", func(t *testing.T) {
			if opType(0b10001110) != OpTypeMovRegOrMemToSegReg {
				t.FailNow()
			}
		})

		t.Run("MovTypeSegRegToRegMemory", func(t *testing.T) {
			if opType(0b10001100) != OpTypeMovSegRegToRegMemory {
				t.FailNow()
			}
		})
	})

	t.Run("OpTypeImmToRegOrMem", func(t *testing.T) {
		for _, op := range []byte{0b10000000, 0b10000001, 0b10000010, 0b10000011} {
			if opType(op) != OpTypeImmToRegOrMem {
				t.FailNow()
			}
		}
	})
	t.Run("add", func(t *testing.T) {
		t.Run("OpTypeAddRegMemWithReg", func(t *testing.T) {
			for _, op := range []byte{0b00000000, 0b00000001, 0b00000010, 0b00000011} {
				if opType(op) != OpTypeAddRegMemWithReg {
					t.FailNow()
				}
			}
		})
		t.Run("OpTypeAddImmToAcc", func(t *testing.T) {
			for _, op := range []byte{0b00000100, 0b00000101} {
				if opType(op) != OpTypeAddImmToAcc {
					t.FailNow()
				}
			}
		})
	})

	t.Run("sub", func(t *testing.T) {
		t.Run("OpTypeSubRegMemWithReg", func(t *testing.T) {
			for _, op := range []byte{0b00101000, 0b00101001, 0b00101010, 0b00101011} {
				if opType(op) != OpTypeSubRegMemWithReg {
					t.FailNow()
				}
			}
		})
		t.Run("OpTypeSubImmToAcc", func(t *testing.T) {
			for _, op := range []byte{0b00101100, 0b00101101} {
				if opType(op) != OpTypeSubImmToAcc {
					t.FailNow()
				}
			}
		})
	})

	t.Run("cmp", func(t *testing.T) {
		t.Run("OpTypeCmpRegMemWithReg", func(t *testing.T) {
			for _, op := range []byte{0b00111000, 0b00111001, 0b00111010, 0b00111011} {
				if opType(op) != OpTypeCmpRegMemWithReg {
					t.FailNow()
				}
			}
		})
		t.Run("OpTypeCmpImmToAcc", func(t *testing.T) {
			for _, op := range []byte{0b00111100, 0b00111101} {
				if opType(op) != OpTypeCmpImmToAcc {
					t.FailNow()
				}
			}
		})
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
		{
			name:            "more movs",
			inFile:          bin3,
			expectedOutFile: asm3,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Dasm(tc.inFile)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if !bytes.Equal(out, tc.expectedOutFile) {
				t.Fatalf("expected asm(%s) to return %s", tc.inFile, tc.expectedOutFile)
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
