# 编写脚本进入脚本当前目录
cd "$(dirname "$0")"

# 执行protoc（老版本）
protoc -I. --go_out=plugins=grpc:. demo.proto
protoc -I. --grpc-gateway_out=./ demo.proto


# https://github.com/grpc-ecosystem/grpc-gateway
# protoc -I . demo.proto --go-grpc_out=./proto/demo
# protoc-gen-go-grpc: program not found or is not executable