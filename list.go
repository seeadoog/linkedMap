package linkedMap

type list[K comparable, V any] struct {
	front *Elem[K, V]
	back  *Elem[K, V]
}

func newList[K comparable, V any]() *list[K, V] {
	l := &list[K, V]{}
	l.front = &Elem[K, V]{}
	l.back = &Elem[K, V]{}
	l.front.next = l.back
	l.back.pre = l.front
	return l
}

// front -> back
func (l *list[K, V]) pushBack(e *Elem[K, V]) {
	pre := l.back.pre
	pre.next = e
	e.next = l.back
	l.back.pre = e
	e.pre = pre
}

func (l *list[K, V]) remove(e *Elem[K, V]) {
	pre := e.pre
	pre.next = e.next
	e.next.pre = pre
}

func (l *list[K, V]) foreach(f func(e *Elem[K, V]) bool) {
	e := l.front.next
	for e != l.back {
		if !f(e) {
			return
		}
		e = e.next
	}
}
