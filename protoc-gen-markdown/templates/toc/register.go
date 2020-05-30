package toc

import (
	"text/template"

	"github.com/appootb/grpc-gen/protoc-gen-markdown/templates/shared"
	pgs "github.com/lyft/protoc-gen-star"
)

func Register(tpl *template.Template, params pgs.Parameters) {
	shared.Register(tpl, params)
	template.Must(tpl.Parse(fileTpl))
}
