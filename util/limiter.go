package util

type limiter struct {
	ch chan struct{}
}

func NewLimiter(limit int) *limiter {
	return &limiter{make(chan struct{}, limit)}
}

func (l *limiter) Begin() {
	l.ch <- struct{}{}
}

func (l *limiter) End() {
	<-l.ch
}
