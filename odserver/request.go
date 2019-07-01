package odserver

import "net/http"

type Request struct {
	*http.Request
	ho *HandlerObject
}
