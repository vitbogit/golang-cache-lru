// Package list содержит определение двухсвязного списка и его элементов для использования в LRU-кэше
package list

import "time"

// LruList реализует двухсвязный список
type LruList struct {
	root Entry // служебный элемент списка
	len  int   // размер списка без служебного элемента root
}

// Init инициализирует (или чистит) двухсвязный список
func (l *LruList) Init() *LruList {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// NewList создает новый двухсвязный список
func NewList() *LruList {
	return new(LruList).Init()
}

// Length возвращает длину двухсвязного списка
func (l *LruList) Length() int {
	return l.len
}

// Back возвращает последний элемент списка или nil
func (l *LruList) Back() *Entry {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

// lazyInit инициализирует список, проверяя, есть ли указатель хотя бы на один элемент
func (l *LruList) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}

// insert вставляет элемент e после at в список
func (l *LruList) insert(e, at *Entry) *Entry {
	e.prev = at
	e.next = at.next
	e.prev.next = e // меняет at
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue оборачивает функцию insert для возможности "собрать" внутри нее новый элемент с
// заданными для него значениями
func (l *LruList) insertValue(k string, v interface{}, expiresAt time.Time, at *Entry) *Entry {
	return l.insert(&Entry{Value: v, Key: k, ExpiresAt: expiresAt}, at)
}

// Remove удаляет e из списка
func (l *LruList) Remove(e *Entry) interface{} {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	l.len--

	return e.Value
}

// move ставит e перед at
func (l *LruList) move(e, at *Entry) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

// PushFront оборачивает функцию insertValue для добавления в начало списка
func (l *LruList) PushFront(k string, v interface{}, expiresAt time.Time) *Entry {
	l.lazyInit()

	return l.insertValue(k, v, expiresAt, &l.root)
}

// MoveToFront перемещает e в начало списка
func (l *LruList) MoveToFront(e *Entry) {
	if e == nil || e.list != l || l.root.next == e {
		return
	}

	l.move(e, &l.root)
}
