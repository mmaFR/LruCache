package LruCache

import (
	"container/list"
	"time"
)

type entry struct {
	key          EntryKey
	value        interface{}
	ttl          time.Duration
	maxAge       time.Duration
	accessTime   time.Time
	creationTime time.Time
	lruElement   *list.Element
}

// SetLruLink sets the link between the cache entry the LRU entry list.
func (e *entry) SetLruLink(link *list.Element) {
	link.Value = e
	e.lruElement = link
}

// GetLruLink sets the link between the cache entry the LRU entry list.
func (e *entry) GetLruLink() *list.Element {
	return e.lruElement
}

// Key returns the entry key.
func (e *entry) Key() EntryKey {
	return e.key
}

// Value returns the entry value.
func (e *entry) Value() interface{} {
	e.UpdateAccessTime()
	return e.value
}

// SetTTL sets the entry TTL.
func (e *entry) SetTTL(ttl time.Duration) {
	e.ttl = ttl
}

// GetTTL returns the entry TTL value.
func (e *entry) GetTTL() time.Duration {
	return e.ttl
}

// SetMaxAge set the entry max age to maxAge.
func (e *entry) SetMaxAge(maxAge time.Duration) {
	e.maxAge = maxAge
}

// GetMaxAge returns the entry max age.
func (e *entry) GetMaxAge() time.Duration {
	return e.maxAge
}

// UpdateAccessTime updates the entry last access time to time.Now()
func (e *entry) UpdateAccessTime() {
	e.accessTime = time.Now()
}

// GetAccessTime returns the entry last access time.
func (e *entry) GetAccessTime() time.Time {
	return e.accessTime
}

// GetAge returns the entry age.
func (e *entry) GetAge() time.Duration {
	return time.Now().Sub(e.creationTime)
}

// GetElapsedTimeFromLastAccess returns the elapsed from the entry last access time.
func (e *entry) GetElapsedTimeFromLastAccess() time.Duration {
	return time.Now().Sub(e.accessTime)
}

// GetDelayToTTL returns the remaining delay before to reach the entry ttl.
func (e *entry) GetDelayToTTL() time.Duration {
	var dtt time.Duration = e.ttl - e.GetElapsedTimeFromLastAccess()
	if dtt < 0 {
		return 0
	} else {
		return dtt
	}
}

// GetDelayToMaxAge returns the remaining delay before to reach the entry max age.
func (e *entry) GetDelayToMaxAge() time.Duration {
	var dtma time.Duration = e.maxAge - e.GetAge()
	if dtma < 0 {
		return 0
	} else {
		return dtma
	}
}

// GetDurationBeforeFlush returns the remaining entry lifetime.
func (e *entry) GetDurationBeforeFlush() time.Duration {
	var dtt, dtma time.Duration
	dtt = e.GetDelayToTTL()
	dtma = e.GetDelayToMaxAge()
	if dtt < dtma {
		return dtt
	} else {
		return dtma
	}
}

// ExceedTTL returns true if the entry reached the TTL.
func (e *entry) ExceedTTL() bool {
	if e.GetElapsedTimeFromLastAccess() > e.GetTTL() {
		return true
	} else {
		return false
	}
}

// ExceedMaxAge return true if the entry reached the max age.
func (e *entry) ExceedMaxAge() bool {
	if e.GetAge() > e.GetMaxAge() {
		return true
	} else {
		return false
	}
}

// IsExpired returns true is the cache entry expired
func (e *entry) IsExpired() bool {
	switch {
	case e.ExceedTTL():
		return true
	case e.ExceedMaxAge():
		return true
	default:
		return false
	}
}

func NewEntry(key EntryKey, value interface{}, ttl time.Duration, maxAge time.Duration) Entry {
	var ne = new(entry)
	ne.key = key
	ne.value = value
	ne.ttl = ttl
	ne.maxAge = maxAge
	ne.creationTime = time.Now()
	ne.accessTime = ne.creationTime
	return ne
}
