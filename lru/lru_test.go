package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	key1 := "key1"
	value1 := String("value1")
	lru.Add(key1, value1)
	if v, ok := lru.Get(key1); !ok || string(v.(String)) != "value1" {
		t.Fatalf("cache hit %s:%s failed", key1, value1)
	}
	key2 := "key2"
	value2 := String("value2")
	if _, ok := lru.Get(key2); ok {
		t.Fatalf("cache hit %s:%s failed", key2, value2)
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	capacity := len(k1 + k2 + v1 + v2)
	lru := New(int64(capacity), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))
	expLen := 2
	if _, ok := lru.Get(k1); ok || lru.Len() != expLen {
		t.Fatal("RemoveOldest key1 failed")
	}
}

func TestCache_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("v2"))
	lru.Add("k3", String("v3"))
	lru.Add("k4", String("v4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys: %v, fact keys: %v", expect, keys)
	}

}
