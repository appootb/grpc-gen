package main

import (
	"github.com/appootb/grpc-gen/v2/protoc-gen-markdown/generator"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	pgs.
		Init(pgs.DebugEnv("DEBUG_GRPC_GEN")).
		RegisterModule(generator.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
