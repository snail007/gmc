// Copyright 2020 The GMC Author. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
// More information at https://github.com/snail007/gmc

// Package list implements a doubly linked list.
//
// To iterate over a list (where l is a *LinkList):
//	for e := l.Front(); e != nil; e = e.Next() {
//		// do something with e.Value
//	}
//
package glinklist

import (
	"sync"
	"sync/atomic"
)

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *LinkList

	// The value stored with this element.
	Value interface{}
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// LinkList represents a doubly linked list.
// The zero value for LinkList is an empty list ready to use.
type LinkList struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  *int64  // current list length excluding (this) sentinel element
	sync.RWMutex
}

// Init initializes or clears list l.
func (l *LinkList) Init() *LinkList {
	l.Lock()
	defer l.Unlock()
	return l.init()
}

func (l *LinkList) init() *LinkList {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = new(int64)
	return l
}

// New returns an initialized list.
func New() *LinkList { return new(LinkList).init() }

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *LinkList) Len() int64 { return atomic.LoadInt64(l.len) }

// Front returns the first element of list l or nil if the list is empty.
func (l *LinkList) Front() *Element {
	l.RLock()
	defer l.RUnlock()
	return l.front()
}

func (l *LinkList) front() *Element {
	if atomic.LoadInt64(l.len) == 0 {
		return nil
	}
	return l.root.next
}

// Back returns the last element of list l or nil if the list is empty.
func (l *LinkList) Back() *Element {
	l.RLock()
	defer l.RUnlock()
	return l.back()
}

