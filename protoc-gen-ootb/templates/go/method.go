package golang

const methodTpl = `
{{ if .ClientStreaming }}
	func (w *wrapper{{ serverName .Service }}) {{ .Name.UpperCamelCase }}(srv {{ .Service.Name.UpperCamelCase }}_{{ .Name.UpperCamelCase }}Server) error {
		return w.{{ serverName .Service }}.{{ .Name.UpperCamelCase }}(srv)
	}
{{ else if .ServerStreaming }}
	func (w *wrapper{{ serverName .Service }}) {{ .Name.UpperCamelCase }}(req *{{ inputMessage . }}, srv {{ .Service.Name.UpperCamelCase }}_{{ .Name.UpperCamelCase }}Server) error {
		return w.{{ serverName .Service }}.{{ .Name.UpperCamelCase }}(req, srv)
	}
{{ else }}
	func (w *wrapper{{ serverName .Service }}) {{ .Name.UpperCamelCase }}(ctx context.Context, req *{{ inputMessage . }}) (*{{ outputMessage . }}, error) {
		if w.UnaryInterceptor() == nil {
			return w.{{ serverName .Service }}.{{ .Name.UpperCamelCase }}(ctx, req)
		}
		info := &grpc.UnaryServerInfo{
			Server:     w.{{ serverName .Service }},
			FullMethod: "/{{ .Package.ProtoName }}.{{ .Service.Name }}/{{ .Name.UpperCamelCase }}",
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return w.{{ serverName .Service }}.{{ .Name.UpperCamelCase }}(ctx, req.(*{{ inputMessage . }}))
		}
		resp, err := w.UnaryInterceptor()(ctx, req, info, handler)
		if err != nil {
			return nil, err
		}
		return resp.(*{{ outputMessage . }}), nil
	}
{{ end }}
`
