package shared

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/iancoleman/strcase"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

type Func struct {
	pgsgo.Context
}

func (fn Func) Optional(f pgs.Field) bool {
	return f.Syntax() == pgs.Proto2 && !f.Required()
}

func (fn Func) Accessor(ctx RuleContext) string {
	if ctx.AccessorOverride != "" {
		return ctx.AccessorOverride
	}
	return fmt.Sprintf("m.Get%s()", fn.Name(ctx.Field))
}

func (fn Func) ErrorName(m pgs.Message) pgs.Name {
	return fn.Name(m) + "ValidationError"
}

func (fn Func) ErrorIdxCause(ctx RuleContext, idx, cause string, reason ...interface{}) string {
	f := ctx.Field
	n := fn.Name(f)

	var fld string
	if idx != "" {
		fld = fmt.Sprintf(`fmt.Sprintf("%s[%%v]", %s)`, n, idx)
	} else if ctx.Index != "" {
		fld = fmt.Sprintf(`fmt.Sprintf("%s[%%v]", %s)`, n, ctx.Index)
	} else {
		fld = fmt.Sprintf("%q", n)
	}

	causeFld := ""
	if cause != "nil" && cause != "" {
		causeFld = fmt.Sprintf("cause: %s,", cause)
	}

	keyFld := ""
	if ctx.OnKey {
		keyFld = "key: true,"
	}

	return fmt.Sprintf(`%s{
		field: %s,
		reason: %q,
		%s%s
	}`,
		fn.ErrorName(f.Message()),
		fld,
		fmt.Sprint(reason...),
		causeFld,
		keyFld)
}

func (fn Func) Error(ctx RuleContext, reason ...interface{}) string {
	return fn.ErrorIdxCause(ctx, "", "nil", reason...)
}

func (fn Func) ErrorCause(ctx RuleContext, cause string, reason ...interface{}) string {
	return fn.ErrorIdxCause(ctx, "", cause, reason...)
}

func (fn Func) ErrorIdx(ctx RuleContext, idx string, reason ...interface{}) string {
	return fn.ErrorIdxCause(ctx, idx, "nil", reason...)
}

func (fn Func) Lookup(f pgs.Field, name string) string {
	return fmt.Sprintf(
		"_%s_%s_%s",
		fn.Name(f.Message()),
		fn.Name(f),
		name,
	)
}

func (fn Func) Lit(x interface{}) string {
	val := reflect.ValueOf(x)

	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.String:
		return fmt.Sprintf("%q", x)
	case reflect.Uint8:
		return fmt.Sprintf("0x%X", x)
	case reflect.Slice:
		els := make([]string, val.Len())
		for i, l := 0, val.Len(); i < l; i++ {
			els[i] = fn.Lit(val.Index(i).Interface())
		}
		return fmt.Sprintf("%T{%s}", val.Interface(), strings.Join(els, ", "))
	default:
		return fmt.Sprint(x)
	}
}

func (fn Func) IsBytes(f interface {
	ProtoType() pgs.ProtoType
}) bool {
	return f.ProtoType() == pgs.BytesT
}

func (fn Func) ByteStr(x []byte) string {
	elms := make([]string, len(x))
	for i, b := range x {
		elms[i] = fmt.Sprintf(`\x%X`, b)
	}

	return fmt.Sprintf(`"%s"`, strings.Join(elms, ""))
}

func (fn Func) OneOfTypeName(f pgs.Field) pgsgo.TypeName {
	return pgsgo.TypeName(fn.OneofOption(f)).Pointer()
}

func (fn Func) InType(f pgs.Field, x interface{}) string {
	switch f.Type().ProtoType() {
	case pgs.BytesT:
		return "string"
	case pgs.MessageT:
		switch x.(type) {
		case []*duration.Duration:
			return "time.Duration"
		default:
			return pgsgo.TypeName(fmt.Sprintf("%T", x)).Element().String()
		}
	case pgs.EnumT:
		if f.Type().IsRepeated() {
			return strings.TrimLeft(fn.Type(f).Value().String(), "[]")
		} else {
			return fn.Type(f).Value().String()
		}
	default:
		return fn.Type(f).Value().String()
	}
}

func (fn Func) InKey(f pgs.Field, x interface{}) string {
	switch f.Type().ProtoType() {
	case pgs.BytesT:
		return fn.ByteStr(x.([]byte))
	case pgs.MessageT:
		switch x := x.(type) {
		case *duration.Duration:
			dur, _ := ptypes.Duration(x)
			return fn.Lit(int64(dur))
		default:
			return fn.Lit(x)
		}
	default:
		return fn.Lit(x)
	}
}

func (fn Func) DurationLit(dur *duration.Duration) string {
	return fmt.Sprintf(
		"time.Duration(%d * time.Second + %d * time.Nanosecond)",
		dur.GetSeconds(), dur.GetNanos())
}

func (fn Func) DurationStr(dur *duration.Duration) string {
	d, _ := ptypes.Duration(dur)
	return d.String()
}

func (fn Func) DurationGt(a, b *duration.Duration) bool {
	ad, _ := ptypes.Duration(a)
	bd, _ := ptypes.Duration(b)
	return ad > bd
}

func (fn Func) TimestampLit(ts *timestamp.Timestamp) string {
	return fmt.Sprintf(
		"time.Unix(%d, %d)",
		ts.GetSeconds(), ts.GetNanos(),
	)
}

func (fn Func) TimestampGt(a, b *timestamp.Timestamp) bool {
	at, _ := ptypes.Timestamp(a)
	bt, _ := ptypes.Timestamp(b)

	return bt.Before(at)
}

func (fn Func) TimestampStr(ts *timestamp.Timestamp) string {
	t, _ := ptypes.Timestamp(ts)
	return t.String()
}

func (fn Func) Unwrap(ctx RuleContext, name string) (RuleContext, error) {
	ctx, err := ctx.Unwrap("wrapper")
	if err != nil {
		return ctx, err
	}

	ctx.AccessorOverride = fmt.Sprintf("%s.Get%s()", name,
		pgsgo.PGGUpperCamelCase(ctx.Field.Type().Embed().Fields()[0].Name()))

	return ctx, nil
}

func (fn Func) MessageType(message pgs.Message) pgsgo.TypeName {
	return pgsgo.TypeName(fn.Name(message))
}

func (fn Func) ExternalEnums(file pgs.File) []pgs.Enum {
	var out []pgs.Enum

	for _, msg := range file.AllMessages() {
		for _, fld := range msg.Fields() {
			if en := fld.Type().Enum(); fld.Type().IsEnum() &&
				en.Package().ProtoName() != fld.Package().ProtoName() &&
				fn.PackageName(en) != fn.PackageName(fld) {
				out = append(out, en)
			}
		}
	}

	return out
}

func (fn Func) EnumPackages(enums []pgs.Enum) map[pgs.FilePath]pgs.Name {
	out := make(map[pgs.FilePath]pgs.Name, len(enums))

	for _, en := range enums {
		out[fn.ImportPath(en)] = fn.PackageName(en)
	}

	return out
}

func (fn Func) SnakeCase(name string) string {
	return strcase.ToSnake(name)
}
