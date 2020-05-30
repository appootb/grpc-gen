package shared

import (
	"fmt"

	"github.com/appootb/protobuf/go/permission"
	"github.com/golang/protobuf/proto"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/genproto/googleapis/api/annotations"
)

type Func struct {
	pgsgo.Context
}

func (fn Func) Access(svc pgs.Service) map[string]permission.TokenLevel {
	out := make(map[string]permission.TokenLevel)
	defaultTokenLevel := permission.TokenLevel_NONE_TOKEN
	if fn.Scope(svc) == permission.VisibleScope_INNER_SCOPE {
		defaultTokenLevel = permission.TokenLevel_INNER_TOKEN
	}

	for _, method := range svc.Methods() {
		fullPath := fmt.Sprintf("/%s.%s/%s", svc.Package().ProtoName(), svc.Name(), method.Name().UpperCamelCase())
		out[fullPath] = defaultTokenLevel

		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)

		for _, desc := range descs {
			if desc.Field == 2507 {
				ext, _ := proto.GetExtension(opts, desc)
				if access, ok := ext.(*permission.Token); ok {
					out[fullPath] = access.Required
					break
				}
			}
		}
	}

	return out
}

func (fn Func) Scope(svc pgs.Service) permission.VisibleScope {
	opts := svc.Descriptor().GetOptions()
	descs, _ := proto.ExtensionDescs(opts)

	for _, desc := range descs {
		if desc.Field == 1507 {
			ext, _ := proto.GetExtension(opts, desc)
			if visible, ok := ext.(*permission.ServiceVisible); ok {
				return visible.Scope
			}
		}
	}

	return permission.VisibleScope_DEFAULT_SCOPE
}

func (fn Func) GatewayDefined(svc pgs.Service) bool {
	for _, method := range svc.Methods() {
		opts := method.Descriptor().GetOptions()
		descs, _ := proto.ExtensionDescs(opts)

		for _, desc := range descs {
			if desc.Field == 72295728 {
				ext, _ := proto.GetExtension(opts, desc)
				if _, ok := ext.(*annotations.HttpRule); ok {
					return true
				}
			}
		}
	}

	return false
}
