package list

import (
	"errors"
	"fmt"
)

type (
	// Node represents a LinkedList node that has a value of generic type T.
	Node[T any] struct {
		prev, next *Node[T]
		list       *DLL[T]

		Value T
	}

	// DLL is a collection of Nodes that together comprise the doubly linked list.
	//
	// All nodes are of generic type T.
	DLL[T any] struct {
		head Node[T]
		tail Node[T]
		len  int
	}
)

// Prev is a getter for the node's previous element pointer.
func (n *Node[T]) Prev() *Node[T] {
	return n.prev
}

// Next is a getter for the node's next element pointer.
func (n *Node[T]) Next() *Node[T] {
	return n.next
}

// NewList creates a new doubly linked list for the provided type.
func NewDoublyLinkedList[T any]() *DLL[T] {
	// create dummy nodes for head and tail to make traversal easier
	head := Node[T]{}
	tail := Node[T]{}
	head.prev = nil
	head.next = &tail
	tail.prev = &head
	tail.next = nil

	return &DLL[T]{
		len:  0,
		head: head,
		tail: tail,
	}
}

// Len returns the length of the list.
func (l *DLL[T]) Len() int {
	return l.len
}

// InsertAtFront creates a new Node with value `val` and inserts that node at the front of the list.
func (l *DLL[T]) InsertAtFront(val T) (*Node[T], error) {
	newNode := Node[T]{
		prev:  &l.head,
		next:  l.head.next,
		list:  l,
		Value: val,
	}

	if err := l.MoveToFront(&newNode); err != nil {
		return nil, fmt.Errorf("unable to move to front: %w", err)
	}
	l.len++

	return &newNode, nil
}

func (l *DLL[T]) MoveToFront(n *Node[T]) error {
	if n == nil {
		return errors.New("cannot move a nil node")
	}

	l.head.next.prev = n
	l.head.next = n

	return nil
}

// RemoveFromBack is a helper method to remove elements from the end of the list (useful in LRU when
// evicting elements).
func (l *DLL[T]) RemoveFromBack() (*Node[T], error) {
	toRemove := l.tail.prev
	if err := l.Remove(toRemove); err != nil {
		return nil, err
	}

	return toRemove, nil
}

// Remove removes an arbitrary node from anywhere in the list.
func (l *DLL[T]) Remove(n *Node[T]) error {
	switch {
	case n == &l.head || n == &l.tail:
		// If the dummy tail's prev element is head, then the list is empty (from the outsider's perspective)
		// and it only contains the dummy nodes.
		return errors.New("list is empty: unable to dummy nodes")
	case n == nil:
		return errors.New("unable to remove nil node")
	}

	n.prev.next = n.next
	n.next.prev = n.prev
	l.len--

	return nil
}
