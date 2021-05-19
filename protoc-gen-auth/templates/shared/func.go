package shared

import (
	"fmt"
	"sort"

	"github.com/appootb/protobuf/go/api"
	"github.com/appootb/protobuf/go/permission"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/httprule"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/genproto/googleapis/api/annotations"
)

type Func struct {
	pgsgo.Context
}

func (fn Func) ServiceRoles(svc pgs.Service) map[string][]string {
	out := make(map[string][]string)
	for _, method := range svc.Methods() {
		fullPath := fmt.Sprintf("/%s.%s/%s", svc.Package().ProtoName(), svc.Name(), method.Name().UpperCamelCase())
		urlRoles := map[string]int{}
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)

		for _, desc := range descs {
			if desc.TypeDescriptor().Number() == 4507 {
				ext, _ := proto.GetExtension(opts, desc)
				if roles, ok := ext.([]string); ok {
					for _, role := range roles {
						urlRoles[role]++
					}
				}
			}
		}

		for role := range urlRoles {
			out[fullPath] = append(out[fullPath], role)
		}
		sort.Strings(out[fullPath])
	}

	return out
}

func (fn Func) Access(svc pgs.Service) map[string]Subjects {
	out := make(map[string]Subjects)
	defaultAudience := permission.Subject_NONE
	if fn.Scope(svc) == permission.VisibleScope_SERVER {
		defaultAudience = permission.Subject_SERVER
	}

	for _, method := range svc.Methods() {
		fullPath := fmt.Sprintf("/%s.%s/%s", svc.Package().ProtoName(), svc.Name(), method.Name().UpperCamelCase())
		if defaultAudience == permission.Subject_SERVER {
			// Ignore method option within server scope.
			out[fullPath] = append(out[fullPath], defaultAudience)
			continue
		}

		audiences := map[permission.Subject]int{}
		urlRoles := map[string]int{}
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)

		for _, desc := range descs {
			if desc.TypeDescriptor().Number() == 4507 {
				ext, _ := proto.GetExtension(opts, desc)
				if roles, ok := ext.([]string); ok {
					for _, role := range roles {
						urlRoles[role]++
					}
				}
			}
			if desc.TypeDescriptor().Number() == 2507 {
				ext, _ := proto.GetExtension(opts, desc)
				if auds, ok := ext.([]permission.Subject); ok {
					for _, aud := range auds {
						switch aud {
						case permission.Subject_LOGGED_IN:
							audiences[permission.Subject_WEB]++
							audiences[permission.Subject_PC]++
							audiences[permission.Subject_MOBILE]++
						case permission.Subject_CLIENT:
							audiences[permission.Subject_GUEST]++
							audiences[permission.Subject_WEB]++
							audiences[permission.Subject_PC]++
							audiences[permission.Subject_MOBILE]++
						case permission.Subject_ANY:
							audiences[permission.Subject_GUEST]++
							audiences[permission.Subject_WEB]++
							audiences[permission.Subject_PC]++
							audiences[permission.Subject_MOBILE]++
							audiences[permission.Subject_SERVER]++
						default:
							audiences[aud]++
						}
					}
				}
			}
		}

		if len(audiences) == 0 {
			if len(urlRoles) == 0 {
				audiences[defaultAudience]++
			} else {
				audiences[permission.Subject_WEB]++
				audiences[permission.Subject_PC]++
				audiences[permission.Subject_MOBILE]++
			}
		}
		for aud := range audiences {
			out[fullPath] = append(out[fullPath], aud)
		}
		sort.Sort(out[fullPath])
	}

	return out
}

func (fn Func) Scope(svc pgs.Service) permission.VisibleScope {
	opts := svc.Descriptor().GetOptions()
	descs, _ := proto.ExtensionDescs(opts)

	for _, desc := range descs {
		if desc.TypeDescriptor().Number() == 1507 {
			ext, _ := proto.GetExtension(opts, desc)
			if scope, ok := ext.(*permission.VisibleScope); ok {
				return *scope
			}
		}
	}

	return permission.VisibleScope_CLIENT
}

func (fn Func) GatewayDefined(svc pgs.Service) bool {
	for _, method := range svc.Methods() {
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)

		for _, desc := range descs {
			if desc.TypeDescriptor().Number() == 72295728 {
				ext, _ := proto.GetExtension(opts, desc)
				if _, ok := ext.(*annotations.HttpRule); ok {
					return true
				}
			}
		}
	}

	return false
}

func (fn Func) WebsocketURLPattern(svc pgs.Service) map[string]httprule.Template {
	patterns := make(map[string]httprule.Template)
	for _, method := range svc.Methods() {
		key := fmt.Sprintf("ws_pattern_%s_%s_0", svc.Name(), method.Name().UpperCamelCase())
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)
		for _, desc := range descs {
			if desc.TypeDescriptor().Number() == 3507 {
				ext, _ := proto.GetExtension(opts, desc)
				if rule, ok := ext.(*api.WebsocketRule); ok {
					c, err := httprule.Parse(rule.Url)
					if err != nil {
						continue
					}
					patterns[key] = c.Compile()
				}
			}
		}
	}
	return patterns
}

func (fn Func) WebsocketDefined(svc pgs.Service) bool {
	for _, method := range svc.Methods() {
		if fn.IsWebsocket(method) {
			return true
		}
	}
	return false
}

func (fn Func) IsWebsocket(method pgs.Method) bool {
	opts := method.Descriptor().GetOptions()
	descs, _ := proto.ExtensionDescs(opts)
	for _, desc := range descs {
		if desc.TypeDescriptor().Number() == 3507 {
			ext, _ := proto.GetExtension(opts, desc)
			if _, ok := ext.(*api.WebsocketRule); ok {
				return true
			}
		}
	}
	return false
}

func (fn Func) GolangInputMessageName(method pgs.Method) string {
	messageName := method.Input().Name().UpperCamelCase().String()
	if method.Input().Package() == method.Package() {
		return messageName
	}
	return fn.PackageName(method.Input()).String() + "." + messageName
}

func (fn Func) GolangOutputMessageName(method pgs.Method) string {
	messageName := method.Output().Name().UpperCamelCase().String()
	if method.Output().Package() == method.Package() {
		return messageName
	}
	return fn.PackageName(method.Output()).String() + "." + messageName
}

func (fn Func) GolangImports(file pgs.File) []string {
	imps := make(map[pgs.FilePath]int)
	for _, service := range file.Services() {
		for _, method := range service.Methods() {
			if method.Input().Package() != method.Package() {
				imps[fn.ImportPath(method.Input())]++
			}
			if method.Output().Package() != method.Package() {
				imps[fn.ImportPath(method.Output())]++
			}
		}
	}
	files := make([]string, 0, len(imps))
	for f := range imps {
		files = append(files, string(f))
	}
	sort.Strings(files)
	return files
}
