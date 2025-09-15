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
	head *ListItem
	tail *ListItem
	lens int
}

func NewList() List {
	return &list{}
}

func (l *list) Len() int {
	return l.lens
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.lens == 0 {
		firstItem := &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}
		l.head = firstItem
		l.tail = firstItem
		l.lens++
		return firstItem
	}

	newItem := &ListItem{
		Value: v,
		Next:  l.head,
		Prev:  nil,
	}
	l.head.Prev = newItem
	l.head = newItem
	l.lens++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.lens == 0 {
		firstItem := &ListItem{
			Value: v,
			Next:  nil,
			Prev:  nil,
		}
		l.head = firstItem
		l.tail = firstItem
		l.lens++
		return firstItem
	}

	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.tail,
	}
	l.tail.Next = newItem
	l.tail = newItem
	l.lens++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Next == nil && i.Prev == nil {
		l.head = nil
		l.tail = nil
		l.lens = 0
	} else {
		switch {
		case l.head == i:
			l.head = i.Next
			i.Next.Prev = nil
		case l.tail == i:
			l.tail = i.Prev
			i.Prev.Next = nil
		default:
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
		}
		l.lens--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	if l.tail == i {
		l.tail = i.Prev
		l.tail.Next = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i
}
