package core

import "time"

func expirySample() float32 {
	var keyLimit int = 20
	var expiredCount = 0

	for key, Object := range store {
		if Object.ExpiresAt != -1 {
			keyLimit--

			if Object.ExpiresAt <= time.Now().UnixMilli() {
				delete(store, key)
				expiredCount++
			}
		}

		if keyLimit == 0 {
			break
		}
	}
	return float32(expiredCount) / float32(20)
}
