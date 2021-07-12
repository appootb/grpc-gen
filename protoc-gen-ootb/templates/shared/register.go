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
		"comment":           pgs.C80,
		"package":           fns.PackageName,
		"serverName":        fns.ServerName,
		"name":              fns.Name,
		"hasWebApi":         fns.hasWebApi,
		"hasWebStream":      fns.hasWebStream,
		"webStreamPatterns": fns.webStreamPatterns,
		"subjectsName":      fns.subjectsName,
		"serviceSubjects":   fns.serviceSubjects,
		"rolesName":         fns.rolesName,
		"serviceRoles":      fns.serviceRoles,
		"serviceScope":      fns.serviceScope,
		"inputMessage":      fns.inputMessage,
		"outputMessage":     fns.outputMessage,
		"externalPackages":  fns.externalPackages,
	})
}
