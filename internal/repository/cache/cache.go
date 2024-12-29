// package cache
package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	def "github.com/vitbogit/golang-cache-lru/internal/repository"
	"github.com/vitbogit/golang-cache-lru/internal/repository/cache/list"
)

const (
	defaultEvicterFrequency = 50 * time.Nanosecond
)

var _ def.ILRUCache = (*LRU)(nil)

// LRU имплементирует потокобезопасный LRU-кэш с поддержкой TTL
type LRU struct {
	size      int
	evictList *list.LruList
	items     map[string]*list.Entry

	mu         sync.Mutex
	defaultTTL time.Duration
	done       chan struct{}
}

// NewCache создает новый кэш
func NewCache(size int, defaultTTL time.Duration) *LRU {
	if size <= 0 {
		log.Fatal().Msg("cache size can not be zero")
	}
	if defaultTTL <= 0 {
		log.Fatal().Msg("default TTL for cache can not be zero")
	}

	res := LRU{
		size:       size,
		evictList:  list.NewList(),
		items:      make(map[string]*list.Entry),
		defaultTTL: defaultTTL,
		done:       make(chan struct{}),
	}

	// Тикер, который будет раз в defaultEvicterFrequency времени
	// запускать deleteExpired() для удаления старых элементов
	go func(done <-chan struct{}) {
		ticker := time.NewTicker(defaultEvicterFrequency)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				res.deleteExpired()
			}
		}
	}(res.done)

	return &res
}

// EvictAll ручная инвалидация всего кэша
func (c *LRU) EvictAll(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Очистка мапы значений
	for k := range c.items {
		delete(c.items, k)
	}

	// Очистка двухсвязного списка
	c.evictList.Init()

	return nil
}

// Put запись данных в кэш
func (c *LRU) Put(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if len(key) == 0 || ttl < 0 {
		return fmt.Errorf("некорректные входные данные")
	}

	if ttl == 0 {
		ttl = c.defaultTTL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	// Перезапись существующего элемента
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value = value
		ent.ExpiresAt = now.Add(ttl)
		return nil
	}

	// Добавление в список
	ent := c.evictList.PushFront(key, value, now.Add(ttl)) // может "переполнить" список
	if c.evictList.Length() > c.size {                     // удаление лишнего элемента сзади
		c.removeOldest()
	}

	// Добавление в мапу
	c.items[key] = ent

	return nil
}

// Get получение данных из кэша по ключу
func (c *LRU) Get(ctx context.Context, key string) (value interface{}, expiresAt time.Time, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ent, ok := c.items[key]; ok {
		// Дополнительная проверка на expired
		if time.Now().After(ent.ExpiresAt) {
			// возвращаем nil error для not found
			return nil, time.Time{}, nil
		}

		// Успешно найдено
		return ent.Value, ent.ExpiresAt, nil
	}

	// возвращаем nil error для not found
	return nil, time.Time{}, nil
}

// Evict ручное удаление данных по ключу
func (c *LRU) Evict(ctx context.Context, key string) (value interface{}, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return ent.Value, nil
	}

	return nil, nil
}

// GetAll получение всего наполнения кэша в виде двух слайсов: слайса ключей и слайса значений.
// Пары ключ-значения из кэша располагаются на соответствующих позициях в слайсах.
func (c *LRU) GetAll(ctx context.Context) (keys []string, values []interface{}, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys = make([]string, 0, len(c.items))
	values = make([]interface{}, 0, len(c.items))

	now := time.Now()

	for ent := c.evictList.Back(); ent != nil; ent = ent.PrevEntry() {
		// Дополнительная проверка на expired
		if now.After(ent.ExpiresAt) {
			continue
		}

		keys = append(keys, ent.Key)
		values = append(values, ent.Value)
	}
	return keys, values, nil
}

// removeOldest удаляет старейший элемент. Подразумевается, что уже вызван lock.
func (c *LRU) removeOldest() {
	if ent := c.evictList.Back(); ent != nil {
		c.removeElement(ent)
	}
}

// removeElement удаляет указанный элемент. Подразумевается, что уже вызван lock.
func (c *LRU) removeElement(e *list.Entry) {
	c.evictList.Remove(e)  // удаление из списка
	delete(c.items, e.Key) // удаление из мапы
}

// deleteExpired вызывается каждые defaultEvicterFrequency времени специальной горутиной
// для удаления expired элементов
func (c *LRU) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	var nextEnt *list.Entry
	for ent := c.evictList.Back(); ; ent = nextEnt {
		if ent == nil {
			break
		}

		nextEnt = ent.PrevEntry()

		// Дополнительная проверка на expired
		if now.After(ent.ExpiresAt) {
			c.removeElement(ent)
		}
	}
}
