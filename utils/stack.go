package utils

import "container/list"

type Stack struct {
	list *list.List
}

func NewStack() *Stack {
	list := list.New()
	return &Stack{list}
}

func (self *Stack) Push(value interface{}) {
	self.list.PushBack(value)
}

func (self *Stack) Pop() interface{} {
	e := self.list.Back()
	if e != nil {
		self.list.Remove(e)
		return e.Value
	}
	return nil
}

func (self *Stack) Peak() interface{} {
	e := self.list.Back()
	if e != nil {
		return e.Value
	}
	return nil
}

func (self *Stack) Len() int {
	return self.list.Len()
}

func (self *Stack) Empty() bool {
	return self.list.Len() == 0
}
