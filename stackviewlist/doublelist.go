package stackviewlist

import (
	"fmt"
	"iter"
)

type doubleList[T any] struct {
	head *node[T]
	end  *node[T]
}

type node[T any] struct {
	value T
	prev  *node[T]
	next  *node[T]
}

func newEmptyDoubleList[T any]() doubleList[T] {
	return doubleList[T]{}
}

func newDoubleList[T any](items []T) doubleList[T] {
	d := newEmptyDoubleList[T]()
	for _, item := range items {
		d.append(item)
	}
	return d
}

func (d *doubleList[T]) isEmpty() bool {
	return d.head == nil
}

func (d *doubleList[T]) length() int {
	if d.head == nil {
		return 0
	}

	var res int
	for range d.All() {
		res++
	}
	return res
}

func (d *doubleList[T]) clear() *doubleList[T] {
	d.head = nil
	d.end = nil
	return d
}

func (d *doubleList[T]) append(value T) *doubleList[T] {
	n := node[T]{value: value}
	if d.end == nil {
		d.head = &n
		d.end = &n
		return d
	}
	n.prev = d.end
	d.end.next = &n
	d.end = &n
	return d
}

func (d *doubleList[T]) prepend(value T) *doubleList[T] {
	n := node[T]{value: value}
	if d.head == nil {
		d.head = &n
		d.end = &n
		return d
	}
	n.next = d.head
	d.head.prev = &n
	d.head = &n
	return d
}

func (d *doubleList[T]) remove(n *node[T]) *doubleList[T] {
	if d.isEmpty() {
		return d
	}
	if d.head == n && d.end == n {
		d.head = nil
		d.end = nil
		return d
	}
	if d.head == n {
		d.head = d.head.next
		if d.head != nil {
			d.head.prev = nil
		}
		n.next = nil
		return d
	}
	if d.end == n {
		d.end = d.end.prev
		d.end.next = nil
		n.prev = nil
		return d
	}
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
	return d
}

func (d *doubleList[T]) String() string {
	s := "[[ "
	cur := d.head
	for cur != nil {
		s += fmt.Sprint(cur.value) + " "
		cur = cur.next
	}
	s += " ]]"
	return s
}

func (d *doubleList[T]) All() iter.Seq[*T] {
	return func(yield func(value *T) bool) {
		cur := d.head
		for cur != nil {
			if !yield(&cur.value) {
				break
			}
			cur = cur.next
		}
	}
}
