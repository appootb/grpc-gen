## Shipping package
PROJECT_NAME = example

_cleango:
	@rm -rf proto/*.go
	@rm -rf proto/*.md

_gengo:
	@protoc -Iproto \
		-I${GOPATH}/src/github.com/appootb/substratum/proto \
		-I${GOPATH}/src/github.com/googleapis/googleapis \
		--go_out=plugins=grpc:proto \
		--grpc-gateway_out=logtostderr=true:proto \
		--ootb_out="lang=go:proto" \
		--markdown_out=":proto" \
		proto/*.proto

gen: clean _gengo

clean: _cleango
