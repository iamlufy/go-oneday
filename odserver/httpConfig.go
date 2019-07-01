package odserver

const (
	CONTENT_TYPE    = "Content-Type"
	CONTENT_LENGTH  = "Content-Length"
	CONTENT_BINARY  = "application/octet-stream"
	CONTENT_JSON    = "application/json"
	CONTENT_HTML    = "text/html"
	CONTENT_PLAIN   = "text/plain"
	CONTENT_XHTML   = "application/xhtml+xml"
	CONTENT_XML     = "text/xml"
	DEFAULT_CHARSET = "UTF-8"
)

type headerMaps map[string]string

type httpConfig struct {
	header headerMaps
}

func NewHttpConfig() *httpConfig {
	return &httpConfig{header: make(map[string]string),}
}
func DefaultConfig()*httpConfig {
	config := &httpConfig{header: make(map[string]string),}
	config.ProduceJSON()
	return config
}

//设置全局http 头部信息，优先级比单个设置低
func (h *httpConfig) SetHeader(key, value string) *httpConfig {
	h.header[key] = value
	return h
}

//设置全局http 头部信息，优先级比单个设置低
func (h *httpConfig) RemoveHeader(key string) *httpConfig {
	delete(h.header, key)
	return h
}

func (h *httpConfig) Header() (headerMaps) {
	return h.header
}
func (h *httpConfig) Produces(s string) *httpConfig{
	h.header[CONTENT_TYPE] = s
	return h
}
func (h *httpConfig) ProduceJSON() *httpConfig {
	h.Produces(CONTENT_JSON)
	return h
}

//func (h *httpConfig)Write(p []byte) (n int, err error){
//	return nil
//}
