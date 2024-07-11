package timers

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	interval struct {
		running  bool
		duration time.Duration
		until    time.Time
		channel  workflow.Channel
	}

	// Interval provides helpers to manage a recurring interval in the context of a temporal workflow.
	Interval interface {
		// Next blocks until the the end of the interval. After that, it prepares the interval for the next iteration.
		Next(ctx workflow.Context)

		// NextWith blocks until the the end of the interval. After that, it prepares the interval for the next iteration
		// with the specified duration.
		NextWith(ctx workflow.Context, duration time.Duration)

		// Restart restarts the current interval.
		Restart(ctx workflow.Context)

		// RestartWith restarts the current interval with the specified duration.
		RestartWith(ctx workflow.Context, duration time.Duration)

		// Cancel stops the current interval.
		Cancel(ctx workflow.Context)
	}
)

// NextWith blocks until the the end of the interval. After that, it prepares the interval for the next iteration
// with the specified duration.
func (t *interval) NextWith(ctx workflow.Context, duration time.Duration) {
	t.running = true
	t.wait(ctx)
	t.update(ctx, duration)
	t.running = false
}

// RestartWith stops the current interval and starts a new one with the specified duration.
func (t *interval) RestartWith(ctx workflow.Context, duration time.Duration) {
	if t.running {
		t.channel.Send(ctx, duration)
	} else {
		t.update(ctx, duration)
	}
}

// Next blocks until the the end of the interval. After that, it prepares the interval for the next iteration.
func (t *interval) Next(ctx workflow.Context) {
	t.NextWith(ctx, t.duration)
}

func (t *interval) Restart(ctx workflow.Context) {
	t.RestartWith(ctx, t.duration)
}

func (t *interval) Cancel(ctx workflow.Context) {
	t.channel.Send(ctx, time.Duration(0))
}

// wait blocks until the timer expires or a message is received on the channel. The timer is cancelled if the duration is 0,
// otherwise it is reset.
func (t *interval) wait(ctx workflow.Context) {
	done := false

	for !done && ctx.Err() == nil {
		_ctx, cancel := workflow.WithCancel(ctx)
		duration := time.Duration(0)
		timer := workflow.NewTimer(_ctx, t.duration)
		selector := workflow.NewSelector(_ctx)

		// when the channel receives a message
		selector.AddReceive(t.channel, func(channel workflow.ReceiveChannel, more bool) {
			channel.Receive(_ctx, &duration)
			cancel()

			if duration == 0 {
				done = true // the timer is cancelled, so the interval is over.
			} else {
				t.update(_ctx, t.duration)
			}
		})

		// when the timer finishes
		selector.AddFuture(timer, func(future workflow.Future) {
			if err := future.Get(_ctx, nil); err == nil {
				done = true
			}
		})

		selector.Select(ctx)
	}
}

// update updates the interval's duration and the time at which the interval should stop.
// The duration parameter specifies the new interval duration.
func (t *interval) update(ctx workflow.Context, duration time.Duration) {
	t.duration = duration
	t.until = Now(ctx).Add(duration)
}

func NewInterval(ctx workflow.Context, duration time.Duration) Interval {
	return &interval{
		duration: duration,
		until:    Now(ctx).Add(duration),
		channel:  workflow.NewChannel(ctx),
	}
}