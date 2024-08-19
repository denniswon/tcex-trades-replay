SHELL:=/bin/bash

proto_clean:
	rm -rfv app/pb

proto_gen:
	mkdir app/pb
	protoc -I app/proto/ --go_out=paths=source_relative:app/pb app/proto/*.proto

setup:
	sh scripts/setup.sh

build:
	go build -o tcex

run: build
	./tcex
