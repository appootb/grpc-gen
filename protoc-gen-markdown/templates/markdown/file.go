package markdown

const fileTpl = `
# {{ .InputPath }} API Document
*Document generated by protoc-gen-markdown. DO NOT EDIT!*
> APIs
{{ range $svc := .Services }}
* [{{ headerTitle $svc }}](#{{ anchorName $svc }}) - {{ tocComment $svc.SourceCodeInfo }}
{{ range $method := $svc.Methods }}
{{ $url := (webUrl $method) }}
	* [{{ headerTitle $method }}{{ if $url }} ({{ $url }}){{ end }}](#{{ anchorName $method }}) - {{ tocComment $method.SourceCodeInfo }}
{{ end }}
{{ end }}
{{ range .Services }}
{{ template "service" . }}
{{ end }}
********
## *Embed Enums & Messages*
{{ range $enum := (embedEnums .) }}
{{ template "enum" $enum }}
{{ end }}
{{ range $message := (embedMessages .) }}
<h3 id="{{ anchorName $message }}">{{ headerTitle $message }}</h3>
{{ template "message" $message }}
{{ end }}
`
