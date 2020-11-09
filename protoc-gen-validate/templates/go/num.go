package golang

const numTpl = `{{ $opt := optional .Field }}
	{{ if $opt }}if m.{{ name .Field }} != nil { {{ end }}
	{{ template "const" . }}
	{{ template "ltgt" . }}
	{{ template "in" . }}
{{ if $opt }} } {{ end }}
`
