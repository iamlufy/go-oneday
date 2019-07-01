package odsession

import (
	"time"
	"sync"
	"sync/atomic"
	"net/http"
)

type HttpSession interface {
	//全局ID
	SessionId() string
	//上次访问时间
	LastAccessedTime() time.Time
	//设置最大访问间隔，超过该时间则session失效
	SetMaxInactiveInterval(d time.Duration)
	////获取最大访问间隔
	MaxInactiveInterval() time.Duration
	//获取
	Load(key string) (value interface{}, ok bool)
	//存储
	Store(key string, value interface{})
	//删除
	Remove(key string)
	/*
 	* 无效化该session ，并且解除所有跟该session的绑定关系
 	*/
	Invalidate(w http.ResponseWriter)
	/**
	 * 过期操作
	 */
	Expire(w http.ResponseWriter)
	//是否是新创建的
	//isNew() bool
	//除去所有的属性
	RemoveAll()
	/**
	* Update the accessed time information for this session.  This method
	* should be called by the context when a request comes in for a particular
	* session, even if the application does not reference it.
	*/
	access()
}

type StandardSession struct {
	lock sync.Mutex

	createTime          time.Time
	lastAccessedTime    time.Time
	id                  string
	maxInactiveInterval time.Duration
	accessTime          time.Time
	accessNumber        int16
	attributes          sync.Map //存储绑定在session里的数据
	manager             Manager  //session管理器,session本身不对自己进行管理，需要向manager递交请求从而管理自身
	expiring            bool     //是否过期
	activityCheck       bool     //是否检查有效链接
	accessCount         uint32   //有效链接个数
	isValidChan         chan bool
}

//生成标准session
func NewStandardSession(manager Manager) *StandardSession {
	s := &StandardSession{
		manager:             manager,
		createTime:          time.Now(),
		id:                  manager.SessionId(),
		maxInactiveInterval: time.Duration(manager.MaxLifeTime())}
	manager.addSession(s)
	s.lastAccessedTime = s.createTime
	s.activityCheck = true
	s.isValidChan = make(chan bool, 1)
	s.isValidChan <- true
	return s
}

func (session *StandardSession) SessionId() string {

	return session.id
}
func (session *StandardSession) LastAccessedTime() time.Time {
	return session.lastAccessedTime
}
func (session *StandardSession) SetMaxInactiveInterval(d time.Duration) {
	session.maxInactiveInterval = d
}
func (session *StandardSession) MaxInactiveInterval() time.Duration {
	return session.maxInactiveInterval
}

func (session *StandardSession) Load(key string) (value interface{}, ok bool) {
	return session.attributes.Load(key)
}
func (session *StandardSession) Store(key string, value interface{}) {
	session.attributes.Store(key, value)
}
func (session *StandardSession) Remove(key string) {
	session.attributes.Delete(key)
}

func (session *StandardSession) RemoveAll() {
	session.attributes = sync.Map{}
}

func (session *StandardSession) access() {
	session.lastAccessedTime = time.Now()
	if session.activityCheck {
		atomic.AddUint32(&session.accessCount, 1)
	}
}
func (session *StandardSession) Invalidate(w http.ResponseWriter) {
	if !session.IsValid() {
		panic(" procession had done")
	}
	session.Expire(w)

}

func (session *StandardSession) Expire(w http.ResponseWriter) {
	if !session.IsValid() {
		return
	}
	//double check
	session.lock.Lock()
	defer session.lock.Unlock()
	defer func() {
		session.expiring = false
	}()
	if session.expiring {
		return
	}
	session.expiring = true
	session.SetIsValid(false)
	session.RemoveAll()
	//向管理器递交移除该session请求
	session.manager.RemoveSession(session.id,w)
}

//实现对IsValid原子修改
func (session *StandardSession) SetIsValid(isValid bool) {
	<-session.isValidChan
	session.isValidChan <- isValid
}

//实现对IsValid原子读取
func (session *StandardSession) IsValid() bool {
	b := <-session.isValidChan
	session.isValidChan <- b
	return b
}
