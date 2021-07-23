package shared

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func Register(tpl *template.Template, params pgs.Parameters) {
	fns := Funcs{
		Context: pgsgo.InitContext(params),
	}

	tpl.Funcs(map[string]interface{}{
		"package":         fns.PackageName,
		"docFileName":     fns.docFileName,
		"anchorName":      fns.anchorName,
		"tocComment":      fns.tocComment,
		"leadingComment":  fns.leadingComment,
		"trailingComment": fns.trailingComment,
		"webUrl":          fns.webURL,
		"webDoc":          fns.webDoc,
		"inputMessage":    fns.inputMessage,
		"outputMessage":   fns.outputMessage,
		"embedEnums":      fns.embedEnums,
		"embedMessages":   fns.embedMessages,
		"fieldDoc":        fns.fieldDoc,
		"jsonDemo":        fns.messageJSONDemo,
		"jsonWellKnown":   fns.wellKnownJSONDemo,
	})
}
