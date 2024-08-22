SHELL:=/bin/bash

# protobuf generation targets

proto_clean:
	rm -rfv app/pb

proto_gen:
	mkdir app/pb
	protoc -I app/proto/ --go_out=paths=source_relative:app/pb app/proto/*.proto

# environment setup targets

setup:
	sh scripts/setup.sh

# server targets

build:
	go build -o tcex

run: build
	./tcex

# demo client targets

client:
	cd demo && yarn dev

browser:
	open http://localhost:3000

demo: build
	make -j 3 run client browser

# test targets

run_test:
	cd test && yarn start

test: build
	make -j 2 run run_test
