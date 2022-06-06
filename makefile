
APP_RELATIVE_PATH=$(shell a=`basename $$PWD` && echo $$a)
API_PROTO_FILES=$(shell cd api && find . -name *.proto)
PB_FILES=$(shell cd api && find . -name *.pb.go)

.PHONY: init
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	# wire 安装
	go get github.com/google/wire/cmd/wire
	#go install github.com/google/wire/cmd/wire

	go install github.com/china-xs/gin-tpl/cmd/protoc-gen-go-gin@latest
	go install github.com/china-xs/gin-tpl/cmd/proto@latest
	# gin 绑定 参数增加tag
	#go install github.com/favadi/protoc-go-inject-tag
	go get -u github.com/favadi/protoc-go-inject-tag
	# 请求参数校验
	go install github.com/envoyproxy/protoc-gen-validate@latest
	# 错误状态返回提示语，使用kratos error
	go get  github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2
	# swagger 依赖包
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

.PHONY: wire
# generate wire
wire:
	cd cmd/$(APP_RELATIVE_PATH) && wire
.PHONY: run
run:
	cd cmd/$(APP_RELATIVE_PATH) && go run .

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
           $(API_PROTO_FILES)

.PHONY: tag
tag:
	cd api && \
	 for name in $(PB_FILES); \
		do \
		protoc-go-inject-tag -input=$$name; \
		done

.PHONY: gen
gen:
	make http
	make tag