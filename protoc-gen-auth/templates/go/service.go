package golang

const serviceTpl = `
{{- $hasGateway := (hasGw .) }}
{{- $hasWebsocket := (hasWebsocket .) }}
var _{{ .Name.LowerCamelCase }}ServiceSubjects = map[string][]permission.Subject{
	{{- range $url, $subs := (access .) }}
	"{{ $url }}": {
		{{- range $sub := $subs }}
		permission.Subject_{{ $sub }},
		{{- end }}
	},
	{{- end }}
}

var _{{ .Name.LowerCamelCase }}ServiceRoles = map[string][]string{
	{{- range $url, $roles := (serviceRoles .) }}
	"{{ $url }}": {
		{{- range $role := $roles }}
		"{{ $role }}",
		{{- end }}
	},
	{{- end }}
}

{{- if $hasGateway }}
type wrapper{{ .Name.UpperCamelCase }}Server struct {
	{{ .Name.UpperCamelCase }}Server
	service.Implementor
}
{{ range .Methods }}
{{ template "method" . }}
{{ end }}
{{- end }}

// Register scoped server.
func Register{{ .Name.UpperCamelCase }}ScopeServer(component string, auth service.Authenticator, impl service.Implementor, srv {{ .Name.UpperCamelCase }}Server) error {
	// Register service required subjects.
	auth.RegisterServiceSubjects(component, _{{ .Name.LowerCamelCase }}ServiceSubjects, _{{ .Name.LowerCamelCase }}ServiceRoles)

	// Register scoped gRPC server.
	for _, gRPC := range impl.GetGRPCServer(permission.VisibleScope_{{ (scope .) }}) {
		Register{{ .Name.UpperCamelCase }}Server(gRPC, srv)
	}

	{{- if $hasGateway }}
	// Register scoped gateway handler server.
	wrapper := wrapper{{ .Name.UpperCamelCase }}Server{
		{{ .Name.UpperCamelCase }}Server: srv,
		Implementor: impl,
	}
	{{- end }}
	{{- if or $hasGateway $hasWebsocket }}
	for _, mux := range impl.GetGatewayMux(permission.VisibleScope_{{ (scope .) }}) {
		{{- if $hasGateway }}
		// Register gateway handler.
		if err := Register{{ .Name.UpperCamelCase }}HandlerServer(impl.Context(), mux, &wrapper); err != nil {
			return err
		}
		{{- end }}
		{{- if $hasWebsocket }}
		// Register websocket handler.
		if err := Register{{ .Name.UpperCamelCase }}WsHandlerServer(mux, srv, impl.StreamInterceptor()); err != nil {
			return err
		}
		{{- end }}
	}
	{{ else }}
	// No gateway generated.
	{{- end }}
	return nil
}

{{ if $hasWebsocket }}
func Register{{ .Name.UpperCamelCase }}WsHandlerServer(mux *runtime.ServeMux, srv {{ .Name.UpperCamelCase }}Server, streamInterceptor grpc.StreamServerInterceptor) error {
	{{ range .Methods }}
	{{ template "websocket" . }}
	{{ end }}
	return nil
}

var (
	{{- range $name, $rule := (websocketURL .) }}
	{{ $name }} = runtime.MustPattern(runtime.NewPattern(1, []int{ {{ range $rule.OpCodes }}{{ . }},{{ end }} }, []string{ {{ range $rule.Pool }}"{{ . }}",{{ end }} }, "", runtime.AssumeColonVerbOpt(true)))
	{{- end }}
)
{{ end }}
`
