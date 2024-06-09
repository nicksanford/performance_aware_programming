package inst

import (
	"bytes"
	"testing"
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
			name:   "single instruction",
			inFile: []byte{0b10001001, 0b11011001},
			expectedOutFile: []byte(`bits 16

mov cx, bx
`),
		},
		{
			name: "multiple instructions",
			inFile: []byte{
				0b10001001, 0b11011001, 0b10001000, 0b11100101,
				0b10001001, 0b11011010, 0b10001001, 0b11011110,
				0b10001001, 0b11111011, 0b10001000, 0b11001000,
				0b10001000, 0b11101101, 0b10001001, 0b11000011,
				0b10001001, 0b11110011, 0b10001001, 0b11111100,
				0b10001001, 0b11000101},
			expectedOutFile: []byte(`bits 16

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
`),
		},
		{
			name: "more movs",
			inFile: []byte{
				0b10001001, 0b11011110, 0b10001000, 0b11000110, 0b10110001, 0b00001100,
				0b10110101, 0b11110100, 0b10111001, 0b00001100, 0b00000000, 0b10111001,
				0b11110100, 0b11111111, 0b10111010, 0b01101100, 0b00001111, 0b10111010,
				0b10010100, 0b11110000, 0b10001010, 0b00000000, 0b10001011, 0b00011011,
				0b10001011, 0b01010110, 0b00000000, 0b10001010, 0b01100000, 0b00000100,
				0b10001010, 0b10000000, 0b10000111, 0b00010011, 0b10001001, 0b00001001,
				0b10001000, 0b00001010, 0b10001000, 0b01101110, 0b00000000},
			expectedOutFile: []byte(`bits 16

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
`),
		},
		// {
		// 	name: "challenge movs",
		// 	inFile: []byte{
		// 		0b10001011, 0b01000001, 0b11011011, 0b10001001, 0b10001100, 0b11010100,
		// 		0b11111110, 0b10001011, 0b01010111, 0b11100000, 0b11000110, 0b00000011,
		// 		0b00000111, 0b11000111, 0b10000101, 0b10000101, 0b00000011, 0b01011011,
		// 		0b00000001, 0b10001011, 0b00101110, 0b00000101, 0b00000000, 0b10001011,
		// 		0b00011110, 0b10000010, 0b00001101, 0b10100001, 0b11111011, 0b00001001,
		// 		0b10100001, 0b00010000, 0b00000000, 0b10100011, 0b11111010, 0b00001001,
		// 		0b10100011, 0b00001111, 0b00000000, 0b00001010},
		// 	expectedOutFile: []byte(`bits 16

		// mov ax, [bx + di - 37]
		// mov [si - 300], cx
		// mov dx, [bx - 32]
		// mov [bp + di], byte 7
		// mov [di + 901], word 347
		// mov bp, [5]
		// mov bx, [3458]
		// mov ax, [2555]
		// mov ax, [16]
		// mov [2554], ax
		// mov [15], ax
		// `),
		// },
		{
			name: "add sub cmp jnz",
			inFile: []byte{
				0b00000011, 0b00011000, 0b00000011, 0b01011110, 0b00000000, 0b10000011,
				0b11000110, 0b00000010, 0b10000011, 0b11000101, 0b00000010, 0b10000011,
				0b11000001, 0b00001000, 0b00000011, 0b01011110, 0b00000000, 0b00000011,
				0b01001111, 0b00000010, 0b00000010, 0b01111010, 0b00000100, 0b00000011,
				0b01111011, 0b00000110, 0b00000001, 0b00011000, 0b00000001, 0b01011110,
				0b00000000, 0b00000001, 0b01011110, 0b00000000, 0b00000001, 0b01001111,
				0b00000010, 0b00000000, 0b01111010, 0b00000100, 0b00000001, 0b01111011,
				0b00000110, 0b10000000, 0b00000111, 0b00100010, 0b10000011, 0b10000010,
				0b11101000, 0b00000011, 0b00011101, 0b00000011, 0b01000110, 0b00000000,
				0b00000010, 0b00000000, 0b00000001, 0b11011000, 0b00000000, 0b11100000,
				0b00000101, 0b11101000, 0b00000011, 0b00000100, 0b11100010, 0b00000100,
				0b00001001, 0b00101011, 0b00011000, 0b00101011, 0b01011110, 0b00000000,
				0b10000011, 0b11101110, 0b00000010, 0b10000011, 0b11101101, 0b00000010,
				0b10000011, 0b11101001, 0b00001000, 0b00101011, 0b01011110, 0b00000000,
				0b00101011, 0b01001111, 0b00000010, 0b00101010, 0b01111010, 0b00000100,
				0b00101011, 0b01111011, 0b00000110, 0b00101001, 0b00011000, 0b00101001,
				0b01011110, 0b00000000, 0b00101001, 0b01011110, 0b00000000, 0b00101001,
				0b01001111, 0b00000010, 0b00101000, 0b01111010, 0b00000100, 0b00101001,
				0b01111011, 0b00000110, 0b10000000, 0b00101111, 0b00100010, 0b10000011,
				0b00101001, 0b00011101, 0b00101011, 0b01000110, 0b00000000, 0b00101010,
				0b00000000, 0b00101001, 0b11011000, 0b00101000, 0b11100000, 0b00101101,
				0b11101000, 0b00000011, 0b00101100, 0b11100010, 0b00101100, 0b00001001,
				0b00111011, 0b00011000, 0b00111011, 0b01011110, 0b00000000, 0b10000011,
				0b11111110, 0b00000010, 0b10000011, 0b11111101, 0b00000010, 0b10000011,
				0b11111001, 0b00001000, 0b00111011, 0b01011110, 0b00000000, 0b00111011,
				0b01001111, 0b00000010, 0b00111010, 0b01111010, 0b00000100, 0b00111011,
				0b01111011, 0b00000110, 0b00111001, 0b00011000, 0b00111001, 0b01011110,
				0b00000000, 0b00111001, 0b01011110, 0b00000000, 0b00111001, 0b01001111,
				0b00000010, 0b00111000, 0b01111010, 0b00000100, 0b00111001, 0b01111011,
				0b00000110, 0b10000000, 0b00111111, 0b00100010, 0b10000011, 0b00111110,
				0b11100010, 0b00010010, 0b00011101, 0b00111011, 0b01000110, 0b00000000,
				0b00111010, 0b00000000, 0b00111001, 0b11011000, 0b00111000, 0b11100000,
				0b00111101, 0b11101000, 0b00000011, 0b00111100, 0b11100010, 0b00111100,
				0b00001001, 0b01110101, 0b00000010, 0b01110101, 0b11111100, 0b01110101,
				0b11111010, 0b01110101, 0b11111100, 0b01110100, 0b11111110, 0b01111100,
				0b11111100, 0b01111110, 0b11111010, 0b01110010, 0b11111000, 0b01110110,
				0b11110110, 0b01111010, 0b11110100, 0b01110000, 0b11110010, 0b01111000,
				0b11110000, 0b01110101, 0b11101110, 0b01111101, 0b11101100, 0b01111111,
				0b11101010, 0b01110011, 0b11101000, 0b01110111, 0b11100110, 0b01111011,
				0b11100100, 0b01110001, 0b11100010, 0b01111001, 0b11100000, 0b11100010,
				0b11011110, 0b11100001, 0b11011100, 0b11100000, 0b11011010, 0b11100011,
				0b11011000,
			},
			expectedOutFile: []byte(`bits 16

add bx, [bx+si]
add bx, [bp]
add si, 2
add bp, 2
add cx, 8
add bx, [bp + 0]
add cx, [bx + 2]
add bh, [bp + si + 4]
add di, [bp + di + 6]
add [bx+si], bx
add [bp], bx
add [bp + 0], bx
add [bx + 2], cx
add [bp + si + 4], bh
add [bp + di + 6], di
add byte [bx], 34
add word [bp + si + 1000], 29
add ax, [bp]
add al, [bx + si]
add ax, bx
add al, ah
add ax, 1000
add al, -30
add al, 9
sub bx, [bx+si]
sub bx, [bp]
sub si, 2
sub bp, 2
sub cx, 8
sub bx, [bp + 0]
sub cx, [bx + 2]
sub bh, [bp + si + 4]
sub di, [bp + di + 6]
sub [bx+si], bx
sub [bp], bx
sub [bp + 0], bx
sub [bx + 2], cx
sub [bp + si + 4], bh
sub [bp + di + 6], di
sub byte [bx], 34
sub word [bx + di], 29
sub ax, [bp]
sub al, [bx + si]
sub ax, bx
sub al, ah
sub ax, 1000
sub al, -30
sub al, 9
cmp bx, [bx+si]
cmp bx, [bp]
cmp si, 2
cmp bp, 2
cmp cx, 8
cmp bx, [bp + 0]
cmp cx, [bx + 2]
cmp bh, [bp + si + 4]
cmp di, [bp + di + 6]
cmp [bx+si], bx
cmp [bp], bx
cmp [bp + 0], bx
cmp [bx + 2], cx
cmp [bp + si + 4], bh
cmp [bp + di + 6], di
cmp byte [bx], 34
cmp word [4834], 29
cmp ax, [bp]
cmp al, [bx + si]
cmp ax, bx
cmp al, ah
cmp ax, 1000
cmp al, -30
cmp al, 9
test_label0:
jnz test_label1
jnz test_label0
test_label1:
jnz test_label0
jnz test_label1
label:
je label
jl label
jle label
jb label
jbe label
jp label
jo label
js label
jne label
jnl label
jg label
jnb label
ja label
jnp label
jno label
jns label
loop label
loopz label
loopnz label
jcxz label`),
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
