package LruCache

import (
	"reflect"
	"strconv"
	"testing"
)

func TestNewIntKey(t *testing.T) {
	type args struct {
		key int
	}
	type testCase struct {
		name string
		args args
		want EntryKey
	}
	var f func(int) EntryKey = func(i int) EntryKey {
		var v = entryIntKey(i)
		return &v
	}
	tests := []testCase{
		{
			name: "Want a string key equal to A",
			args: args{key: 1},
			want: f(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIntKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIntKey() = %v, want %v", got, tt.want)
				if got.String() != strconv.Itoa(tt.args.key) {
					t.Errorf("NewStringKey().String() = %s, want %d", got.String(), tt.args.key)
				}
			}
		})
	}
}

func TestNewStringKey(t *testing.T) {
	type args struct {
		key string
	}
	type testCase struct {
		name string
		args args
		want EntryKey
	}
	var f func(string) *entryStringKey = func(a string) *entryStringKey {
		var v = entryStringKey(a)
		return &v
	}
	tests := []testCase{
		{
			name: "Want a string key equal to A",
			args: args{key: "A"},
			want: f("A"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStringKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStringKey() = %v, want %v", got, tt.want)
				if got.String() != tt.args.key {
					t.Errorf("NewStringKey().String() = %s, want %s", got.String(), tt.args.key)
				}
			}
		})
	}
}

func Test_entryIntKey_String(t *testing.T) {
	type testCase struct {
		name string
		ek   entryIntKey
		want string
	}
	tests := []testCase{
		{
			name: "Want an int key equal to 1",
			ek:   entryIntKey(1),
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ek.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryStringKey_String(t *testing.T) {
	type testCase struct {
		name string
		ek   entryStringKey
		want string
	}
	tests := []testCase{
		{
			name: "Want a string key equal to A",
			ek:   entryStringKey("A"),
			want: "A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ek.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
