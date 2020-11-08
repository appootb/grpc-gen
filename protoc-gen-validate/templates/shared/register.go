package shared

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func Register(tpl *template.Template, params pgs.Parameters) {
	fn := Func{
		Context: pgsgo.InitContext(params),
	}

	tpl.Funcs(map[string]interface{}{
		"cmt":      pgs.C80,
		"render":   Render(tpl),
		"disabled": Disabled,
		"ignored":  Ignored,
		"required": RequiredOneOf,
		"context":  RulesContext,
		"has":      Has,
		"needs":    Needs,
		//
		"pkg":           fn.PackageName,
		"typ":           fn.Type,
		"name":          fn.Name,
		"accessor":      fn.Accessor,
		"err":           fn.Error,
		"errname":       fn.ErrorName,
		"errIdx":        fn.ErrorIdx,
		"errCause":      fn.ErrorCause,
		"errIdxCause":   fn.ErrorIdxCause,
		"lookup":        fn.Lookup,
		"lit":           fn.Lit,
		"isBytes":       fn.IsBytes,
		"byteStr":       fn.ByteStr,
		"oneof":         fn.OneOfTypeName,
		"inType":        fn.InType,
		"inKey":         fn.InKey,
		"durGt":         fn.DurationGt,
		"durLit":        fn.DurationLit,
		"durStr":        fn.DurationStr,
		"tsGt":          fn.TimestampGt,
		"tsLit":         fn.TimestampLit,
		"tsStr":         fn.TimestampStr,
		"unwrap":        fn.Unwrap,
		"msgTyp":        fn.MessageType,
		"externalEnums": fn.ExternalEnums,
		"enumPackages":  fn.EnumPackages,
		"snakeCase":     fn.SnakeCase,
	})
}
