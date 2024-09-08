package retry

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

var (
	DefaultAllowAttempt   uint32 = 3
	DefaultFactor         uint32 = 2
	DefaultMaxBackOffTime        = time.Second * 10 * 60
	ErrOverMaximumAttempt        = fmt.Errorf("over maximum attempt")
)

type (
	ExecuteFn func() error
	RetryFn   func(*Retry)
)

type Retry struct {
	factor         uint32
	allowAttempt   uint32
	currentAttempt uint32
	timer          *time.Timer
	maxBackOffTime time.Duration
}

func WithFactor(factor uint32) RetryFn {
	return func(retry *Retry) {
		retry.factor = factor
	}
}
func WithMaxBackOffTime(maxBackOffTime time.Duration) RetryFn {
	return func(retry *Retry) {
		retry.maxBackOffTime = maxBackOffTime
	}
}

func New(allowAttempt uint32, fns ...RetryFn) *Retry {
	if allowAttempt == 0 {
		allowAttempt = DefaultAllowAttempt
	}

	retry := &Retry{
		allowAttempt:   allowAttempt,
		factor:         DefaultFactor,
		timer:          time.NewTimer(0),
		maxBackOffTime: DefaultMaxBackOffTime,
	}
	for _, fn := range fns {
		fn(retry)
	}
	return retry
}

func (b *Retry) Duration() time.Duration {
	if atomic.LoadUint32(&b.currentAttempt) >= b.allowAttempt {
		return b.maxBackOffTime
	}
	atomic.AddUint32(&b.currentAttempt, 1)
	backDuration := time.Duration(math.Pow(float64(b.factor), float64(b.currentAttempt))) * time.Second
	if backDuration > b.maxBackOffTime {
		backDuration = b.maxBackOffTime
	}
	return backDuration
}

func (b *Retry) MaxBackOffTime() time.Duration {
	return b.maxBackOffTime
}

func (b *Retry) Reset() {
	atomic.SwapUint32(&b.currentAttempt, 0)
}

func (b *Retry) Run(fn ExecuteFn) error {
	if atomic.LoadUint32(&b.currentAttempt) >= atomic.LoadUint32(&b.allowAttempt) {
		return fmt.Errorf("over maximum attempt")
	}

	<-b.timer.C
	err := fn()
	if err == nil {
		b.timer.Stop()
		return nil
	}
	b.timer.Reset(b.Duration())
	if atomic.LoadUint32(&b.currentAttempt) == atomic.LoadUint32(&b.allowAttempt) {
		return err
	}
	return b.Run(fn)
}
