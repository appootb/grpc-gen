package golang

const clientStreamTpl = `
	mux.Handle("GET", ws_pattern_{{ .Service.Name }}_{{ .Name.UpperCamelCase }}_0, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fn := func(c *websocket.Conn) {
			ctx, err := runtime.AnnotateIncomingContext(r.Context(), mux, r, "/{{ .Package.ProtoName }}.{{ .Service.Name }}/{{ .Name.UpperCamelCase }}")
			if err != nil {
				_ = c.WriteClose(http.StatusBadRequest)
				return
			}
			inbound, outbound := runtime.MarshalerForRequest(mux, r)
			stream := gateway.NewWebsocketStream(ctx, c, inbound, outbound)
			if streamInterceptor == nil {
				err = srv.{{ .Name.UpperCamelCase }}(&{{ .Service.Name.LowerCamelCase }}{{ .Name.UpperCamelCase }}Server{stream})
				if err != nil {
					_ = c.WriteClose(http.StatusInternalServerError)
				}
				return
			}
			handler := func(_ interface{}, stream grpc.ServerStream) error {
				return srv.{{ .Name.UpperCamelCase }}(&{{ .Service.Name.LowerCamelCase }}{{ .Name.UpperCamelCase }}Server{stream})
			}
			info := &grpc.StreamServerInfo{
				FullMethod:     "/{{ .Package.ProtoName }}.{{ .Service.Name }}/{{ .Name.UpperCamelCase }}",
				IsClientStream: {{ .ClientStreaming }},
				IsServerStream: {{ .ServerStreaming }},
			}
			if err = streamInterceptor(srv, stream, info, handler); err != nil {
				_ = c.WriteClose(http.StatusInternalServerError)
			}
		}
		websocket.Handler(fn).ServeHTTP(w, r)
	})
`
