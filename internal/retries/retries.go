package retries

import (
	"math/rand"
	"time"
)

func Interval(attempts int) time.Duration {

	d := time.Second * time.Duration(1<<attempts)
	max := 30 * time.Minute
	if d > max {
		d = max
	}
	jitter := time.Duration(rand.Int63n(int64(time.Second)))
	return d + jitter
}
