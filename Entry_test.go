package LruCache

import (
	"container/list"
	"reflect"
	"testing"
	"time"
)

func Test_entry_ExceedMaxAge(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Expect TRUE",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       time.Duration(time.Second * 5),
				accessTime:   time.Time{},
				creationTime: time.Now().Add(time.Second * -7),
			},
			want: true,
		},
		{
			name: "Expect FALSE",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       time.Duration(time.Second * 5),
				accessTime:   time.Time{},
				creationTime: time.Now().Add(time.Second * -2),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.ExceedMaxAge(); got != tt.want {
				t.Errorf("ExceedMaxAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_ExceedTTL(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Expect TRUE",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 5),
				maxAge:       0,
				accessTime:   time.Now().Add(time.Second * -7),
				creationTime: time.Time{},
			},
			want: true,
		},
		{
			name: "Expect FALSE",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 5),
				maxAge:       0,
				accessTime:   time.Now().Add(time.Second * -2),
				creationTime: time.Time{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.ExceedTTL(); got != tt.want {
				t.Errorf("ExceedTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_GetAccessTime(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	var loc *time.Location
	var err error
	if loc, err = time.LoadLocation("Europe/Paris"); err != nil {
		panic(err)
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Date(2022, 1, 1, 1, 1, 1, 1, loc),
				creationTime: time.Time{},
			},
			want: time.Date(2022, 1, 1, 1, 1, 1, 1, loc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.GetAccessTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAccessTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_GetAge(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "Expect 5s",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Now().Add(time.Second * -5),
			},
			want: time.Duration(time.Second * 5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.GetAge(); got-tt.want > time.Duration(time.Millisecond*500) {
				t.Errorf("Time delta = %v, want < 500ms", got-tt.want)
			}
		})
	}
}

func Test_entry_GetDelayToMaxAge(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "Expect 5s",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       time.Duration(time.Second * 10),
				accessTime:   time.Time{},
				creationTime: time.Now().Add(time.Second * -5),
			},
			want: time.Duration(time.Second * 5),
		},
		{
			name: "Expect 0s",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       time.Duration(time.Second * 10),
				accessTime:   time.Time{},
				creationTime: time.Now().Add(time.Second * -15),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.GetDelayToMaxAge(); got-tt.want > time.Duration(time.Millisecond*500) {
				t.Errorf("Time delta = %v, want < 500ms", got-tt.want)
			}
		})
	}
}

func Test_entry_GetDelayToTTL(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "Expect 5s",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 10),
				maxAge:       0,
				accessTime:   time.Now().Add(time.Second * -5),
				creationTime: time.Time{},
			},
			want: time.Duration(time.Second * 5),
		},
		{
			name: "Expect 0s",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 10),
				maxAge:       0,
				accessTime:   time.Now().Add(time.Second * -15),
				creationTime: time.Time{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.GetDelayToTTL(); got-tt.want > time.Duration(time.Millisecond*500) {
				t.Errorf("Time delta = %v, want < 500ms", got-tt.want)
			}
		})
	}
}

func Test_entry_GetDurationBeforeFlush(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "Expect 10m due to TTL",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Minute * 15),
				maxAge:       time.Duration(time.Minute * 20),
				accessTime:   time.Now().Add(time.Minute * -5),
				creationTime: time.Now().Add(time.Minute * -10),
			},
			want: time.Duration(time.Minute * 10),
		},
		{
			name: "Expect 5m due to max age",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Minute * 15),
				maxAge:       time.Duration(time.Minute * 20),
				accessTime:   time.Now().Add(time.Minute * -5),
				creationTime: time.Now().Add(time.Minute * -15),
			},
			want: time.Duration(time.Minute * 5),
		},
		{
			name: "Expect 0s due to TTL timeout",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Minute * 5),
				maxAge:       time.Duration(time.Minute * 20),
				accessTime:   time.Now().Add(time.Minute * -7),
				creationTime: time.Now().Add(time.Minute * -10),
			},
			want: 0,
		},
		{
			name: "Expect 0s due to max age timeout",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Minute * 10),
				maxAge:       time.Duration(time.Minute * 20),
				accessTime:   time.Now().Add(time.Minute * -5),
				creationTime: time.Now().Add(time.Minute * -30),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			var td, got time.Duration
			got = e.GetDurationBeforeFlush()
			td = got - tt.want
			if td < time.Duration(time.Millisecond*-500) || td > time.Duration(time.Millisecond*500) {
				t.Errorf("Time delta = %v, want < 500ms", got-tt.want)
			}
		})
	}
}

