package LruCache

import (
	"container/list"
	"reflect"
	"sort"
	"testing"
	"time"
)

func Second(i int) time.Duration {
	return time.Duration(i) * time.Second
}

func feedLRU(entries ...Entry) *list.List {
	var l = list.New()
	for _, v := range entries {
		v.SetLruLink(l.PushFront(v))
	}
	return l
}

func TestNewCache(t *testing.T) {
	type args struct {
		size uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "Create cache",
			args: args{size: 128},
			want: 128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCache(tt.args.size)
			if got.(*cache).capacity != tt.args.size {
				t.Errorf("New cache capacity = %d, want %d", got.(*cache).capacity, tt.want)
			}
		})
	}
}

func Test_cache_Add(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))

	var e2a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e2b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))
	var e2c = NewEntry(NewStringKey("C"), "C entry", Second(10), Second(15))

	var e3a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e3b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))
	var e3c = NewEntry(NewStringKey("C"), "C entry", Second(10), Second(15))
	var e3d = NewEntry(NewStringKey("D"), "D entry", Second(10), Second(15))
	var e3e = NewEntry(NewStringKey("E"), "E entry", Second(10), Second(15))

	var e4a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e4aBis = NewEntry(NewStringKey("A"), "A entry, new version", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		args   Entry
		want   bool
	}{
		{
			name: "Add an entry in an empty cache",
			fields: fields{
				cacheMap: make(map[string]Entry, 16),
				cacheLRU: list.New(),
				capacity: 16,
			},
			args: e1a,
			want: false,
		},
		{
			name: "Add an entry in a non-empty cache",
			fields: fields{
				cacheMap: map[string]Entry{
					"A": e2a,
					"B": e2b,
				},
				cacheLRU: feedLRU(e2a, e2b),
				capacity: 16,
			},
			args: e2c,
			want: false,
		},
		{
			name: "Add an entry in a full cache",
			fields: fields{
				cacheMap: map[string]Entry{"A": e3a, "B": e3b, "C": e3c, "D": e3d},
				cacheLRU: feedLRU(e3a, e3b, e3c, e3d),
				capacity: 4,
			},
			args: e3e,
			want: true,
		},
		{
			name: "Overwrite an entry in the cache",
			fields: fields{
				cacheMap: map[string]Entry{"A": e4a},
				cacheLRU: feedLRU(e4a),
				capacity: 16,
			},
			args: e4aBis,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got := c.Add(tt.args)
			// Check the position in the lru cache
			if tt.args.GetLruLink() != c.cacheLRU.Front() {
				t.Error("The new cache entry should be the first entry in the LRU cache, and it is not.")
			}
			// Check it is present in the cache
			if _, exist := c.cacheMap[tt.args.Key().String()]; !exist {
				t.Error("The new cache entry was not found in the cache.")
			}
			// Check if the LRU entry was removed when cache is full
			if got != tt.want {
				t.Errorf("Got %t in return instead of %t", got, tt.want)
			}
			// Check is not expired
			if tt.args.IsExpired() {
				t.Error("New cache entry is expired and it should not.")
			}
			// Check entry value
			if tt.args.Value() != c.cacheMap[tt.args.Key().String()].Value() {
				t.Errorf("Wrong cache entry value, got (%v) instead of (%s).", tt.args.Value(), c.cacheMap[tt.args.Key().String()].Value())
			}
		})
	}
}

