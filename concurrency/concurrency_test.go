package concurrency

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestWithContext(t *testing.T) {
	start := time.Now()

	var v int64
	g, _ := WithContext(context.Background(), 10)
	for i := 0; i < 100; i++ {
		v2 := i
		g.Go(func() error {
			atomic.AddInt64(&v, int64(v2))
			if v2 == 49 {
				return errors.New("this is a test error")
			}
			time.Sleep(time.Second)
			return nil
		})
	}

	e := g.Wait()

	t.Logf("time elapsed: %s", time.Since(start))
	t.Log(v)
	t.Log(e)
}
