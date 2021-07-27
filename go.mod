module github.com/appootb/grpc-gen

go 1.14

require (
	github.com/appootb/substratum v0.0.0-20210727065501-c8636c517cb1
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.14.6
	github.com/iancoleman/strcase v0.1.2
	github.com/lyft/protoc-gen-star v0.5.2
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.29.1 // indirect
	google.golang.org/protobuf v1.23.1-0.20200526195155-81db48ad09cc // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
