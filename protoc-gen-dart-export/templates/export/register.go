package export

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
)

func Register(tpl *template.Template, _ pgs.Parameters) {
	template.Must(tpl.Parse(fileTpl))
}
