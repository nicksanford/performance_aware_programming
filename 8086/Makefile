.PHONY: build

DATA_URL := https://raw.githubusercontent.com/cmuratori/computer_enhance/main/perfaware/part1

build:
	go build -o bin/8086 ./cmd/main.go 

clean:
	rm -rf bin data

test:
	go test -v -race

data:
	wget \
		$(DATA_URL)/listing_0037_single_register_mov.asm \
		$(DATA_URL)/listing_0037_single_register_mov  \
		$(DATA_URL)/listing_0038_many_register_mov.asm \
		$(DATA_URL)/listing_0038_many_register_mov \
		$(DATA_URL)/listing_0039_more_movs.asm \
		$(DATA_URL)/listing_0039_more_movs \
		$(DATA_URL)/listing_0041_add_sub_cmp_jnz.asm \
		$(DATA_URL)/listing_0041_add_sub_cmp_jnz \
		-P ./data
