package golang

const serviceTpl = `
	var _level{{ .Name.UpperCamelCase }} = map[string]permission.TokenLevel{
		{{- range $k, $v := (access .) }}
		"{{ $k }}": permission.TokenLevel_{{ $v }},
		{{- end }}
	}

	// Register scoped server.
	func Register{{ .Name.UpperCamelCase }}ScopeServer(auth service.Authenticator, impl service.Implementor, srv {{ .Name.UpperCamelCase }}Server) error {
		// Register service required token level.
		auth.RegisterServiceTokenLevel(_level{{ .Name.UpperCamelCase }})

		// Register scoped gRPC server.
		for _, grpc := range impl.GetScopedGRPCServer(permission.VisibleScope_{{ (scope .) }}) {
			Register{{ .Name.UpperCamelCase }}Server(grpc, srv)
		}

		{{- if (hasGw .) }}
		// Register scoped gateway handler server.
		for _, mux := range impl.GetScopedGatewayMux(permission.VisibleScope_{{ (scope .) }}) {
			err := Register{{ .Name.UpperCamelCase }}HandlerServer(impl.Context(), mux, srv)
			if err != nil {
				return err
			}
		}
		{{- else }}// No gateway generated.
		{{- end }}
		return nil
	}
`
