package shared

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/appootb/substratum/proto/go/api"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/genproto/googleapis/api/annotations"
)

const (
	TODOComment = "***TODO***"
)

type Funcs struct {
	pgsgo.Context
}

func (fns Funcs) anchorName(name pgs.Name) string {
	return name.Transform(strings.ToLower, strings.ToLower, "").String()
}

func (fns Funcs) docFileName(file pgs.File) string {
	name := path.Base(file.Name().String())
	extName := path.Ext(name)
	return strings.ReplaceAll(name, extName, ".pb.md")
}

func (fns Funcs) tocComment(srcInfo pgs.SourceCodeInfo) string {
	comment := srcInfo.LeadingComments()
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return TODOComment
}

func (fns Funcs) leadingComment(srcInfo pgs.SourceCodeInfo) string {
	comment := TODOComment
	if srcInfo.LeadingComments() != "" {
		comment = srcInfo.LeadingComments()
	}
	comment = strings.Trim(comment, "\n")
	return strings.Replace(comment, "\n", "\n> ", -1)
}

func (fns Funcs) trailingComment(srcInfo pgs.SourceCodeInfo) string {
	comment := TODOComment
	if srcInfo.TrailingComments() != "" {
		comment = srcInfo.TrailingComments()
	} else if srcInfo.LeadingComments() != "" {
		comment = srcInfo.LeadingComments()
	}
	comment = strings.Trim(comment, "\n")
	return strings.Replace(comment, "\n", "<br>", -1)
}

func (fns Funcs) webURL(method pgs.Method) string {
	var (
		webStreamRule *api.WebsocketRule
		webApiRules   *annotations.HttpRule
	)

	if ok, _ := method.Extension(api.E_Websocket, &webStreamRule); ok {
		return webStreamRule.Url
	}
	if ok, _ := method.Extension(annotations.E_Http, &webApiRules); ok {
		switch p := webApiRules.Pattern.(type) {
		case *annotations.HttpRule_Get:
			return p.Get
		case *annotations.HttpRule_Put:
			return p.Put
		case *annotations.HttpRule_Post:
			return p.Post
		case *annotations.HttpRule_Delete:
			return p.Delete
		case *annotations.HttpRule_Patch:
			return p.Patch
		case *annotations.HttpRule_Custom:
			return p.Custom.Path
		}
	}
	return ""
}

type WebDoc struct {
	URL         string
	Method      string
	ContentType string
}

func (fns Funcs) webDoc(method pgs.Method) *WebDoc {
	var (
		webStreamRule *api.WebsocketRule
		webApiRules   *annotations.HttpRule
	)

	if ok, _ := method.Extension(api.E_Websocket, &webStreamRule); ok {
		return &WebDoc{
			URL:         webStreamRule.Url,
			Method:      "`WS/WSS`",
			ContentType: "`TextFrame`",
		}
	}
	if ok, _ := method.Extension(annotations.E_Http, &webApiRules); ok {
		switch p := webApiRules.Pattern.(type) {
		case *annotations.HttpRule_Get:
			return &WebDoc{
				URL:    fmt.Sprintf("`%s`", p.Get),
				Method: fmt.Sprintf("`%s`", http.MethodGet),
			}
		case *annotations.HttpRule_Put:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Put),
				Method:      fmt.Sprintf("`%s`", http.MethodPut),
				ContentType: "`application/json`",
			}
		case *annotations.HttpRule_Post:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Post),
				Method:      fmt.Sprintf("`%s`", http.MethodPost),
				ContentType: "`application/json`",
			}
		case *annotations.HttpRule_Delete:
			return &WebDoc{
				URL:    fmt.Sprintf("`%s`", p.Delete),
				Method: fmt.Sprintf("`%s`", http.MethodDelete),
			}
		case *annotations.HttpRule_Patch:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Patch),
				Method:      fmt.Sprintf("`%s`", http.MethodPatch),
				ContentType: "`application/json`",
			}
		case *annotations.HttpRule_Custom:
			return &WebDoc{
				URL:         fmt.Sprintf("`%s`", p.Custom.Path),
				Method:      fmt.Sprintf("`%s`", p.Custom.Kind),
				ContentType: "`application/json`",
			}
		}
	}
	return nil
}

func (fns Funcs) enumName(enum pgs.Enum, pkg pgs.Package) string {
	if pkg.ProtoName() != enum.Package().ProtoName() &&
		fns.PackageName(pkg) != fns.PackageName(enum) {
		return fmt.Sprintf("%s.%s", fns.PackageName(enum).String(), fns.Name(enum).String())
	}
	return fns.Name(enum).String()
}

func (fns Funcs) messageName(message pgs.Message, pkg pgs.Package) string {
	if pkg.ProtoName() != message.Package().ProtoName() &&
		fns.PackageName(pkg) != fns.PackageName(message) {
		return fmt.Sprintf("%s.%s", fns.PackageName(message).String(), fns.Name(message).String())
	}
	return fns.Name(message).String()
}

func (fns Funcs) inputMessage(method pgs.Method) string {
	return fns.messageName(method.Input(), method.Package())
}

func (fns Funcs) outputMessage(method pgs.Method) string {
	return fns.messageName(method.Output(), method.Package())
}

func (fns Funcs) embedEnums(file pgs.File) []pgs.Enum {
	enums := map[string]pgs.Enum{}
	for _, enum := range file.Enums() {
		enums[enum.FullyQualifiedName()] = enum
	}
	//
	var messages []pgs.Message
	for _, svc := range file.Services() {
		for _, method := range svc.Methods() {
			messages = append(messages, method.Input(), method.Output())
			for _, msg := range method.Input().Dependents() {
				messages = append(messages, msg)
			}
			for _, msg := range method.Output().Dependents() {
				messages = append(messages, msg)
			}
		}
	}
	for _, msg := range messages {
		for _, field := range msg.Fields() {
			if field.Type().ProtoType() != pgs.EnumT {
				continue
			}
			var enum pgs.Enum
			if field.Type().IsRepeated() {
				enum = field.Type().Element().Enum()
			} else {
				enum = field.Type().Enum()
			}
			enums[enum.FullyQualifiedName()] = enum
		}
	}
	//
	resp := make([]pgs.Enum, 0, len(enums))
	for _, enum := range enums {
		resp = append(resp, enum)
	}
	sort.Slice(resp, func(i, j int) bool {
		return strings.Compare(resp[i].Name().String(), resp[j].Name().String()) < 0
	})
	return resp
}

func (fns Funcs) embedMessages(file pgs.File) []pgs.Message {
	messages := map[string]pgs.Message{}
	for _, msg := range file.Messages() {
		messages[msg.FullyQualifiedName()] = msg
		for _, dep := range msg.Dependents() {
			messages[dep.FullyQualifiedName()] = dep
		}
	}
	//
	for _, svc := range file.Services() {
		for _, method := range svc.Methods() {
			messages[method.Input().FullyQualifiedName()] = method.Input()
			messages[method.Output().FullyQualifiedName()] = method.Output()
			//
			for _, msg := range method.Input().Dependents() {
				messages[msg.FullyQualifiedName()] = msg
			}
			for _, msg := range method.Output().Dependents() {
				messages[msg.FullyQualifiedName()] = msg
			}
		}
	}
	//
	resp := make([]pgs.Message, 0, len(messages))
	for _, msg := range messages {
		if !msg.IsWellKnown() {
			resp = append(resp, msg)
		}
	}
	sort.Slice(resp, func(i, j int) bool {
		return strings.Compare(resp[i].Name().String(), resp[j].Name().String()) < 0
	})
	return resp
}
