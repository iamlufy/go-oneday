package odsession

import "net/http"

//manager委托给provider进行session声明周期的管理
type Provider interface {
	/**
	manager session管理器
	 */
	createSession(manager Manager) (HttpSession, error)
	//SessionRead(sid string) (HttpSession, error)
	//删除session
	sessionDestroy(session HttpSession,w http.ResponseWriter)
	//SessionGC(maxLifeTime int64)
}
//基于内存的session
type memoryProvider struct {
}

func NewMemoryProvider() *memoryProvider  {
	return new(memoryProvider)
}

func (provider *memoryProvider) createSession(manager Manager) (HttpSession, error) {
	//初始化session并将自己添加到manager中
	s := NewStandardSession(manager)
	s.lastAccessedTime = s.createTime
	s.accessTime = s.createTime
	return s, nil
}
//无效化session,包括清除属性
func (provider *memoryProvider) sessionDestroy(session HttpSession,w http.ResponseWriter) {
	session.Invalidate(w)
}
