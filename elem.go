package linkedMap

type Elem[K comparable, V any] struct {
	Key  K
	Val  V
	next *Elem[K, V]
	pre  *Elem[K, V]
}
