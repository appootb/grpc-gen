package shared

import (
	"fmt"
	"sort"

	"github.com/appootb/protobuf/go/permission"
	"github.com/golang/protobuf/proto"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/genproto/googleapis/api/annotations"
)

type Func struct {
	pgsgo.Context
}

func (fn Func) Access(svc pgs.Service) map[string][]permission.Audience {
	out := make(map[string][]permission.Audience)
	defaultAudience := permission.Audience_NONE
	if fn.Scope(svc) == permission.VisibleScope_SERVER {
		defaultAudience = permission.Audience_SERVER
	}

	for _, method := range svc.Methods() {
		fullPath := fmt.Sprintf("/%s.%s/%s", svc.Package().ProtoName(), svc.Name(), method.Name().UpperCamelCase())
		if defaultAudience == permission.Audience_SERVER {
			// Ignore method option within server scope.
			out[fullPath] = append(out[fullPath], defaultAudience)
			continue
		}

		audiences := map[permission.Audience]int{}
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)

		for _, desc := range descs {
			if desc.TypeDescriptor().Number() == 2507 {
				ext, _ := proto.GetExtension(opts, desc)
				if auds, ok := ext.([]permission.Audience); ok {
					for _, aud := range auds {
						switch aud {
						case permission.Audience_LOGGED_IN:
							audiences[permission.Audience_WEB]++
							audiences[permission.Audience_PC]++
							audiences[permission.Audience_MOBILE]++
						case permission.Audience_CLIENT:
							audiences[permission.Audience_GUEST]++
							audiences[permission.Audience_WEB]++
							audiences[permission.Audience_PC]++
							audiences[permission.Audience_MOBILE]++
						case permission.Audience_ANY:
							audiences[permission.Audience_GUEST]++
							audiences[permission.Audience_WEB]++
							audiences[permission.Audience_PC]++
							audiences[permission.Audience_MOBILE]++
							audiences[permission.Audience_SERVER]++
						default:
							audiences[aud]++
						}
					}
				}
			}
		}

		if len(audiences) == 0 {
			audiences[defaultAudience]++
		}
		for aud := range audiences {
			out[fullPath] = append(out[fullPath], aud)
		}
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

func (fn Func) IsServerStreaming(method pgs.Method) bool {
	return method.ServerStreaming() && !method.ClientStreaming()
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
