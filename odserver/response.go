package odserver

import "net/http"

type responseWriter struct {
	 http.ResponseWriter
}
func NewResponse(rw http.ResponseWriter) responseWriter {
	return responseWriter{rw,}
}
