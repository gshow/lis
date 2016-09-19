package point

import (
	"container/list"
)

var expireQueue *list.List = list.New()

func ExpireQueueAdd(ps *PointShell) *list.Element {
	return expireQueue.PushBack(ps)
}

func ExpireQueueRead() *list.Element {
	ret := expireQueue.Front()
	if ret != nil {
		expireQueue.Remove(ret)
	}

	return ret
}

func ExpireQueueLen() int {
	return expireQueue.Len()
}
