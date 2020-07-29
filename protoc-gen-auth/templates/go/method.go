package golang

const methodTpl = `
{{ if .ClientStreaming }}
func (w *wrapper{{ .Service.Name.UpperCamelCase }}Server) {{ .Name.UpperCamelCase }}(srv {{ .Service.Name.UpperCamelCase }}_{{ .Name.UpperCamelCase }}Server) error {
	return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(srv)
}
{{ else if .ServerStreaming }}
func (w *wrapper{{ .Service.Name.UpperCamelCase }}Server) {{ .Name.UpperCamelCase }}(req *{{ goInput . }}, srv {{ .Service.Name.UpperCamelCase }}_{{ .Name.UpperCamelCase }}Server) error {
	return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(req, srv)
}
{{ else }}
func (w *wrapper{{ .Service.Name.UpperCamelCase }}Server) {{ .Name.UpperCamelCase }}(ctx context.Context, req *{{ goInput . }}) (*{{ goOutput . }}, error) {
	if w.UnaryInterceptor() == nil {
		return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     w.{{ .Service.Name.UpperCamelCase }}Server,
		FullMethod: "/{{ .Package.ProtoName }}.{{ .Service.Name }}/{{ .Name.UpperCamelCase }}",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(ctx, req.(*{{ goInput . }}))
	}
	resp, err := w.UnaryInterceptor()(ctx, req, info, handler)
	if err != nil {
		return nil, err
	}
	return resp.(*{{ goOutput . }}), nil
}
{{ end }}
`
