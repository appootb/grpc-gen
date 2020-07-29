package shared

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func Register(tpl *template.Template, params pgs.Parameters) {
	fn := Func{
		Context: pgsgo.InitContext(params),
	}

	tpl.Funcs(map[string]interface{}{
		"pkg":          fn.PackageName,
		"imports":      fn.GolangImports,
		"scope":        fn.Scope,
		"hasGw":        fn.GatewayDefined,
		"access":       fn.Access,
		"goInput":      fn.GolangInputMessageName,
		"goOutput":     fn.GolangOutputMessageName,
		"hasWebsocket": fn.WebsocketDefined,
		"isWebsocket":  fn.IsWebsocket,
		"websocketURL": fn.WebsocketURLPattern,
	})
}
