package golang

const webStreamTpl = `
func Register{{ .Name.UpperCamelCase }}WsHandlerServer(mux *runtime.ServeMux, srv {{ serverName . }}, streamInterceptor grpc.StreamServerInterceptor) error {
	{{ range .Methods }}
		{{ if .ClientStreaming }}
			{{ template "clientStream" . }}
		{{ else if .ServerStreaming }}
			{{ template "serverStream" . }}
		{{ end }}
	{{ end }}
	return nil
}

var (
	{{- range $name, $rule := (webStreamPatterns .) }}
		{{ $name }} = runtime.MustPattern(runtime.NewPattern(1, []int{ {{ range $rule.OpCodes }}{{ . }},{{ end }} }, []string{ {{ range $rule.Pool }}"{{ . }}",{{ end }} }, ""))
	{{- end }}
)
`
