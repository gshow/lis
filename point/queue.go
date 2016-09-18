package point

import (
	"container/list"
)

type ExpireQueueObject *list.Element

var ExpireQueue *pointQueue = &pointQueue{Queue: list.New()}

type pointQueue struct {
	Queue *list.List
}

func (this *pointQueue) Add(ps *PointShell) *list.Element {
	return this.Queue.PushBack(ps)
}

func (this *pointQueue) Read() *list.Element {
	ret := this.Queue.Front()
	if ret != nil {
		this.Queue.Remove(ret)
	}

	return ret
}

func (this *pointQueue) Len() int {
	return this.Queue.Len()
}
