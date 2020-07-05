package greet

import (
	"sync"
)

// Cache struct for Greet entity
type Cache struct {
	mutex sync.RWMutex
	store map[string]Greet
}

// GetAll gets all data values stored in cache struct
func (c *Cache) GetAll() (store map[string]Greet) {
	return c.store
}

// Get data stored in cache struct related with provided `Greet` entity
func (c *Cache) Get(t Greet) (cached Greet, found bool) {
	key := t.Name
	c.mutex.RLock()
	response, found := c.store[key]
	c.mutex.RUnlock()
	return response, found
}

// Set a new `Greet` entity to cached values
func (c *Cache) Set(t Greet) {
	key := t.Name
	c.mutex.Lock()
	c.store[key] = t
	c.mutex.Unlock()
}

// SetIfAbsent set a new `Greet` entity to cached values and returns the last cached value
func (c *Cache) SetIfAbsent(t Greet) (cached Greet, found bool) {
	response, found := c.Get(t)
	if found {
		return response, found
	}
	c.Set(t)
	return response, found
}

// Release remove some cached value.
func (c *Cache) Release(t Greet) {
	key := t.Name
	_, ok := c.Get(t)
	// remove if present
	if ok {
		c.mutex.Lock()
		delete(c.store, key)
		c.mutex.Unlock()
	}
}

// NewCache initializes cache struct with provided options.
func NewCache() (*Cache, error) {
	c := &Cache{} // new object
	// initialize properties
	c.mutex = sync.RWMutex{}
	c.store = make(map[string]Greet)
	// return it
	return c, nil
}