func (l *LinkList) back() *Element {
	if atomic.LoadInt64(l.len) == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit lazily initializes a zero LinkList value.
func (l *LinkList) lazyInit() {
	if l.root.next == nil {
		l.init()
	}
}

// insert inserts e after at, increments l.len, and returns e.
func (l *LinkList) insert(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	atomic.AddInt64(l.len, 1)
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *LinkList) insertValue(v interface{}, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}

// remove removes e from its list, decrements l.len, and returns e.
func (l *LinkList) removeElement(e *Element) *Element {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil // avoid memory leaks
	e.prev = nil // avoid memory leaks
	e.list = nil
	atomic.AddInt64(l.len, -1)
	return e
}

// move moves e to next to at and returns e.
func (l *LinkList) move(e, at *Element) *Element {
	if e == at {
		return e
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e

	return e
}

// Remove removes e from l if e is an element of list l.
// It returns the element value e.Value.
// The element must not be nil.
func (l *LinkList) Remove(e *Element) interface{} {
	l.Lock()
	defer l.Unlock()
	return l.remove(e)
}

func (l *LinkList) remove(e *Element) interface{} {
	if e.list == l {
		// if e.list == l, l must have been initialized when e was inserted
		// in l or l == nil (e is a zero Element) and l.remove will crash
		l.removeElement(e)
	}
	return e.Value
}

// PushFront inserts a new element e with value v at the front of list l and returns e.
func (l *LinkList) PushFront(v interface{}) *Element {
	l.Lock()
	defer l.Unlock()
	return l.pushFront(v)
}

func (l *LinkList) pushFront(v interface{}) *Element {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

// PushBack inserts a new element e with value v at the back of list l and returns e.
func (l *LinkList) PushBack(v interface{}) *Element {
	l.Lock()
	defer l.Unlock()
	return l.pushBack(v)
}

func (l *LinkList) pushBack(v interface{}) *Element {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

// InsertBefore inserts a new element e with value v immediately before mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *LinkList) InsertBefore(v interface{}, mark *Element) *Element {
	l.Lock()
	defer l.Unlock()
	// see comment in LinkList.Remove about initialization of l
	return l.insertBefore(v, mark)
}

func (l *LinkList) insertBefore(v interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in LinkList.Remove about initialization of l
	return l.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element e with value v immediately after mark and returns e.
// If mark is not an element of l, the list is not modified.
// The mark must not be nil.
func (l *LinkList) InsertAfter(v interface{}, mark *Element) *Element {
	l.Lock()
	defer l.Unlock()
	return l.insertAfter(v, mark)
}

func (l *LinkList) insertAfter(v interface{}, mark *Element) *Element {
	if mark.list != l {
		return nil
	}
	// see comment in LinkList.Remove about initialization of l
	return l.insertValue(v, mark)
}

// MoveToFront moves element e to the front of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *LinkList) MoveToFront(e *Element) {
	l.Lock()
	defer l.Unlock()
	// see comment in LinkList.Remove about initialization of l
	l.moveToFront(e)
}

func (l *LinkList) moveToFront(e *Element) {
	if e.list != l || l.root.next == e {
		return
	}
	// see comment in LinkList.Remove about initialization of l
	l.move(e, &l.root)
}

// MoveToBack moves element e to the back of list l.
// If e is not an element of l, the list is not modified.
// The element must not be nil.
func (l *LinkList) MoveToBack(e *Element) {
	l.Lock()
	defer l.Unlock()
	// see comment in LinkList.Remove about initialization of l
	l.moveToBack(e)
}

func (l *LinkList) moveToBack(e *Element) {
	if e.list != l || l.root.prev == e {
		return
	}
	// see comment in LinkList.Remove about initialization of l
	l.move(e, l.root.prev)
}

// MoveBefore moves element e to its new position before mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *LinkList) MoveBefore(e, mark *Element) {
	l.Lock()
	defer l.Unlock()
	l.moveBefore(e, mark)
}

func (l *LinkList) moveBefore(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark.prev)
}

// MoveAfter moves element e to its new position after mark.
// If e or mark is not an element of l, or e == mark, the list is not modified.
// The element and mark must not be nil.
func (l *LinkList) MoveAfter(e, mark *Element) {
	l.Lock()
	defer l.Unlock()
	l.moveAfter(e, mark)
}

func (l *LinkList) moveAfter(e, mark *Element) {
	if e.list != l || e == mark || mark.list != l {
		return
	}
	l.move(e, mark)
}

// PushBackList inserts a copy of an other list at the back of list l.
// The lists l and other may be the same. They must not be nil.
func (l *LinkList) PushBackList(other *LinkList) {
	l.Lock()
	defer l.Unlock()
	l.pushBackList(other)
}

func (l *LinkList) pushBackList(other *LinkList) {
	l.lazyInit()
	for i, e := other.Len(), other.front(); i > 0; i, e = i-1, e.Next() {
		l.insertValue(e.Value, l.root.prev)
	}
}

// PushFrontList inserts a copy of an other list at the front of list l.
// The lists l and other may be the same. They must not be nil.
func (l *LinkList) PushFrontList(other *LinkList) {
	l.Lock()
	defer l.Unlock()
	l.pushFrontList(other)
}

func (l *LinkList) pushFrontList(other *LinkList) {
	l.lazyInit()
	for i, e := other.Len(), other.back(); i > 0; i, e = i-1, e.Prev() {
		l.insertValue(e.Value, &l.root)
	}
}

// Range calls f sequentially for each key and value present in the linklist.
// If f returns false, range stops the iteration.
// Don't modify the map during the Range call, that will cause deadlock.
func (l *LinkList) Range(f func(v interface{}) bool) {
	l.RLock()
	defer l.RUnlock()
	if atomic.LoadInt64(l.len) == 0 {
		return
	}
	p := l.root.next
	for {
		if p == nil {
			break
		}
		if !f(p.Value) {
			return
		}
		p = p.Next()
	}
}

// IndexOf returns the value first occurs index in the linklist, index starts with 0.
// if not found return -1.
func (l *LinkList) IndexOf(v interface{}) int {
	l.RLock()
	defer l.RUnlock()
	p := l.front()
	i := 0
	for {
		if p == nil {
			break
		}
		if p.Value == v {
			return i
		}
		p = p.Next()
		i++
	}
	return -1
}

// Clone creates a copy of current list.
func (l *LinkList) Clone() *LinkList {
	newList := New()
	if atomic.LoadInt64(l.len) == 0 {
		return newList
	}
	l.RLock()
	defer l.RUnlock()
	p := l.root.next
	for {
		if p == nil {
			break
		}
		newList.PushBack(p.Value)
		p = p.Next()
	}
	return newList
}
