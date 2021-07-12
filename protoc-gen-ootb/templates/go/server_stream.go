package golang

const serverStreamTpl = `
	mux.Handle("GET", ws_pattern_{{ .Service.Name }}_{{ .Name.UpperCamelCase }}_0, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		fn := func(c *websocket.Conn) {
			ctx, err := runtime.AnnotateIncomingContext(r.Context(), mux, r)
			if err != nil {
				_ = c.WriteClose(http.StatusBadRequest)
				return
			}
			inbound, outbound := runtime.MarshalerForRequest(mux, r)
			stream := webstream.NewWebsocketStream(ctx, c, inbound, outbound)
			req := new({{ inputMessage . }})
			if err = stream.RecvMsg(req); err != nil {
				_ = c.WriteClose(http.StatusBadRequest)
				return
			}
			if streamInterceptor == nil {
				err = srv.{{ .Name.UpperCamelCase }}(req, &{{ .Service.Name.LowerCamelCase }}{{ .Name.UpperCamelCase }}Server{stream})
				if err != nil {
					_ = c.WriteClose(http.StatusInternalServerError)
				}
				return
			}
			handler := func(_ interface{}, stream grpc.ServerStream) error {
				return srv.{{ .Name.UpperCamelCase }}(req, &{{ .Service.Name.LowerCamelCase }}{{ .Name.UpperCamelCase }}Server{stream})
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
