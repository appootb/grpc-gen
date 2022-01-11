package shared

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/appootb/protobuf/go/validate"
	pgs "github.com/lyft/protoc-gen-star"
)

func (fns Funcs) valueCmp(v1, v2 reflect.Value) int {
	switch v1.Kind() {
	case reflect.Float32, reflect.Float64:
		return big.NewFloat(v1.Float()).Cmp(big.NewFloat(v2.Float()))
	case reflect.Int32, reflect.Int64:
		switch {
		case v1.Int() < v2.Int():
			return -1
		case v1.Int() > v2.Int():
			return 1
		}
	case reflect.Uint32, reflect.Uint64:
		switch {
		case v1.Uint() < v2.Uint():
			return -1
		case v1.Uint() > v2.Uint():
			return 1
		}
	}
	return 0
}

func (fns Funcs) ruleIs(v interface{}) string {
	return fmt.Sprintf("IS: `%v`", v)
}

func (fns Funcs) ruleIn(in interface{}, not bool) string {
	list := reflect.ValueOf(in)
	//
	var values []string
	for i := 0; i < list.Len(); i++ {
		values = append(values, fmt.Sprintf("%v", list.Index(i).Interface()))
	}
	v := fmt.Sprintf("IN: `[%s]`", strings.Join(values, ", "))
	if not {
		v = "NOT " + v
	}
	return v
}

