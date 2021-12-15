package LruCache

import (
	"container/list"
	"time"
)

// Cache is the cache interface.
type Cache interface {
	// Add adds a new entry in the cache and returns true if the oldest entry has been removed because the cache was full.
	Add(entry Entry) bool

	// Get returns the entry corresponding to the requested key. It returns nil if the entry doesn't exist or expired.
	Get(key EntryKey) Entry

	// GetWithoutAccessUpdate returns the entry corresponding to the requested key.
	// It returns true if the entry exists.
	// It doesn't the update the entry last access time.
	GetWithoutAccessUpdate(key EntryKey) Entry

	// GetLruEntry returns the oldest cache entry.
	// It doesn't the update the entry last access time.
	GetLruEntry() Entry

	// Contains returns true if the cache contains an entry for the requested key.
	Contains(key EntryKey) bool

	// Remove removes the cache entry corresponding the requested key.
	// It returns true if the entry exists, and the removed entry.
	Remove(key EntryKey) (bool, Entry)

	// RemoveLruEntry removes the least recently used cache entry, and returns it.
	RemoveLruEntry() Entry

	// Keys returns the list of cache entries keys.
	Keys() []EntryKey

	// Len returns the number of entries present in the cache.
	Len() uint32

	// Capacity returns the cache capacity.
	Capacity() uint32

	// Flush clears the cache and returns the number of entries flushed.
	Flush() uint32

	// Resize updates the cache capacity and returns the number of cache entries flushed during the downsizing and the slice of flushed entries.
	Resize(size uint32) (uint32, []Entry)

	// HouseCleaning triggers the cache cleaning and removes the entries expired.
	// It returns the number of flushed entries and the slice of them.
	HouseCleaning() (uint32, []Entry)
}

// Entry is the cache entry interface.
type Entry interface {

	// SetLruLink sets the link between the cache entry the LRU entry list.
	SetLruLink(link *list.Element)

	// GetLruLink sets the link between the cache entry the LRU entry list.
	GetLruLink() *list.Element

	// Key returns the entry key.
	Key() EntryKey

	// Value returns the entry value.
	Value() interface{}

	// SetTTL the entry TTL
	SetTTL(ttl time.Duration)

	// GetTTL returns the entry TTL value.
	GetTTL() time.Duration

	// SetMaxAge set the entry max age to maxAge.
	SetMaxAge(maxAge time.Duration)

	// GetMaxAge returns the entry max age.
	GetMaxAge() time.Duration

	// UpdateAccessTime updates the entry last access time to time.Now()
	UpdateAccessTime()

	// GetAccessTime returns the entry last access time.
	GetAccessTime() time.Time

	// GetAge returns the entry age.
	GetAge() time.Duration

	// GetElapsedTimeFromLastAccess returns the elapsed from the entry last access time.
	GetElapsedTimeFromLastAccess() time.Duration

	// GetDelayToTTL returns the remaining delay before to reach the entry TTL.
	GetDelayToTTL() time.Duration

	// GetDelayToMaxAge returns the remaining delay before to reach the entry max age.
	GetDelayToMaxAge() time.Duration

	// GetDurationBeforeFlush returns the remaining entry lifetime.
	GetDurationBeforeFlush() time.Duration

	// ExceedTTL returns true if the entry reached the TTL.
	ExceedTTL() bool

	// ExceedMaxAge return true if the entry reached the max age.
	ExceedMaxAge() bool

	// IsExpired returns true is the cache entry expired
	IsExpired() bool
}

// EntryKey is the entry key interface.
type EntryKey interface {
	// String returns the key string representation.
	String() string
}
