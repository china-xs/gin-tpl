

.PHONY: init
init:
	go install github.com/srikrsna/protoc-gen-gotag
	go install github.com/favadi/protoc-go-inject-tag

.PHONY: wire
# generate wire
wire:
	cd cmd/blog && wire


.PHONY: tidy
tidy:
	go mod tidy -compat=1.17

.PHONY: http
# generate http code
http:
	cd api && protoc --proto_path=. \
           --proto_path=../third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           --go-gin_out=paths=source_relative:. \
           --validate_out=paths=source_relative,lang=go:. \
           --openapiv2_out . \
           --openapiv2_opt logtostderr=true \
           ./auth/login.proto

.PHONY: tag
tag:
	cd api && protoc-go-inject-tag -input="./auth/*pb.go"