## Shipping package
PROJECT_NAME = example

_cleango:
	@rm -rf proto/*.go
	@rm -rf proto/*.md

_gengo:
	@protoc -Iproto \
		-I${GOPATH}/src/github.com/appootb/substratum/proto \
		-I${GOPATH}/src/github.com/googleapis/googleapis \
		--go_out=paths=source_relative:proto \
		--go-grpc_out=paths=source_relative:proto \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:proto \
		--ootb_out=lang=go,paths=source_relative:proto \
		--validate_out=lang=go,paths=source_relative:proto \
		--markdown_out=paths=source_relative:proto \
		proto/*.proto

gen: clean _gengo

clean: _cleango
