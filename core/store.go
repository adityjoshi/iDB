package core

import (
	"time"
)

type Object struct {
	value     interface{}
	ExpiresAt int64
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