func (fns Funcs) isUnsigned(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func (fns Funcs) minRule(gt, gte interface{}) (string, bool) {
	v, eqV := reflect.ValueOf(gt), reflect.ValueOf(gte)
	if v.IsNil() && eqV.IsNil() {
		if fns.isUnsigned(v.Type().Elem()) {
			return "0", true
		}
		return "-INF", false
	}
	if v.IsNil() {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	} else if eqV.IsNil() {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	}
	//
	if fns.valueCmp(v.Elem(), eqV.Elem()) > 0 {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	} else {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	}
}

func (fns Funcs) maxRule(lt, lte interface{}) (string, bool) {
	v, eqV := reflect.ValueOf(lt), reflect.ValueOf(lte)
	if v.IsNil() && eqV.IsNil() {
		return "+INF", false
	}
	if v.IsNil() {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	} else if eqV.IsNil() {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	}
	//
	if fns.valueCmp(v.Elem(), eqV.Elem()) > 0 {
		return fmt.Sprintf("%v", eqV.Elem().Interface()), true
	} else {
		return fmt.Sprintf("%v", v.Elem().Interface()), false
	}
}

func (fns Funcs) ruleInterval(min, max string, minEq, maxEq bool) string {
	if minEq {
		min = "[" + min
	} else {
		min = "(" + min
	}
	if maxEq {
		max = max + "]"
	} else {
		max = max + ")"
	}
	return fmt.Sprintf("RANGE: `%s, %s`", min, max)
}

func (fns Funcs) floatRules(f *validate.FloatRules) string {
	var rules []string
	//
	if f.Gt != nil || f.Gte != nil || f.Lt != nil || f.Lte != nil {
		min, minEq := fns.minRule(f.Gt, f.Gte)
		max, maxEq := fns.maxRule(f.Lt, f.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if f.Const != nil {
		rules = append(rules, fns.ruleIs(f.GetConst()))
	}
	if len(f.In) > 0 {
		rules = append(rules, fns.ruleIn(f.In, false))
	}
	if len(f.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(f.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) doubleRules(d *validate.DoubleRules) string {
	var rules []string
	//
	if d.Gt != nil || d.Gte != nil || d.Lt != nil || d.Lte != nil {
		min, minEq := fns.minRule(d.Gt, d.Gte)
		max, maxEq := fns.maxRule(d.Lt, d.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if d.Const != nil {
		rules = append(rules, fns.ruleIs(d.GetConst()))
	}
	if len(d.In) > 0 {
		rules = append(rules, fns.ruleIn(d.In, false))
	}
	if len(d.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(d.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) int32Rules(i *validate.Int32Rules) string {
	var rules []string
	//
	if i.Gt != nil || i.Gte != nil || i.Lt != nil || i.Lte != nil {
		min, minEq := fns.minRule(i.Gt, i.Gte)
		max, maxEq := fns.maxRule(i.Lt, i.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if i.Const != nil {
		rules = append(rules, fns.ruleIs(i.GetConst()))
	}
	if len(i.In) > 0 {
		rules = append(rules, fns.ruleIn(i.In, false))
	}
	if len(i.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(i.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) int64Rules(i *validate.Int64Rules) string {
	var rules []string
	//
	if i.Gt != nil || i.Gte != nil || i.Lt != nil || i.Lte != nil {
		min, minEq := fns.minRule(i.Gt, i.Gte)
		max, maxEq := fns.maxRule(i.Lt, i.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if i.Const != nil {
		rules = append(rules, fns.ruleIs(i.GetConst()))
	}
	if len(i.In) > 0 {
		rules = append(rules, fns.ruleIn(i.In, false))
	}
	if len(i.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(i.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) uint32Rules(u *validate.UInt32Rules) string {
	var rules []string
	//
	if u.Gt != nil || u.Gte != nil || u.Lt != nil || u.Lte != nil {
		min, minEq := fns.minRule(u.Gt, u.Gte)
		max, maxEq := fns.maxRule(u.Lt, u.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if u.Const != nil {
		rules = append(rules, fns.ruleIs(u.GetConst()))
	}
	if len(u.In) > 0 {
		rules = append(rules, fns.ruleIn(u.In, false))
	}
	if len(u.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(u.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) uint64Rules(u *validate.UInt64Rules) string {
	var rules []string
	//
	if u.Gt != nil || u.Gte != nil || u.Lt != nil || u.Lte != nil {
		min, minEq := fns.minRule(u.Gt, u.Gte)
		max, maxEq := fns.maxRule(u.Lt, u.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if u.Const != nil {
		rules = append(rules, fns.ruleIs(u.GetConst()))
	}
	if len(u.In) > 0 {
		rules = append(rules, fns.ruleIn(u.In, false))
	}
	if len(u.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(u.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sInt32Rules(s *validate.SInt32Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		min, minEq := fns.minRule(s.Gt, s.Gte)
		max, maxEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sInt64Rules(s *validate.SInt64Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		min, minEq := fns.minRule(s.Gt, s.Gte)
		max, maxEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) fixed32Rules(f *validate.Fixed32Rules) string {
	var rules []string
	//
	if f.Gt != nil || f.Gte != nil || f.Lt != nil || f.Lte != nil {
		min, minEq := fns.minRule(f.Gt, f.Gte)
		max, maxEq := fns.maxRule(f.Lt, f.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if f.Const != nil {
		rules = append(rules, fns.ruleIs(f.GetConst()))
	}
	if len(f.In) > 0 {
		rules = append(rules, fns.ruleIn(f.In, false))
	}
	if len(f.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(f.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) fixed64Rules(f *validate.Fixed64Rules) string {
	var rules []string
	//
	if f.Gt != nil || f.Gte != nil || f.Lt != nil || f.Lte != nil {
		min, minEq := fns.minRule(f.Gt, f.Gte)
		max, maxEq := fns.maxRule(f.Lt, f.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if f.Const != nil {
		rules = append(rules, fns.ruleIs(f.GetConst()))
	}
	if len(f.In) > 0 {
		rules = append(rules, fns.ruleIn(f.In, false))
	}
	if len(f.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(f.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sFixed32Rules(s *validate.SFixed32Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		min, minEq := fns.minRule(s.Gt, s.Gte)
		max, maxEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) sFixed64Rules(s *validate.SFixed64Rules) string {
	var rules []string
	//
	if s.Gt != nil || s.Gte != nil || s.Lt != nil || s.Lte != nil {
		min, minEq := fns.minRule(s.Gt, s.Gte)
		max, maxEq := fns.maxRule(s.Lt, s.Lte)
		rules = append(rules, fns.ruleInterval(min, max, minEq, maxEq))
	}
	//
	if s.Const != nil {
		rules = append(rules, fns.ruleIs(s.GetConst()))
	}
	if len(s.In) > 0 {
		rules = append(rules, fns.ruleIn(s.In, false))
	}
	if len(s.NotIn) > 0 {
		rules = append(rules, fns.ruleIn(s.NotIn, true))
	}
	if len(rules) == 0 {
		return "-"
	}
	return strings.Join(rules, "<br>")
}

func (fns Funcs) boolRules(b *validate.BoolRules) string {
	if b.Const != nil {
		return fns.ruleIs(b.GetConst())
	}
	return "-"
}

func (fns Funcs) stringRules(s *validate.StringRules) string {
	return "-"
}

func (fns Funcs) bytesRules(b *validate.BytesRules) string {
	return "-"
}

func (fns Funcs) enumRules(enum *validate.EnumRules) string {
	return "-"
}

func (fns Funcs) repeatedRules(r *validate.RepeatedRules) string {
	return "-"
}

func (fns Funcs) mapRules(m *validate.MapRules) string {
	return "-"
}

func (fns Funcs) anyRules(any *validate.AnyRules) string {
	return "-"
}

func (fns Funcs) durationRules(dur *validate.DurationRules) string {
	return "-"
}

func (fns Funcs) timestampRules(ts *validate.TimestampRules) string {
	return "-"
}

func (fns Funcs) fieldRules(field pgs.Field) string {
	var rules validate.FieldRules
	if ok, _ := field.Extension(validate.E_Rules, &rules); !ok {
		return "-"
	}

	switch r := rules.GetType().(type) {
	case *validate.FieldRules_Float:
		return fns.floatRules(r.Float)
	case *validate.FieldRules_Double:
		return fns.doubleRules(r.Double)
	case *validate.FieldRules_Int32:
		return fns.int32Rules(r.Int32)
	case *validate.FieldRules_Int64:
		return fns.int64Rules(r.Int64)
	case *validate.FieldRules_Uint32:
		return fns.uint32Rules(r.Uint32)
	case *validate.FieldRules_Uint64:
		return fns.uint64Rules(r.Uint64)
	case *validate.FieldRules_Sint32:
		return fns.sInt32Rules(r.Sint32)
	case *validate.FieldRules_Sint64:
		return fns.sInt64Rules(r.Sint64)
	case *validate.FieldRules_Fixed32:
		return fns.fixed32Rules(r.Fixed32)
	case *validate.FieldRules_Fixed64:
		return fns.fixed64Rules(r.Fixed64)
	case *validate.FieldRules_Sfixed32:
		return fns.sFixed32Rules(r.Sfixed32)
	case *validate.FieldRules_Sfixed64:
		return fns.sFixed64Rules(r.Sfixed64)
	case *validate.FieldRules_Bool:
		return fns.boolRules(r.Bool)
	case *validate.FieldRules_String_:
		return fns.stringRules(r.String_)
	case *validate.FieldRules_Bytes:
		return fns.bytesRules(r.Bytes)
	case *validate.FieldRules_Enum:
		return fns.enumRules(r.Enum)
	case *validate.FieldRules_Repeated:
		return fns.repeatedRules(r.Repeated)
	case *validate.FieldRules_Map:
		return fns.mapRules(r.Map)
	case *validate.FieldRules_Any:
		return fns.anyRules(r.Any)
	case *validate.FieldRules_Duration:
		return fns.durationRules(r.Duration)
	case *validate.FieldRules_Timestamp:
		return fns.timestampRules(r.Timestamp)
	//case nil: // TODO
	default:
		return "-"
	}
}
