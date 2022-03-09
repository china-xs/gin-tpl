

.PHONY: init
init:
	# gotag 书写方式不适合、转换查找文件路径处理错误
	#go install github.com/srikrsna/protoc-gen-gotag
	# gin 绑定 参数增加tag
	go install github.com/favadi/protoc-go-inject-tag
	# 请求参数校验
	go get -u github.com/envoyproxy/protoc-gen-validate

.PHONY: wire
# generate wire
wire:
	cd cmd/blog && wire


.PHONY: tidy
tidy:
	go mod tidy -compat=1.17

.PHONY: http
# generate http、swagger、grpc、gin、validate error
http:
	cd api && protoc --proto_path=. \
           --proto_path=../third_party \
           --go_out=paths=source_relative:. \
           --go-grpc_out=paths=source_relative:. \
           --go-gin_out=paths=source_relative:. \
           --validate_out=paths=source_relative,lang=go:. \
           --go-errors_out=paths=source_relative:. \
           --openapiv2_out . \
           --openapiv2_opt logtostderr=true \
           ./auth/login.proto



.PHONY: tag
tag:
	cd api && protoc-go-inject-tag -input="./auth/*pb.go"


.PHONY: gen
gen:
	make http
	make tag