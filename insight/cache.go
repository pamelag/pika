package insight

import (
	"sync"
	"time"

	"github.com/pamelag/pika/cab"
)

type cache struct {
	mu      sync.RWMutex
	entries map[cab.Trip]int
}

func newCache() *cache {
	return &cache{
		entries: make(map[cab.Trip]int),
	}
}

func (c *cache) add(medallion string, tripDate time.Time, rides int) {
	t := cab.Trip{Medallion: medallion,
		PickupDatetime: tripDate}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[t] = rides
}

func (c *cache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.entries {
		delete(c.entries, k)
	}
}

func (c *cache) get(medallion string, tripDate time.Time) int {
	t := cab.Trip{Medallion: medallion,
		PickupDatetime: tripDate}
	c.mu.RLock()
	count, _ := c.entries[t]
	c.mu.RUnlock()
	return count
}
