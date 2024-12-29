package list

import "time"

// Entry определяет запись в кэше
type Entry struct {
	// Следующий и предыдущий указатели в двусвязном списке элементов.
	next, prev *Entry

	// Родительский список
	list *LruList

	// Ключ
	Key string

	// Значение
	Value interface{}

	// Дата истечения
	ExpiresAt time.Time
}

// PrevEntry возвращает предыдущий элемент
func (e *Entry) PrevEntry() *Entry {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}
