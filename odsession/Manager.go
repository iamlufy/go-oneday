package odsession

import (
	"sync"
	"encoding/base64"
	"net/http"
	"math/rand"
)

const (
	cookieName  = "gosessionid"
	maxLifeTime = 1800
)

//session的管理器，主要作用是下达命令
type Manager interface {
	//生成sessionID方法
	SessionId() string
	//创建session，并且添加到sessionHolder中
	CreateSession(w http.ResponseWriter, r *http.Request) HttpSession
	//增加session到sessionHolder
	addSession(session HttpSession)
	MaxLifeTime() int
	//从sessionHolder获取session
	Session(sid string) HttpSession
	//删除session，从	删除session，sessionHolder，并且清空session的属性
	RemoveSession(sid string, w http.ResponseWriter)
	SessionFromHolder(sid string) (HttpSession, bool)
	IsSessionCreated(host string) bool
}

type ManagerBase struct {
	cookieName    string   // private cookieName
	sessionHolder sync.Map // 存储系统已有的session
	provider      Provider //负责session的生命周期，创建，删除等等
	maxLifeTime   int
}

var sessionProvider Provider

func NewManager(p Provider, cookieName string, maxLifeTime int) (Manager, error) {
	return &ManagerBase{provider: p, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

func Default() (Manager) {
	m, _ := NewManager(NewMemoryProvider(), cookieName, maxLifeTime)
	return m
}

//不支持集群，简单的随机而已
func (manager *ManagerBase) SessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *ManagerBase) SetProvider(p Provider) {
	manager.provider = p
}
func (manager *ManagerBase) MaxLifeTime() int {
	return manager.maxLifeTime
}

func (manager *ManagerBase) addSession(session HttpSession) {
	manager.sessionHolder.Store(session.SessionId(), session)
}

func (manager *ManagerBase) CreateSession(w http.ResponseWriter, r *http.Request) (session HttpSession) {
	session, _ = manager.provider.createSession(manager)

	cookie := http.Cookie{Name: manager.cookieName, Value: session.SessionId(), Path: "/", HttpOnly: true, MaxAge: manager.maxLifeTime}
	http.SetCookie(w, &cookie)
	return
}

func (manager *ManagerBase) RemoveSession(sid string, w http.ResponseWriter) {
	s, exist := manager.SessionFromHolder(sid)
	if exist {
		cookie := http.Cookie{Name: manager.cookieName, Value: sid, Path: "/", HttpOnly: true, MaxAge: -1}
		http.SetCookie(w, &cookie)
		manager.sessionHolder.Delete(sid)
		manager.provider.sessionDestroy(s, w)
	}

}

//从sessionHolder容器中获取
func (manager *ManagerBase) Session(sid string) HttpSession {
	result, exist := manager.SessionFromHolder(sid)
	if exist {
		return result
	}
	return nil

}

//封装从sessionHolder中取出值之后需要进行类型转换
func (manager *ManagerBase) SessionFromHolder(sid string) (HttpSession, bool) {
	s, exist := manager.sessionHolder.Load(sid)
	result, ok := s.(HttpSession)
	if ok {
		return result, exist
	}
	return nil, false
}
func (manager *ManagerBase) IsSessionCreated(host string) bool {
		return true

}
