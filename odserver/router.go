package odserver

import (
	"strings"
	"regexp"
	"time"
	"net/http"
)

/**
提供基本的路由功能，添加路由，查找路由
 */
const (
	GET         = iota
	POST
	PUT
	DELETE
	CONNECTIBNG
	HEAD
	OPTIONS
	PATCH
	TRACE
	paramsSize  = 6
)


type handler map[string]*HandlerObject
type regexpMap map[*regexp.Regexp]*HandlerObject

func NewRouter() *Router {
	return &Router{
		handler:   make(map[string]*HandlerObject),
		regexpMap: make(map[*regexp.Regexp]*HandlerObject),
	}
}

//接口函数单位，即我们编写代码逻辑的函数
type IHandlerFunc interface{}

type FuncObject struct {
	params []string
	//对应编写的接口，IHandlerFunc是个空接口
	f     IHandlerFunc
	*httpConfig
}


func NewFuncObject(f IHandlerFunc) FuncObject {
	return FuncObject{
		f:          f,
		httpConfig: &httpConfig{header: make(map[string]string)},
	}
}

//方法函数映射，0代表GET方法下的接口
type methodFuncs map[int]FuncObject

/**
	关键struct，代表每个实体的请求
 */
type HandlerObject struct {
	//方法函数映射
	methodFuncs methodFuncs
	*Router
	//对应占位符的参数
	params []string

	//请求路径 即start+target的路径
	path      string
	startPath string

}



type IHandler interface {
	Get(f IHandlerFunc) *HandlerObject
	Post(f IHandlerFunc) *HandlerObject
	Put(f IHandlerFunc) *HandlerObject
	Delete(f IHandlerFunc) *HandlerObject
}

func NewHandlerObject(r *Router, startPath string) *HandlerObject {
	return &HandlerObject{
		params:      make([]string, paramsSize),
		Router:      r,
		startPath:   startPath,
		methodFuncs: make(map[int]FuncObject),
	}
}

func NewRequest(r *http.Request) Request {
	return Request{Request: r}
}

func (ho *HandlerObject) Get(f IHandlerFunc) *HandlerObject {
	_,exist:=ho.methodFuncs[GET]
	if exist {
		panic("GetFunc has existed")
	}
	ho.methodFuncs[GET] = NewFuncObject(f)
	return ho
}

func (ho *HandlerObject) Post(f IHandlerFunc) *HandlerObject {
	_,exist:=ho.methodFuncs[POST]
	if exist {
		panic("GetFunc has existed")
	}
	ho.methodFuncs[POST] = NewFuncObject(f)
	return ho
}

func (ho *HandlerObject) Put(f IHandlerFunc) *HandlerObject {
	_,exist:=ho.methodFuncs[PUT]
	if exist {
		panic("GetFunc has existed")
	}
	ho.methodFuncs[PUT] = NewFuncObject(f)
	return ho
}
func (ho *HandlerObject) Delete(f IHandlerFunc) *HandlerObject {
	_,exist:=ho.methodFuncs[DELETE]
	if exist {
		panic("GetFunc has existed")
	}
	ho.methodFuncs[DELETE] = NewFuncObject(f)
	return ho
}

func (ho *HandlerObject) Func(method int) (FuncObject, bool) {
	switch method {
	case GET:
		return ho.getFunc()
	case DELETE:
		return ho.deleteFunc()
	case PUT:
		return ho.putFunc()
	case POST:
		return ho.postFunc()
	case TRACE:
		return ho.traceFunc()
	case PATCH:
		return ho.patchFunc()
	case OPTIONS:
		return ho.optionsFunc()
	case HEAD:
		return ho.headFunc()
	case CONNECTIBNG:
		return ho.connectingFunc()
	}
	return FuncObject{}, false

}

func (ho *HandlerObject) getFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[GET]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}
func (ho *HandlerObject) postFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[POST]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}
func (ho *HandlerObject) putFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[PUT]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}

func (ho *HandlerObject) deleteFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[DELETE]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}

func (ho *HandlerObject) connectingFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[CONNECTIBNG]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}
func (ho *HandlerObject) headFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[HEAD]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}

func (ho *HandlerObject) optionsFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[OPTIONS]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}

func (ho *HandlerObject) patchFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[PATCH]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}

func (ho *HandlerObject) traceFunc() (FuncObject, bool) {
	handler,exist:=ho.methodFuncs[TRACE]
	if exist {
		handler.params = ho.params
	}
	return handler, exist
}

type Router struct {
	handler
	regexpMap
}

func (r *Router) Start(url string) *HandlerObject {
	return NewHandlerObject(r, addSlash(url))
}

func (ho *HandlerObject) And() *HandlerObject {
	if ho.Router == nil || ho.startPath == "" {
		panic("ho.Router is nil or startPath is unknown，maybe u should use Start()")
	}
	return NewHandlerObject(ho.Router, ho.startPath)
}

func (ho *HandlerObject) Target(url string) *HandlerObject {
	//设置完整的路径
	if ho.startPath == "/" {
		ho.path = ho.startPath + deleteSlash(url)
	} else {
		if strings.HasSuffix(ho.startPath, "/") {
			url = deleteSlash(url)
		} else {
			url = addSlash(url)
		}
		ho.path = ho.startPath + url
	}
	//尝试将url转换成正则表达式，如果没有占位符，则转换不成功
	pattern, ok := matcher.ToPattern(ho.path)
	if ok {
		ho.path = pattern
		re, err := regexp.Compile(pattern)
		if err != nil {
			panic("error compile pattern:" + pattern)
		}
		ho.Router.regexpMap[re] = ho
	} else {
		ho.handler[ho.path] = ho
	}
	return ho
}
func addSlash(s string) string {
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	return s
}

func deleteSlash(s string) string {
	if strings.HasPrefix(s, "/") {
		array := strings.SplitN(s, "/", 2)
		s = array[1]
	}
	return s
}

//匹配路径
func (r *Router) doUrlMapping(url string, method int) (FuncObject, bool) {
	ch := make(chan *HandlerObject)
	//精准匹配
	go func() {
		if ho, ok := r.handler[url]; ok {
			ch <- ho
		}
	}()
	//正则匹配
	go func() {
		for k, v := range r.regexpMap {
			if k.MatchString(url) {
				pathArray := strings.Split(url, "/")[1:]
				regexpArray := strings.Split(k.String(), "/")[1:]
				if len(pathArray) == len(regexpArray) {
					//设置参数
					paramsNum := 0
					for i := 0; i < len(pathArray); i++ {
						if matcher.IsPattern(regexpArray[i]) {
							v.params[paramsNum] = pathArray[i]
							paramsNum++
						}
					}
					v.params = v.params[:paramsNum]
				}
				ch <- v
			}
		}
	}()
	select {
	case ho := <-ch:
		{
			return ho.Func(method)
		}
	case <-time.After(2e6):
		{
			return FuncObject{}, false
		}
	}

}
