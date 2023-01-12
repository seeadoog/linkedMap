package linkedMap

type cmap[K comparable, V any] struct {
	index    map[K]int
	data     []V
	maxIndex int
}
