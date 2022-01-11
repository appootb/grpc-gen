package readme

const fileTpl = `
# APIs
*Document generated by protoc-gen-markdown. DO NOT EDIT!*
{{ range $file := . }}
{{ range $svc := $file.Services }}
## [{{ $svc.Name.UpperCamelCase }}]({{ docFileName $file }}#{{ anchorName $svc }}) - {{ tocComment $svc.SourceCodeInfo }}
> {{ leadingComment $svc.SourceCodeInfo }}
{{ range $method := $svc.Methods }}
{{ $url := (webUrl $method) }}
* [{{ $method.Name.UpperCamelCase }}{{ if $url }} ({{ $url }}){{ end }}]({{ docFileName $file }}#{{ anchorName $method }}) - {{ tocComment $method.SourceCodeInfo }}
{{ end }}
{{ end }}
{{ end }}
`
