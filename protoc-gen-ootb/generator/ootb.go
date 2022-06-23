package generator

import (
	"github.com/appootb/grpc-gen/v2/protoc-gen-ootb/templates"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

const (
	OutOfTheBoxGenerator = "ootb"
	LangParam            = "lang"
)

type OutOfTheBox struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func New() pgs.Module {
	return &OutOfTheBox{
		ModuleBase: &pgs.ModuleBase{},
	}
}

func (m *OutOfTheBox) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *OutOfTheBox) Name() string {
	return OutOfTheBoxGenerator
}

func (m *OutOfTheBox) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	lang := m.Parameters().Str(LangParam)
	m.Assert(lang != "", "`lang` parameter must be set")

	// Process file-level templates
	tpls := templates.Template(m.Parameters())[lang]
	m.Assert(tpls != nil, "could not find templates for `lang`: ", lang)

	for _, f := range targets {
		m.Push(f.Name().String())

		// TODO: check

		for _, tpl := range tpls {
			out := templates.FilePathFor(tpl)(f, m.ctx, tpl)
			// A nil path means no output should be generated for this file - as controlled by
			// implementation-specific FilePathFor implementations.
			if out != nil {
				m.AddGeneratorTemplateFile(out.String(), tpl, f)
			}
		}

		m.Pop()
	}

	return m.Artifacts()
}
