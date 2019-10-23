package utils

import "container/list"

type Queue struct {
	list *list.List
}

func NewQueue() *Queue {
	list := list.New()
	return &Queue{list}
}

func (self *Queue) Enqueue(value interface{}) {
	self.list.PushBack(value)
}

func (self *Queue) Dequeue() interface{} {
	e := self.list.Front()
	if e != nil {
		self.list.Remove(e)
		return e.Value
	}
	return nil
}

func (self *Queue) Peak() interface{} {
	e := self.list.Front()
	if e != nil {
		return e.Value
	}
	return nil
}

func (self *Queue) Len() int {
	return self.list.Len()
}

func (self *Queue) Empty() bool {
	return self.list.Len() == 0
}