func Test_cache_Capacity(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "Get cache capacity",
			fields: fields{
				cacheMap: nil,
				cacheLRU: nil,
				capacity: 16,
			},
			want: 16,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			if got := c.Capacity(); got != tt.want {
				t.Errorf("Capacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_Contains(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   EntryKey
		want   bool
	}{
		{
			name: "The cache is supposed to contain A",
			fields: fields{
				cacheMap: map[string]Entry{
					"A": NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15)),
					"B": NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15)),
				},
				cacheLRU: nil,
				capacity: 16,
			},
			args: NewStringKey("A"),
			want: true,
		},
		{
			name: "The cache is not supposed to contain D",
			fields: fields{
				cacheMap: map[string]Entry{
					"A": NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15)),
					"B": NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15)),
				},
				cacheLRU: nil,
				capacity: 16,
			},
			args: NewStringKey("D"),
			want: false,
		},
		{
			name: "The cache is empty",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: nil,
				capacity: 16,
			},
			args: NewStringKey("A"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			if got := c.Contains(tt.args); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_Flush(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))
	var e1c = NewEntry(NewStringKey("C"), "C entry", Second(10), Second(15))
	var e1d = NewEntry(NewStringKey("D"), "D entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "Flush a non-empty cache",
			fields: fields{
				cacheMap: map[string]Entry{
					"A": e1a,
					"B": e1b,
					"C": e1c,
					"D": e1d,
				},
				cacheLRU: feedLRU(e1a, e1b, e1c, e1d),
				capacity: 16,
			},
			want: 4,
		},
		{
			name: "Flush an empty cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			want: 0,
		},
		{
			name: "Flush a cache full",
			fields: fields{
				cacheMap: map[string]Entry{
					"A": e1a,
					"B": e1b,
					"C": e1c,
					"D": e1d,
				},
				cacheLRU: feedLRU(e1a, e1b, e1c, e1d),
				capacity: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got := c.Flush()
			if got != tt.want {
				t.Errorf("Flush() = %v, want %v", got, tt.want)
			}
			if l := len(c.cacheMap); l != 0 {
				t.Errorf("Cache map is not flushed, len = %d", l)
			}
			if l := c.cacheLRU.Len(); l != 0 {
				t.Errorf("LRU cache is not flushed, len = %d", l)
			}
		})
	}
}

func Test_cache_Get(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}
	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	var e2a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e2b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		args   EntryKey
		want   Entry
	}{
		{
			name: "Get an existing entry",
			fields: fields{
				cacheMap: map[string]Entry{
					"A": e1a,
					"B": e1b,
				},
				cacheLRU: feedLRU(e1a, e1b),
				capacity: 16,
			},
			args: e1a.Key(),
			want: e1a,
		},
		{
			name: "Get a non-existing entry",
			fields: fields{
				cacheMap: map[string]Entry{
					"A": e2a,
					"B": e2b,
				},
				cacheLRU: feedLRU(e2a, e2b),
				capacity: 16,
			},
			args: NewStringKey("C"),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got := c.Get(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
			if got != nil {
				if got.(*entry).creationTime.Sub(got.(*entry).accessTime) == 0 {
					t.Error("Entry access time not updated")
				}
			}
		})
	}
}

func Test_cache_GetLruEntry(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	var e2a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		want   Entry
	}{
		{
			name: "Get LRU entry",
			fields: fields{
				cacheMap: map[string]Entry{"A": e1a, "B": e1b},
				cacheLRU: feedLRU(e1a, e1b),
				capacity: 16,
			},
			want: e1a,
		},
		{
			name: "Get LRU entry",
			fields: fields{
				cacheMap: map[string]Entry{"A": e2a},
				cacheLRU: feedLRU(e2a),
				capacity: 16,
			},
			want: e2a,
		},
		{
			name: "Get LRU entry",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			if got := c.GetLruEntry(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLruEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_GetWithoutAccessUpdate(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	var e3a = &entry{
		key:          NewStringKey("A"),
		value:        "A entry",
		ttl:          Second(30),
		maxAge:       Second(60),
		accessTime:   time.Now().Add(Second(-40)),
		creationTime: time.Now().Add(Second(-50)),
		lruElement:   nil,
	}
	var e3b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		args   EntryKey
		want   Entry
	}{
		{
			name: "Get an existing entry",
			fields: fields{
				cacheMap: map[string]Entry{"A": e1a, "B": e1b},
				cacheLRU: nil,
				capacity: 16,
			},
			args: e1a.Key(),
			want: e1a,
		},
		{
			name: "Get a non-existing entry",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: nil,
				capacity: 16,
			},
			args: e1a.Key(),
			want: nil,
		},
		{
			name: "Get an existing entry which is expired",
			fields: fields{
				cacheMap: map[string]Entry{"A": e3a, "B": e3b},
				cacheLRU: feedLRU(e3a, e3b),
				capacity: 16,
			},
			args: e3a.Key(),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got := c.GetWithoutAccessUpdate(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWithoutAccessUpdate() = %v, want %v", got, tt.want)
			}
			if got != nil {
				if got.(*entry).creationTime.Sub(got.(*entry).accessTime) != 0 {
					t.Error("Entry access time was updated")
				}
			}
		})
	}
}

