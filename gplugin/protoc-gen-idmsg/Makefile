.PHONY: proto build install

proto:
	protoc --proto_path="./pb/third" --proto_path="./pb" --go_out=. ./pb/*.proto

build:
	go build -mod=vendor -o protoc-gen-idmsg

install:
	go install .