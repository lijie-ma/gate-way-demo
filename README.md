# gate way demo

### 只是简单的做个备忘，想深入参考
    https://github.com/grpc-ecosystem/grpc-gateway/blob/main/internal/descriptor/grpc_api_configuration_test.go
    https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/
    https://google.aip.dev/127

### 指令
```
     protoc -I ./proto \
  --go_out ./proto --go_opt paths=source_relative \
  --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
   --grpc-gateway_opt grpc_api_configuration=./proto/helloworld/hello_world.yaml \
  ./proto/helloworld/hello_world.proto
```
#### api配置可以单独的yaml 也可以配置到到proto中

代码中 的 proto/google 文件 来此 https://github.com/googleapis/googleapis 库