func Test_cache_HouseCleaning(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a Entry = &entry{
		key:          NewStringKey("A"),
		value:        "A entry",
		ttl:          Second(30),
		maxAge:       Second(60),
		accessTime:   time.Now().Add(Second(-20)),
		creationTime: time.Now(),
		lruElement:   nil,
	}
	var e1b Entry = &entry{
		key:          NewStringKey("B"),
		value:        "B entry",
		ttl:          Second(30),
		maxAge:       Second(60),
		accessTime:   time.Now().Add(Second(-30)),
		creationTime: time.Now(),
		lruElement:   nil,
	}
	var e1c Entry = &entry{
		key:          NewStringKey("C"),
		value:        "C entry",
		ttl:          Second(30),
		maxAge:       Second(60),
		accessTime:   time.Now().Add(Second(-40)),
		creationTime: time.Now(),
		lruElement:   nil,
	}
	var e1d Entry = &entry{
		key:          NewStringKey("D"),
		value:        "D entry",
		ttl:          Second(30),
		maxAge:       Second(60),
		accessTime:   time.Now().Add(Second(-70)),
		creationTime: time.Now(),
		lruElement:   nil,
	}
	var e1e Entry = &entry{
		key:          NewStringKey("E"),
		value:        "E entry",
		ttl:          Second(30),
		maxAge:       Second(60),
		accessTime:   time.Now().Add(Second(-10)),
		creationTime: time.Now(),
		lruElement:   nil,
	}

	tests := []struct {
		name   string
		fields fields
		want   uint32
		want1  []Entry
	}{
		{
			name: "Clean the cache",
			fields: fields{
				cacheMap: map[string]Entry{"A": e1a, "B": e1b, "C": e1c, "D": e1d, "E": e1e},
				cacheLRU: feedLRU(e1a, e1b, e1c, e1d, e1e),
				capacity: 16,
			},
			want:  3,
			want1: []Entry{e1c, e1d, e1b},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got, got1 := c.HouseCleaning()
			if got != tt.want {
				t.Errorf("HouseCleaning() got = %v, want %v", got, tt.want)
			}
			sort.Slice(got1, func(i, j int) bool {
				return got1[i].Key().String() < got1[j].Key().String()
			})
			sort.Slice(tt.want1, func(i, j int) bool {
				return tt.want1[i].Key().String() < tt.want1[j].Key().String()
			})
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("HouseCleaning() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_cache_Keys(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		want   []EntryKey
	}{
		{
			name: "Get keys for a non-empty cache",
			fields: fields{
				cacheMap: map[string]Entry{"A": e1a, "B": e1b},
				cacheLRU: feedLRU(e1a, e1b),
				capacity: 16,
			},
			want: []EntryKey{e1a.Key(), e1b.Key()},
		},
		{
			name: "Get keys for an empty cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			want: []EntryKey{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			if got := c.Keys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_Len(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "Get size of a non-empty cache",
			fields: fields{
				cacheMap: map[string]Entry{"A": e1a, "B": e1b},
				cacheLRU: feedLRU(e1a, e1b),
				capacity: 16,
			},
			want: 2,
		},
		{
			name: "Get size of an empty cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			if got := c.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_Remove(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	var e2a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e2b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		args   EntryKey
		want   bool
		want1  Entry
	}{
		{
			name: "Remove an existing entry",
			fields: fields{
				cacheMap: map[string]Entry{"A": e1a, "B": e1b},
				cacheLRU: feedLRU(e1a, e1b),
				capacity: 16,
			},
			args:  e1b.Key(),
			want:  true,
			want1: e1b,
		},
		{
			name: "Remove a non-existing entry from a non-empty cache",
			fields: fields{
				cacheMap: map[string]Entry{"A": e2a, "B": e2b},
				cacheLRU: feedLRU(e2a, e2b),
				capacity: 16,
			},
			args:  NewStringKey("C"),
			want:  false,
			want1: nil,
		},
		{
			name: "Remove a non-existing entry from an empty cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			args:  NewStringKey("A"),
			want:  false,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got, got1 := c.Remove(tt.args)
			if got != tt.want {
				t.Errorf("Remove() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Remove() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_cache_RemoveLruEntry(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}
	var e1a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e1b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		want   Entry
	}{
		{
			name: "Get LRU entry from a non-empty cache",
			fields: fields{
				cacheMap: map[string]Entry{"A": e1a, "B": e1b},
				cacheLRU: feedLRU(e1a, e1b),
				capacity: 0,
			},
			want: e1a,
		},
		{
			name: "Get LRU entry from an empty cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 0,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got := c.RemoveLruEntry()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveLruEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cache_Resize(t *testing.T) {
	type fields struct {
		cacheMap map[string]Entry
		cacheLRU *list.List
		capacity uint32
	}

	var e4a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e4b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))
	var e4c = NewEntry(NewStringKey("C"), "C entry", Second(10), Second(15))
	var e4d = NewEntry(NewStringKey("D"), "D entry", Second(10), Second(15))

	var e5a = NewEntry(NewStringKey("A"), "A entry", Second(10), Second(15))
	var e5b = NewEntry(NewStringKey("B"), "B entry", Second(10), Second(15))
	var e5c = NewEntry(NewStringKey("C"), "C entry", Second(10), Second(15))
	var e5d = NewEntry(NewStringKey("D"), "D entry", Second(10), Second(15))

	tests := []struct {
		name   string
		fields fields
		args   uint32
		want   uint32
		want1  []Entry
	}{
		{
			name: "Upsize a cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			args:  16,
			want:  0,
			want1: []Entry{},
		},
		{
			name: "Upsize a cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			args:  32,
			want:  0,
			want1: []Entry{},
		},
		{
			name: "Downsize an empty cache",
			fields: fields{
				cacheMap: make(map[string]Entry),
				cacheLRU: list.New(),
				capacity: 16,
			},
			args:  8,
			want:  0,
			want1: []Entry{},
		},
		{
			name: "Downsize a non-empty cache with enough free room",
			fields: fields{
				cacheMap: map[string]Entry{"A": e4a, "B": e4b, "C": e4c, "D": e4d},
				cacheLRU: feedLRU(e4a, e4b, e4c, e4d),
				capacity: 16,
			},
			args:  8,
			want:  0,
			want1: []Entry{},
		},
		{
			name: "Downsize a non-empty cache with not enough free room",
			fields: fields{
				cacheMap: map[string]Entry{"A": e5a, "B": e5b, "C": e5c, "D": e5d},
				cacheLRU: feedLRU(e5a, e5b, e5c, e5d),
				capacity: 16,
			},
			args:  2,
			want:  2,
			want1: []Entry{e5a, e5b},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cache{
				cacheMap: tt.fields.cacheMap,
				cacheLRU: tt.fields.cacheLRU,
				capacity: tt.fields.capacity,
			}
			got, got1 := c.Resize(tt.args)
			if got != tt.want {
				t.Errorf("Resize() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Resize() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
