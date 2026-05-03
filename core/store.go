package core

import (
	"time"

	"github.com/adityjoshi/iDB/config"
)

var store map[string]*Object

type Object struct {
	value     interface{}
	ExpiresAt int64
}

func init() {
	store = make(map[string]*Object)
}

func NewObj(value interface{}, expiresAtMS int64) *Object {
	var expiresAt int64 = -1

	if expiresAtMS > 0 {
		expiresAt = time.Now().UnixMilli() + expiresAtMS
	}

	return &Object{
		value:     value,
		ExpiresAt: expiresAt,
	}
}

func Put(key string, obj *Object) {
	if len(store) >= config.KeysLimit {
		evict()
	}
	store[key] = obj
}

func Get(key string) *Object {
	value := store[key]

	if value != nil {
		if value.ExpiresAt != -1 && value.ExpiresAt <= time.Now().UnixMilli() {
			delete(store, key)
			return nil
		}
	}

	return value
}

func Del(key string) bool {
	if _, ok := store[key]; ok {
		delete(store, key)
		return true
	}

	return false

}
