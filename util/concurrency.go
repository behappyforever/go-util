package util

import (
	"context"
	"github.com/golang/protobuf/ptypes/any"
	"golang.org/x/time/rate"
	"log"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type fun func(context.Context, *any.Any) error

type callback func(context.Context, int64)

// anyArr: all pending data
// f: handle per any
// concurrency: goroutine num
// cb: record success num
// permitsPerSecond: tokens per second
func AsyncParallelExecutionWithLimit(ctx context.Context, anyArr []*any.Any, f fun, concurrency int, cb callback, permitsPerSecond int) {
	if len(anyArr) == 0 {
		return
	}
	go func() {
		defer Recovery()
		start := time.Now()
		l := NewLimiter(concurrency)
		rateLimit := rate.Inf
		if permitsPerSecond > 0 {
			rateLimit = rate.Limit(permitsPerSecond)
		}
		r := rate.NewLimiter(rateLimit, 1)
		var count int64
		var wg sync.WaitGroup
		wg.Add(len(anyArr))
		for _, e := range anyArr {
			l.Begin()
			go func(e *any.Any) {
				defer Recovery()
				defer l.End()
				defer wg.Done()
				r.Wait(ctx)
				if err := f(ctx, e); err == nil {
					atomic.AddInt64(&count, 1)
				}
			}(e)
		}
		wg.Wait()
		if cb != nil {
			cb(ctx, count)
		}
		log.Printf( "%s elapsed time, %s, totalNum:%d, successNum:%d", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), time.Since(start).String(), len(anyArr), count)
	}()
}

// shorthand for AsyncParallelExecutionWithLimit(..., 0)
func AsyncParallelExecution(ctx context.Context, anyArr []*any.Any, f fun, cb callback, concurrency int) {
	AsyncParallelExecutionWithLimit(ctx, anyArr, f, concurrency, cb, 0)
}
