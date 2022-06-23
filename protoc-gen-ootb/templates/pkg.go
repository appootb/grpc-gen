package templates

import (
	"text/template"

	golang "github.com/appootb/grpc-gen/v2/protoc-gen-ootb/templates/go"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type RegisterFn func(tpl *template.Template, params pgs.Parameters)
type FilePathFn func(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath

func Template(params pgs.Parameters) map[string][]*template.Template {
	return map[string][]*template.Template{
		"go": {makeTemplate("go", golang.Register, params)},
	}
}

func FilePathFor(tpl *template.Template) FilePathFn {
	switch tpl.Name() {
	default:
		return func(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath {
			out := ctx.OutputPath(f)
			out = out.SetExt(".ootb." + tpl.Name())
			return &out
		}
	}
}

func makeTemplate(ext string, fn RegisterFn, params pgs.Parameters) *template.Template {
	tpl := template.New(ext)
	fn(tpl, params)
	return tpl
}
