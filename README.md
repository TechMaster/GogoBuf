# Hello World

This is hello world using micro

## Contents

- main.go - is the main definition of the service, handler and client
- proto - contains the protobuf definition of the API

## Dependencies

Install the following

- [consul](https://www.consul.io/intro/getting-started/install.html)
- [micro](https://github.com/micro/micro)
- [protoc-gen-micro](https://github.com/micro/protoc-gen-micro)

## Run Consul
```shell
consul agent -dev -ui
```

## Run Service

```shell
go run main.go
```

## Query Service

```
$ micro call greeter Greeter.Hello '{"name": "John"}'

{
	"greeting": "Hello John"
}
```

## Call through micro api that reverse proxy http request to gRPC call
### Start micro api that listens at 8080
```shel
$ micro api
```

### Make http request using curl
```shell
$ curl -d 'service=greeter' \
       -d 'method=Greeter.Hello' \
       -d 'request={"name": "John"}' \
       http://localhost:8080/rpc

{"greeting":"Hello John"}
```

## Use gogobuf để bổ xung các annotation cho trường
Cài đặt

```
go get github.com/gogo/protobuf/protoc-gen-gogofast
go get github.com/gogo/protobuf/protoc-gen-gogofaster
go get github.com/gogo/protobuf/protoc-gen-gogoslick

```
### Phân biệt các lệnh generate
Lệnh này không hỗ trợ extension
```
protoc --gofast_out=. --micro_out=. greeter.proto
```


### Bổ xung annotation vào từng trường
1. Phải chuyển syntax sang proto2
2. Thêm vào đây ```import "github.com/gogo/protobuf/gogoproto/gogo.proto";```
3. Sau mỗi trường thêm các khối như
```
[(gogoproto.nullable) = false,
(gogoproto.jsontag) = "MyField1",
(gogoproto.moretags) = "pg:\",array\", sql:\",pk\""];

```
Đây là nội dung file greeter.proto
```proto
syntax = "proto2";

import "github.com/gogo/protobuf/gogoproto/gogo.proto";

service Greeter {
	rpc Hello(HelloRequest) returns (HelloResponse) {}
	rpc GoodBye(HelloRequest) returns (HelloResponse) {}
}

message HelloRequest {
	optional string name = 1 [(gogoproto.nullable) = false,
							 (gogoproto.jsontag) = "MyField1",
							 (gogoproto.moretags) = "pg:\",array\", sql:\",pk\""];
}

message Permission {
	// Tên của service (ví dụ 'blog')
	required string service = 1 [(gogoproto.moretags) = "sql:\",pk\""];

	// Các phương thức của service
	required string method = 2 [(gogoproto.moretags) = "sql:\",pk\""];

	// Mảng các quyền hạn (Role) cho phép truy cập vào method (nếu null tức là public method)
	repeated int32 roles = 3 [(gogoproto.moretags) = "pg:\",array\""];

	// Method có phải là public hay không
	required bool is_public = 4;
}


message HelloResponse {
	optional string greeting = 2;
}

```
Chạy lệnh này để generate thành công annotated field
```
protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogoslick_out=. --micro_out=. \
greeter.proto
```

# Muốn bổ xung descriptor cho rpc
1. Import ```import "google/protobuf/descriptor.proto";```


```proto
syntax = "proto2";

import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "google/protobuf/descriptor.proto";

extend google.protobuf.MethodOptions {
	optional string description = 50056;
	optional bool internalMethod = 50057;
}

service Greeter {
	rpc Hello(HelloRequest) returns (HelloResponse) {}
	rpc GoodBye(HelloRequest) returns (HelloResponse) {}
	rpc VuiQua(HelloRequest) returns (HelloResponse) {
		option(description) = "VuiQua Day";
		option(internalMethod) = false;
	}
}

```
```
protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf \
--micro_out=. \
--gogoslick_out=\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:. \
greeter.proto
```
