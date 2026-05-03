package core

import "github.com/adityjoshi/iDB/config"

func evict() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()
	}
}
