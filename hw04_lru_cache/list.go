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
	l.pushItemToFront(item)

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++

	item := &ListItem{Value: v}
	l.pushItemToBack(item)

	return item
}

func (l *list) Remove(i *ListItem) {
	l.len--

	l.removeItem(i)
}

func (l *list) MoveToFront(i *ListItem) {
	if l.head == i {
		return
	}

	l.removeItem(i)
	l.pushItemToFront(i)
}

func (l *list) pushItemToFront(i *ListItem) {
	if l.head == nil {
		l.head = i
		l.tail = i

		return
	}

	i.Next = l.head
	l.head.Prev = i
	l.head = i
}

func (l *list) pushItemToBack(i *ListItem) {
	if l.head == nil {
		l.head = i
		l.tail = i

		return
	}

	i.Prev = l.tail
	l.tail.Next = i
	l.tail = i
}

func (l *list) removeItem(i *ListItem) {
	if l.head == i && l.tail == i {
		l.head = nil
		l.tail = nil

		return
	}

	if l.head == i {
		l.head = i.Next
		l.head.Prev = nil

		return
	}

	if l.tail == i {
		l.tail = i.Prev
		l.tail.Next = nil

		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
}
