package epss

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type jobPool struct {
	limit   int
	running int32
	list    list.List
	mu      sync.RWMutex
}

func (t *jobPool) Add(f func()) {
	t.mu.Lock()
	t.list.PushBack(f)
	t.mu.Unlock()

	for {
		size := atomic.LoadInt32(&t.running)
		// 限制 数
		if t.limit != -1 && int(size) >= t.limit {
			return
		}

		// 双重 判断 在 load 和 com 之间 任务有可能 结束了 一个 那就
		if atomic.CompareAndSwapInt32(&t.running, size, size+1) {
			break
		}
	}

	go t.jobFunc()
}

func (t *jobPool) Size() int {
	return int(atomic.LoadInt32(&t.running))
}

func (t *jobPool) Jobs() int {
	return t.list.Len()
}

func (t *jobPool) getJob() (value interface{}) {
	t.mu.Lock()
	if v := t.list.Front(); v != nil {
		value = t.list.Remove(v)
	}
	t.mu.Unlock()
	return
}

func (t *jobPool) jobFunc() {
	for {
		if job := t.getJob(); job != nil {
			//do
			job.(func())()
		} else {
			break
		}
	}

	atomic.AddInt32(&t.running, -1)
}
