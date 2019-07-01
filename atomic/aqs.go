package atomic

import "sync/atomic"

type QueuedSynchronizer struct {
	Status *int32
	Head   *Node
	Tail   *Node
}

type Node struct {
}

type Queue struct {
}

func (qs QueuedSynchronizer) compareAndSetStatus(newStatus int32) {
	atomic.CompareAndSwapInt32(qs.Status, *qs.Status, newStatus)
}

/**
独占式获取同步状态，子类实现重写该方法。实现该方法需要查询当前状态并判断同步状态是否符合预期，然后再进行CAS设置同步状态
此方法总是由执行线程调用,获取。如果该方法返回失败，则应该将执行协程加入到同步队列中
 */
func (qs QueuedSynchronizer) tryAcquire(arg int32) bool {
	panic("UnsupportedOperationException ")
}
/**
独占式释放同步状态，等待获取同步状态的协程将有机会获取同步状态
 */
func (qs QueuedSynchronizer) tryRelease(arg int32) bool {
	panic("UnsupportedOperationException ")
}

func (qs QueuedSynchronizer) setExclusiveOwnerThread( goId string) {

}
func acquire(q QueuedSynchronizer,arg int32,tryAcquire func(arg int32)bool) {
	if !tryAcquire(arg) {
		
	}
	
}
func release() {

}
