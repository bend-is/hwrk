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
	len        int
	head, tail *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int         { return l.len }
func (l *list) Back() *ListItem  { return l.tail }
func (l *list) Front() *ListItem { return l.head }

func (l *list) PushFront(v interface{}) *ListItem {
	l.len++

	item := &ListItem{Value: v}

	if l.head == nil {
		l.head = item
		l.tail = item

		return item
	}

	item.Next = l.head
	l.head.Prev = item
	l.head = item

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++

	item := &ListItem{Value: v}

	if l.head == nil {
		l.head = item
		l.tail = item

		return item
	}

	item.Prev = l.tail
	l.tail.Next = item
	l.tail = item

	return item
}

func (l *list) Remove(i *ListItem) {
	l.len--

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Next = l.head
	l.head.Prev = i
	l.head = i
}
