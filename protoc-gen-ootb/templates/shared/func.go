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

type Funcs struct {
	pgsgo.Context
}

func (fns Funcs) isWebStream(method pgs.Method) bool {
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

func (fns Funcs) hasWebApi(svc pgs.Service) bool {
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

func (fns Funcs) hasWebStream(svc pgs.Service) bool {
	for _, method := range svc.Methods() {
		if fns.isWebStream(method) {
			return true
		}
	}
	return false
}

func (fns Funcs) webStreamPatterns(svc pgs.Service) map[string]httprule.Template {
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

func (fns Funcs) serviceScope(svc pgs.Service) permission.VisibleScope {
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

func (fns Funcs) methodMessageName(method pgs.Method, message pgs.Message) string {
	if method.Package().ProtoName() != message.Package().ProtoName() &&
		fns.PackageName(method) != fns.PackageName(message) {
		return fmt.Sprintf("*%s.%s", fns.PackageName(message).String(), fns.Name(message).String())
	}
	return "*" + fns.Name(message).String()
}

func (fns Funcs) inputMessage(method pgs.Method) string {
	return fns.methodMessageName(method, method.Input())
}

func (fns Funcs) outputMessage(method pgs.Method) string {
	return fns.methodMessageName(method, method.Output())
}

func (fns Funcs) subjectsName(svc pgs.Service) string {
	return fmt.Sprintf("_%sServiceSubjects", svc.Name().LowerCamelCase())
}

func (fns Funcs) serviceSubjects(svc pgs.Service) map[string][]string {
	out := make(map[string][]string)
	defaultAudience := permission.Subject_NONE
	if fns.serviceScope(svc) == permission.VisibleScope_SERVER {
		defaultAudience = permission.Subject_SERVER
	}

	for _, method := range svc.Methods() {
		fullPath := fmt.Sprintf("/%s.%s/%s", svc.Package().ProtoName(), svc.Name(), method.Name().UpperCamelCase())
		if defaultAudience == permission.Subject_SERVER {
			// Ignore method option within server scope.
			out[fullPath] = append(out[fullPath], permission.Subject_name[int32(defaultAudience)])
			continue
		}

		methodRoles := map[string]int{}
		methodAudiences := map[permission.Subject]int{}
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)

		for _, desc := range descs {
			// RBAC
			if desc.TypeDescriptor().Number() == 4507 {
				ext, _ := proto.GetExtension(opts, desc)
				if roles, ok := ext.([]string); ok {
					for _, role := range roles {
						methodRoles[role]++
					}
				}
			}
			//
			if desc.TypeDescriptor().Number() == 2507 {
				ext, _ := proto.GetExtension(opts, desc)
				if auds, ok := ext.([]permission.Subject); ok {
					for _, aud := range auds {
						switch aud {
						case permission.Subject_LOGGED_IN:
							methodAudiences[permission.Subject_WEB]++
							methodAudiences[permission.Subject_PC]++
							methodAudiences[permission.Subject_MOBILE]++
						case permission.Subject_CLIENT:
							methodAudiences[permission.Subject_GUEST]++
							methodAudiences[permission.Subject_WEB]++
							methodAudiences[permission.Subject_PC]++
							methodAudiences[permission.Subject_MOBILE]++
						case permission.Subject_ANY:
							methodAudiences[permission.Subject_GUEST]++
							methodAudiences[permission.Subject_WEB]++
							methodAudiences[permission.Subject_PC]++
							methodAudiences[permission.Subject_MOBILE]++
							methodAudiences[permission.Subject_SERVER]++
						default:
							methodAudiences[aud]++
						}
					}
				}
			}
		}

		//
		if len(methodAudiences) == 0 {
			if len(methodRoles) == 0 {
				methodAudiences[defaultAudience]++
			} else {
				methodAudiences[permission.Subject_WEB]++
				methodAudiences[permission.Subject_PC]++
				methodAudiences[permission.Subject_MOBILE]++
			}
		}
		for aud := range methodAudiences {
			out[fullPath] = append(out[fullPath], permission.Subject_name[int32(aud)])
		}
		sort.Strings(out[fullPath])
	}
	return out
}

func (fns Funcs) rolesName(svc pgs.Service) string {
	return fmt.Sprintf("_%sServiceRoles", svc.Name().LowerCamelCase())
}

func (fns Funcs) serviceRoles(svc pgs.Service) map[string][]string {
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

func (fns Funcs) externalPackages(file pgs.File) map[pgs.FilePath]pgs.Name {
	out := make(map[pgs.FilePath]pgs.Name)

	for _, service := range file.Services() {
		for _, method := range service.Methods() {
			if method.Input().Package().ProtoName() != file.Package().ProtoName() &&
				fns.PackageName(method.Input()) != fns.PackageName(file) {
				out[fns.ImportPath(method.Input())] = fns.PackageName(method.Input())
			}
			if method.Output().Package().ProtoName() != file.Package().ProtoName() &&
				fns.PackageName(method.Output()) != fns.PackageName(file) {
				out[fns.ImportPath(method.Output())] = fns.PackageName(method.Output())
			}
		}
	}
	return out
}
