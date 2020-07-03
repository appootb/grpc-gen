package golang

const serviceTpl = `
	var _level{{ .Name.UpperCamelCase }} = map[string][]permission.Audience{
		{{- range $url, $auds := (access .) }}
		"{{ $url }}": {
			{{- range $aud := $auds }}
			permission.Audience_{{ $aud }},
			{{- end }}
		},
		{{- end }}
	}

	{{- if (hasGw .) }}
	type wrapper{{ .Name.UpperCamelCase }}Server struct {
		{{ .Name.UpperCamelCase }}Server
		service.Implementor
	}
	{{ range .Methods }}
	{{ template "method" . }}
	{{ end }}
	{{- end }}

	// Register scoped server.
	func Register{{ .Name.UpperCamelCase }}ScopeServer(auth service.Authenticator, impl service.Implementor, srv {{ .Name.UpperCamelCase }}Server) error {
		// Register service required token level.
		auth.RegisterServiceTokenLevel(_level{{ .Name.UpperCamelCase }})

		// Register scoped gRPC server.
		for _, gRPC := range impl.GetScopedGRPCServer(permission.VisibleScope_{{ (scope .) }}) {
			Register{{ .Name.UpperCamelCase }}Server(gRPC, srv)
		}

		{{- if (hasGw .) }}
		// Register scoped gateway handler server.
		wrapper := wrapper{{ .Name.UpperCamelCase }}Server{
			{{ .Name.UpperCamelCase }}Server: srv,
			Implementor: impl,
		}
		for _, mux := range impl.GetScopedGatewayMux(permission.VisibleScope_{{ (scope .) }}) {
			err := Register{{ .Name.UpperCamelCase }}HandlerServer(impl.Context(), mux, &wrapper)
			if err != nil {
				return err
			}
		}
		{{ else }}
		// No gateway generated.
		{{- end }}
		return nil
	}
`
