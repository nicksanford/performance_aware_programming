.PHONY: build clan test watch

DATA_URL := https://raw.githubusercontent.com/cmuratori/computer_enhance/main/perfaware/part1

build:
	go build -o bin/8086 ./cmd/main.go 

clean:
	rm -rf bin data

test:
	go test -v -race

watch:
	find . | entr -s 'go test -v -race -run=TestAsm/add_sub_cmp_jnz'

data:
	wget \
		$(DATA_URL)/listing_0037_single_register_mov.asm \
		$(DATA_URL)/listing_0037_single_register_mov  \
		$(DATA_URL)/listing_0038_many_register_mov.asm \
		$(DATA_URL)/listing_0038_many_register_mov \
		$(DATA_URL)/listing_0039_more_movs.asm \
		$(DATA_URL)/listing_0040_challenge_movs \
		$(DATA_URL)/listing_0040_challenge_movs.asm \
		$(DATA_URL)/listing_0041_add_sub_cmp_jnz.asm \
		$(DATA_URL)/listing_0041_add_sub_cmp_jnz \
		$(DATA_URL)/listing_0043_immediate_movs \
		$(DATA_URL)/listing_0043_immediate_movs.asm \
		$(DATA_URL)/listing_0044_register_movs \
		$(DATA_URL)/listing_0044_register_movs.asm \
		$(DATA_URL)/listing_0046_add_sub_cmp \
		$(DATA_URL)/listing_0046_add_sub_cmp.asm \
		-P ./data
