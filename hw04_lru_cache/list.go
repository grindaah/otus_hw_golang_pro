package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front  *ListItem
	back   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.front != nil {
		li := &ListItem{
			Next:  l.front,
			Prev:  nil,
			Value: v,
		}
		l.front.Prev = li
		l.front = li
		l.length++
		return li
	} else {
		return l.insertFirstItem(v)
	}
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.back != nil {
		li := &ListItem{
			Next:  nil,
			Prev:  l.back,
			Value: v,
		}
		l.back.Next = li
		l.back = li
		l.length++
		return li
	} else {
		return l.insertFirstItem(v)
	}
}

func (l *list) insertFirstItem(v interface{}) *ListItem {
	li := &ListItem{
		Next:  nil,
		Prev:  nil,
		Value: v,
	}
	l.back = li
	l.front = li
	l.length++
	return li
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i == l.front {
		l.front = i.Next
	}
	if i == l.back {
		l.back = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}
	if i == l.back {
		l.back = l.back.Prev
	}
	i.Prev.Next = i.Next
	i.Next = l.front
	l.front.Prev = i
	i.Prev = nil
	l.front = i
}

func NewList() List {
	return new(list)
}
