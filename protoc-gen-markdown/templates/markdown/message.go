package markdown

const messageTpl = `
> {{ leadingComment .SourceCodeInfo }}

|Field|Proto Type|JSON Type|Validate|Comment|{{ if .Syntax.SupportsRequiredPrefix }}Default|Required|{{ end }}
|---|---|---|---|---|{{ if .Syntax.SupportsRequiredPrefix }}---|---|{{ end }}
{{ range $field := .Fields }}{{ fieldDoc . }}
{{ end }}

> JSON Demo

#TODO
`
