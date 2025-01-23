.PHONY: proto

registry_prefix=""

# public param
out_bin=mist_server
IpFile=cmd/cluster/data/IP2LOCATION-LITE-DB11.BIN

GV = '$(shell date "+%Y-%m-%d %H:%M:%S") $(shell git rev-parse HEAD)'
NOW = $(shell date "+%Y-%m-%d %H:%M:%S")
PROJECT_DIR = $(shell dirname "$0")

# commands
proto:
	protoc \
		--proto_path="./protocol/proto/base" \
		--proto_path="./protocol/proto" \
		--proto_path="./protocol/proto_src" \
		--go_out=. \
		--go-grpc_out=. \
		--idmsg_out=./protocol/pb \
		--swagger_out=allow_merge=true:. \
		protocol/proto/base/*.proto \
		protocol/proto/*.proto

	cp ./protocol/pb/msg.ts ./msg.ts
	cp ./protocol/pb/msg.ts ./msg_mini.ts

	protoc \
		--proto_path="./protocol/proto/base" \
		--proto_path="./protocol/proto" \
		--proto_path="./protocol/proto_src" \
		--gotag_out=xxx="bson+\"-\"":. \
		protocol/proto/base/*.proto \
		protocol/proto/*.proto


