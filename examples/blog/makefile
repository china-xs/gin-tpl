
APP_RELATIVE_PATH=$(shell a=`basename $$PWD` && echo $$a)
API_PROTO_FILES=$(shell cd api && find . -name *.proto)
PB_FILES=$(shell cd api && find . -name *.pb.go)

.PHONY: init
init:
	# gotag 书写方式不适合、转换查找文件路径处理错误
	#go install github.com/srikrsna/protoc-gen-gotag
	# pb.go 文件生成包
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	# wire 安装
	go get -u github.com/google/wire/cmd/wire
	# 自动生成 gin 接口 自家用接口还要调整
	go install github.com/china-xs/gin-tpl/cmd/protoc-gen-go-gin@latest
	# gin 绑定 参数增加tag
	go install github.com/favadi/protoc-go-inject-tag
	# 请求参数校验
	go get -u github.com/envoyproxy/protoc-gen-validate
	# 错误状态返回提示语，使用kratos error
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2
	# swagger 依赖包
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

.PHONY: wire
# generate wire
wire:
	cd cmd/$(APP_RELATIVE_PATH) && wire


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