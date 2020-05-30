package markdown

const methodTpl = `
<h3 id="{{ anchor .Name }}">{{ .Name.UpperCamelCase }}</h3>

> {{ leadingComment .SourceCodeInfo }}

{{ $gateway := (gatewayDoc .) }}
{{ if $gateway }}
* HTTP Gateway

	* URL: {{ $gateway.URL }}
	* Method: {{ $gateway.Method }}
{{ if $gateway.JsonRequired }}	* Content-Type: {{ $gateway.ContentType }}{{ end }}
{{ end }}
* Request Type: ***{{ .Input.Name.UpperCamelCase }}***

> {{ leadingComment .Input.SourceCodeInfo }}

|Field|proto type|JSON type|Comment|Default|Required|
|---|---|---|---|---|---|
{{ range $v := (messageDoc .Input) }}|{{ $v.Name }}|{{ $v.ProtoType }}|{{ $v.JsonType }}|{{ $v.Comment }}|{{ $v.Default }}|{{ $v.Required }}|
{{ end }}

{{ if $gateway }}
{{ if $gateway.JsonRequired }}
> JSON Demo

{{ (jsonDemo .Input) }}
{{ end }}
{{ end }}

* Response Type: ***{{ .Output.Name.UpperCamelCase }}***

> {{ leadingComment .Output.SourceCodeInfo }}

|Field|proto type|JSON type|Comment|Default|Required|
|---|---|---|---|---|---|
{{ range $v := (messageDoc .Output) }}|{{ $v.Name }}|{{ $v.ProtoType }}|{{ $v.JsonType }}|{{ $v.Comment }}|{{ $v.Default }}|{{ $v.Required }}|
{{ end }}

{{ if $gateway }}
> JSON Demo

{{ (jsonDemo .Output) }}
{{ end }}
`