func Test_entry_GetElapsedTimeFromLastAccess(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "Getting elapsed time from the last access",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Now(),
				creationTime: time.Time{},
			},
			want: time.Now().Add(time.Second * 5).Sub(time.Now()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			time.Sleep(time.Second * 5)
			if got := e.GetElapsedTimeFromLastAccess(); got-tt.want > time.Duration(time.Millisecond*500) {
				t.Errorf("Time delta = %v, want < 500ms", got-tt.want)
			}
		})
	}
}

func Test_entry_GetMaxAge(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       time.Duration(time.Second * 45),
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			want: time.Duration(time.Second * 45),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.GetMaxAge(); got != tt.want {
				t.Errorf("GetMaxAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_GetTTL(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "GetTTL",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 45),
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			want: time.Duration(time.Second * 45),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.GetTTL(); got != tt.want {
				t.Errorf("GetTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_Key(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	var nki func(int) EntryKey = func(i int) EntryKey {
		return NewIntKey(i)
	}
	var nks func(string) EntryKey = func(s string) EntryKey {
		return NewStringKey(s)
	}
	tests := []struct {
		name   string
		fields fields
		want   EntryKey
	}{
		{
			name: "Use int key",
			fields: fields{
				key:          nki(1),
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			want: nki(1),
		},
		{
			name: "Use string key",
			fields: fields{
				key:          nks("A"),
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			want: nks("A"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.Key(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_SetMaxAge(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	type args struct {
		maxAge time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "SetMAxAge",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			args: args{maxAge: time.Duration(time.Second * 60)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			e.SetMaxAge(tt.args.maxAge)
			if e.maxAge != tt.args.maxAge {
				t.Errorf("Mas age = %v, want %v", e.maxAge, tt.args.maxAge)
			}
		})
	}
}

func Test_entry_SetTTL(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	type args struct {
		ttl time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Set TTL",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			args: args{ttl: time.Duration(time.Second * 45)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			e.SetTTL(tt.args.ttl)
			if e.ttl != tt.args.ttl {
				t.Errorf("TTL = %v, want %v", e.ttl, tt.args.ttl)
			}
		})
	}
}

func Test_entry_UpdateAccessTime(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Access time update",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Now().Add(time.Minute * -5),
				creationTime: time.Time{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			e.UpdateAccessTime()
			if tDelta := time.Now().Sub(e.accessTime); tDelta > time.Duration(time.Second*1) {
				t.Errorf("Time delta = %v, want < 1s", tDelta)
			}
		})
	}
}

func Test_entry_Value(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	var ni func(int) *int = func(i int) *int {
		return &i
	}
	var ns func(string) *string = func(s string) *string {
		return &s
	}
	var nb func(bool) *bool = func(b bool) *bool {
		return &b
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "Int value",
			fields: fields{
				key:          nil,
				value:        ni(1),
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			want: ni(1),
		},
		{
			name: "String value",
			fields: fields{
				key:          nil,
				value:        ns("A"),
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			want: ns("A"),
		},
		{
			name: "Bool value",
			fields: fields{
				key:          nil,
				value:        nb(true),
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
			},
			want: nb(true),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_GetLruLink(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
		lruElement   *list.Element
	}
	tests := []struct {
		name   string
		fields fields
		want   *list.Element
	}{
		{
			name: "Get an lru link",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
				lruElement:   &list.Element{},
			},
			want: &list.Element{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
				lruElement:   tt.fields.lruElement,
			}
			if got := e.GetLruLink(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLruLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_SetLruLink(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
		lruElement   *list.Element
	}
	type args struct {
		link *list.Element
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          0,
				maxAge:       0,
				accessTime:   time.Time{},
				creationTime: time.Time{},
				lruElement:   nil,
			},
			args: args{link: &list.Element{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
				lruElement:   tt.fields.lruElement,
			}
			e.SetLruLink(tt.args.link)
			if !reflect.DeepEqual(e.lruElement, tt.args.link) {
				t.Error("Registering the lruElement in the cache entry failed")
			}
			if !reflect.DeepEqual(e, tt.args.link.Value) {
				t.Error("Registering the cache entry in the lruElement failed")
			}
		})
	}
}

func Test_entry_IsExpired(t *testing.T) {
	type fields struct {
		key          EntryKey
		value        interface{}
		ttl          time.Duration
		maxAge       time.Duration
		accessTime   time.Time
		creationTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "MaxAge timeout",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 5),
				maxAge:       time.Duration(time.Second * 15),
				accessTime:   time.Now().Add(time.Second * -2),
				creationTime: time.Now().Add(time.Second * -20),
			},
			want: true,
		},
		{
			name: "TTL timeout",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 5),
				maxAge:       time.Duration(time.Second * 15),
				accessTime:   time.Now().Add(time.Second * -7),
				creationTime: time.Now().Add(time.Second * -10),
			},
			want: true,
		},
		{
			name: "MaxAge timeout & TTL timeout",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 5),
				maxAge:       time.Duration(time.Second * 15),
				accessTime:   time.Now().Add(time.Second * -7),
				creationTime: time.Now().Add(time.Second * -20),
			},
			want: true,
		},
		{
			name: "No timeout",
			fields: fields{
				key:          nil,
				value:        nil,
				ttl:          time.Duration(time.Second * 5),
				maxAge:       time.Duration(time.Second * 15),
				accessTime:   time.Now().Add(time.Second * -2),
				creationTime: time.Now().Add(time.Second * -10),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &entry{
				key:          tt.fields.key,
				value:        tt.fields.value,
				ttl:          tt.fields.ttl,
				maxAge:       tt.fields.maxAge,
				accessTime:   tt.fields.accessTime,
				creationTime: tt.fields.creationTime,
			}
			if got := e.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entry_NewEntry(t *testing.T) {
	type args struct {
		key    EntryKey
		value  interface{}
		ttl    time.Duration
		maxAge time.Duration
	}
	type testCase struct {
		name string
		args args
		want Entry
	}
	var tNow = time.Now()
	tests := []testCase{
		{
			name: "New Entry A",
			args: args{
				key:    NewStringKey("A"),
				value:  "A entry",
				ttl:    time.Duration(time.Second * 10),
				maxAge: time.Duration(time.Second * 15),
			},
			want: &entry{
				key:          NewStringKey("A"),
				value:        "A entry",
				ttl:          time.Duration(time.Second * 10),
				maxAge:       time.Duration(time.Second * 15),
				accessTime:   tNow,
				creationTime: tNow,
				lruElement:   nil,
			},
		},
		{
			name: "New Entry 1",
			args: args{
				key:    NewIntKey(1),
				value:  100,
				ttl:    time.Duration(time.Second * 10),
				maxAge: time.Duration(time.Second * 15),
			},
			want: &entry{
				key:          NewIntKey(1),
				value:        100,
				ttl:          time.Duration(time.Second * 10),
				maxAge:       time.Duration(time.Second * 15),
				accessTime:   tNow,
				creationTime: tNow,
				lruElement:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEntry(tt.args.key, tt.args.value, tt.args.ttl, tt.args.maxAge)
			if g, w := got.(*entry).key, tt.want.(*entry).key; !reflect.DeepEqual(g, w) {
				t.Errorf("Wrong key, got (%s) and want (%s)", g, w)
			}
			if g, w := got.(*entry).value, tt.want.(*entry).value; g != w {
				t.Errorf("Wrong value, got (%s) and want (%s)", g, w)
			}
			if g, w := got.(*entry).ttl, tt.want.(*entry).ttl; g != w {
				t.Errorf("Wrong TTL, got (%s) and want (%s)", g, w)
			}
			if g, w := got.(*entry).maxAge, tt.want.(*entry).maxAge; g != w {
				t.Errorf("Wrong max age, got (%s) and want (%s)", g, w)
			}
			if g, w := got.(*entry).accessTime, tt.want.(*entry).accessTime; g.Sub(w) > time.Millisecond*2 {
				t.Errorf("Wrong access time delta, got %s and want < %s", g.Sub(w), time.Millisecond*2)
			}
			if g, w := got.(*entry).creationTime, tt.want.(*entry).creationTime; g.Sub(w) > time.Millisecond*2 {
				t.Errorf("Wrong creation time delta, got %s and want < %s", g.Sub(w), time.Millisecond*2)
			}
			if g := got.(*entry).lruElement; g != nil {
				t.Errorf("Wrong lru pointer, got (%v) and want (nil)", g)
			}
		})
	}
}
