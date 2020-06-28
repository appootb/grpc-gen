package golang

const _ = `
{{ if (wrapperStream .) }}
func (w *wrapper{{ .Service.Name.UpperCamelCase }}Server) {{ .Name.UpperCamelCase }}(req *{{ goInput . }}, srv {{ .Service.Name.UpperCamelCase }}_{{ .Name.UpperCamelCase }}Server) error {
	if w.StreamServerInterceptor() == nil {
		return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(req, srv)
	}
	info := &grpc.StreamServerInfo{
		FullMethod:     "/{{ .Package.ProtoName }}.{{ .Service.Name }}/{{ .Name.UpperCamelCase }}",
		IsServerStream: true,
	}
	handler := func(_ interface{}, stream grpc.ServerStream) error {
		return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(req, &{{ .Service.Name.LowerCamelCase }}{{ .Name.UpperCamelCase }}Server{stream})
	}
	return w.StreamServerInterceptor()(w.{{ .Service.Name.UpperCamelCase }}Server, srv, info, handler)
}
{{ else if or .ServerStreaming .ClientStreaming }}
`

const methodTpl = `
{{ if (wrapperStream .) }}
func (w *wrapper{{ .Service.Name.UpperCamelCase }}Server) {{ .Name.UpperCamelCase }}(req *{{ goInput . }}, srv {{ .Service.Name.UpperCamelCase }}_{{ .Name.UpperCamelCase }}Server) error {
	return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(req, srv)
}
{{ else if or .ServerStreaming .ClientStreaming }}
func (w *wrapper{{ .Service.Name.UpperCamelCase }}Server) {{ .Name.UpperCamelCase }}(srv {{ .Service.Name.UpperCamelCase }}_{{ .Name.UpperCamelCase }}Server) error {
	return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(srv)
}
{{ else }}
func (w *wrapper{{ .Service.Name.UpperCamelCase }}Server) {{ .Name.UpperCamelCase }}(ctx context.Context, req *{{ goInput . }}) (*{{ goOutput . }}, error) {
	if w.UnaryServerInterceptor() == nil {
		return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(ctx, req)
	}
	info := &grpc.UnaryServerInfo{
		Server:     w.{{ .Service.Name.UpperCamelCase }}Server,
		FullMethod: "/{{ .Package.ProtoName }}.{{ .Service.Name }}/{{ .Name.UpperCamelCase }}",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return w.{{ .Service.Name.UpperCamelCase }}Server.{{ .Name.UpperCamelCase }}(ctx, req.(*{{ goInput . }}))
	}
	resp, err := w.UnaryServerInterceptor()(ctx, req, info, handler)
	if err != nil {
		return nil, err
	}
	return resp.(*{{ goOutput . }}), nil
}
{{ end }}
`
