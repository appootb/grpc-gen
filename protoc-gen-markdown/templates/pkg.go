package templates

import (
	"text/template"

	"github.com/appootb/grpc-gen/protoc-gen-markdown/templates/markdown"
	"github.com/appootb/grpc-gen/protoc-gen-markdown/templates/readme"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type RegisterFn func(tpl *template.Template, params pgs.Parameters)
type FilePathFn func(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath

func ProtoTemplate(params pgs.Parameters) []*template.Template {
	return []*template.Template{
		makeTemplate("md", markdown.Register, params),
	}
}

func ReadMeTemplate(params pgs.Parameters) *template.Template {
	return makeTemplate("readme", readme.Register, params)
}

func FilePathFor(tpl *template.Template) FilePathFn {
	switch tpl.Name() {
	default:
		return func(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath {
			out := ctx.OutputPath(f)
			out = out.SetExt("." + tpl.Name())
			return &out
		}
	}
}

func makeTemplate(ext string, fn RegisterFn, params pgs.Parameters) *template.Template {
	tpl := template.New(ext)
	fn(tpl, params)
	return tpl
}
