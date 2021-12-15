package LruCache

import "strconv"

// STRING KEY

// entryStringKey is the private type representing a string key
type entryStringKey string

// String returns the key string representation
func (ek *entryStringKey) String() string {
	return string(*ek)
}

// NewStringKey returns a new key struct implementing the EntryKey interface
func NewStringKey(key string) EntryKey {
	var k = new(entryStringKey)
	*k = entryStringKey(key)
	return k
}

// INTEGER KEY

// entryStringKey is the private type representing a string key
type entryIntKey int

// String returns the key string representation
func (ek *entryIntKey) String() string {
	return strconv.Itoa(int(*ek))
}

// NewIntKey returns a new key struct implementing the EntryKey interface
func NewIntKey(key int) EntryKey {
	var k = new(entryIntKey)
	*k = entryIntKey(key)
	return k
}
