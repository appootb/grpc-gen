package generator

import (
	"github.com/appootb/grpc-gen/v2/protoc-gen-markdown/templates"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

const (
	MarkdownGenerator = "markdown"
)

type Markdown struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
}

func New() pgs.Module {
	return &Markdown{
		ModuleBase: &pgs.ModuleBase{},
	}
}

func (m *Markdown) InitContext(ctx pgs.BuildContext) {
	m.ModuleBase.InitContext(ctx)
	m.ctx = pgsgo.InitContext(ctx.Parameters())
}

func (m *Markdown) Name() string {
	return MarkdownGenerator
}

func (m *Markdown) Execute(targets map[string]pgs.File, _ map[string]pgs.Package) []pgs.Artifact {
	var (
		outDir pgs.FilePath
	)

	// Process file-level templates
	tpls := templates.ProtoTemplate(m.Parameters())

	for _, f := range targets {
		m.Push(f.Name().String())

		// TODO: check

		for _, tpl := range tpls {
			out := templates.FilePathFor(tpl)(f, m.ctx, tpl)
			// A nil path means no output should be generated for this file - as controlled by
			// implementation-specific FilePathFor implementations.
			if out != nil {
				outDir = out.Dir()
				m.AddGeneratorTemplateFile(out.String(), tpl, f)
			}
		}

		m.Pop()
	}

	// README.md
	tocTpl := templates.ReadMeTemplate(m.Parameters())
	tocOut := pgs.JoinPaths(outDir.String(), "README.md")

	m.Push("readme")

	m.AddGeneratorTemplateFile(tocOut.String(), tocTpl, targets)

	m.Pop()

	return m.Artifacts()
}
