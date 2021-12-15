package LruCache

import "container/list"

type cache struct {
	cacheMap map[string]Entry
	cacheLRU *list.List
	capacity uint32
}

// Add adds a new entry in the cache and returns true if the oldest entry has been removed because the cache was full.
func (c *cache) Add(cacheEntry Entry) bool {
	var removeOldest bool
	if c.Len() == c.capacity {
		c.RemoveLruEntry()
		removeOldest = true
	}
	c.cacheMap[cacheEntry.Key().String()] = cacheEntry
	cacheEntry.SetLruLink(c.cacheLRU.PushFront(cacheEntry))
	return removeOldest
}

// Get returns the entry corresponding to the requested key. It returns nil if the entry doesn't exist or expired.
func (c *cache) Get(key EntryKey) Entry {
	if cacheEntry := c.GetWithoutAccessUpdate(key); cacheEntry != nil {
		cacheEntry.UpdateAccessTime()
		return cacheEntry
	} else {
		return nil
	}
}

// GetWithoutAccessUpdate returns the entry corresponding to the requested key.
// It returns true if the entry exists.
// It doesn't the update the entry last access time.
func (c *cache) GetWithoutAccessUpdate(key EntryKey) Entry {
	if cacheEntry, exists := c.cacheMap[key.String()]; exists {
		if cacheEntry.IsExpired() {
			c.Remove(cacheEntry.Key())
			return nil
		}
		return cacheEntry
	} else {
		return nil
	}
}

// GetLruEntry returns the oldest cache entry.
// It doesn't the update the entry last access time.
func (c *cache) GetLruEntry() Entry {
	if c.cacheLRU.Len() == 0 {
		return nil
	} else {
		return c.cacheLRU.Back().Value.(Entry)
	}
}

// Contains returns true if the cache contains an entry for the requested key.
func (c *cache) Contains(key EntryKey) bool {
	var b bool
	_, b = c.cacheMap[key.String()]
	return b
}

// Remove removes the cache entry corresponding the requested key.
// It returns true if the entry exists, and the removed entry.
func (c *cache) Remove(key EntryKey) (bool, Entry) {
	if cacheEntry, exists := c.cacheMap[key.String()]; exists {
		c.cacheLRU.Remove(cacheEntry.GetLruLink())
		delete(c.cacheMap, key.String())
		return true, cacheEntry
	} else {
		return false, nil
	}
}

// RemoveLruEntry removes the least recently used cache entry, and returns it.
func (c *cache) RemoveLruEntry() Entry {
	var removedEntry Entry = c.GetLruEntry()
	if removedEntry == nil {
		return nil
	} else {
		_, removedEntry = c.Remove(removedEntry.Key())
		return removedEntry
	}
}

// Keys returns the list of cache entries keys.
func (c *cache) Keys() []EntryKey {
	var ek = make([]EntryKey, c.Len())
	var idx uint32
	for _, v := range c.cacheMap {
		ek[idx] = v.Key()
		idx++
	}
	return ek
}

// Len returns the number of entries present in the cache.
func (c *cache) Len() uint32 {
	return uint32(len(c.cacheMap))
}

// Capacity returns the cache capacity.
func (c *cache) Capacity() uint32 {
	return c.capacity
}

// Flush clears the cache and returns the number of entries flushed.
func (c *cache) Flush() uint32 {
	var numberOfEntries = c.Len()
	c.cacheMap = make(map[string]Entry, c.capacity)
	c.cacheLRU = list.New()
	return numberOfEntries
}

// Resize updates the cache capacity and returns the number of cache entries flushed during the downsizing and the slice of flushed entries.
func (c *cache) Resize(size uint32) (uint32, []Entry) {
	var flushedEntries = make([]Entry, 0, 0)
	var numberOfEntries = c.Len()
	var numberOfDeletion uint32
	if size < numberOfEntries {
		numberOfDeletion = numberOfEntries - size
		flushedEntries = make([]Entry, numberOfDeletion)
		for i := uint32(0); i < numberOfDeletion; i++ {
			flushedEntries[i] = c.RemoveLruEntry()
		}
	}
	c.capacity = size
	return numberOfDeletion, flushedEntries
}

// HouseCleaning triggers the cache cleaning and removes the entries expired.
// It returns the number of flushed entries and the slice of them.
func (c *cache) HouseCleaning() (uint32, []Entry) {
	var flushedEntry = make([]Entry, 0)
	var numberOfDeletions uint32
	for _, cacheEntry := range c.cacheMap {
		if cacheEntry.IsExpired() {
			flushedEntry = append(flushedEntry, cacheEntry)
			c.Remove(cacheEntry.Key())
			numberOfDeletions++
		}
	}
	return numberOfDeletions, flushedEntry
}

func NewCache(size uint32) Cache {
	var nc = new(cache)
	nc.cacheMap = make(map[string]Entry, size)
	nc.cacheLRU = list.New()
	nc.capacity = size
	return nc
}
