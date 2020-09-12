package kronika

import (
	"context"
	"testing"
	"time"
)

func TestWaitFor(t *testing.T) {
	interval := 200 * time.Millisecond
	ctx := context.Background()

	t1 := time.Now()
	WaitFor(ctx, interval)
	t2 := time.Now()

	actualInterval := t2.Sub(t1)
	if actualInterval < interval {
		t.Fatalf("Difference between times should be more than %d, but was %d", interval, actualInterval)
	}
}

func TestWaitForCancel(t *testing.T) {
	interval := 200 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	t1 := time.Now()
	WaitFor(ctx, interval) // this should return immediately and so not reach our interval
	t2 := time.Now()

	actualInterval := t2.Sub(t1)
	if actualInterval > interval {
		t.Fatalf("Difference between times should be less than %d, but was %d", interval, actualInterval)
	}
}

func TestWaitUntil(t *testing.T) {
	interval := 200 * time.Millisecond
	ctx := context.Background()

	t1 := time.Now().Add(interval)
	WaitUntil(ctx, t1)
	t2 := time.Now()

	actualInterval := t2.Sub(t1)
	if actualInterval < 0 {
		t.Fatalf("t2 should be after t1 but was not. t2: %d, t1: %d", t2.UnixNano(), t1.UnixNano())
	}
}

func TestWaitUntilCancel(t *testing.T) {
	interval := 200 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	t1 := time.Now().Add(interval)
	WaitUntil(ctx, t1) // this should return immediately and so not reach t1
	t2 := time.Now()

	actualInterval := t2.Sub(t1)
	if actualInterval > 0 {
		t.Fatalf("t2 should be before t1 but was not. t2: %d, t1: %d", t2.UnixNano(), t1.UnixNano())
	}
}

func TestEvery(t *testing.T) {
	interval := 200 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel whenever it exits

	lastT := time.Now()
	start := lastT.Add(interval)

	i := 0
	noOfTicks := 10

	for currT := range Every(ctx, start, interval) {
		if i == noOfTicks {
			cancel() // cancel the context
		} else {
			i++
		}

		diff := currT.Sub(lastT)
		slop := diff * 2 / 10 // variation allowance see https://golang.org/src/time/tick_test.go

		if diff-slop > interval || diff+slop < interval {
			t.Fatalf("Interval since last tick should be at least %d but was %d.", interval, diff)
		}

		lastT = currT
	}

	if i != noOfTicks {
		t.Fatalf("Ticker should have run %d times but ran %d times.", noOfTicks, i)
	}
}

func TestEveryCancelBeforeFistEvent(t *testing.T) {
	interval := 5 * time.Second
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	start := time.Now()

	for range Every(ctx, start, interval) {
	}

	length := time.Now().Sub(start)
	if length > time.Millisecond*10 {
		t.Fatalf(
			"Ticker should have stopped almost immediately, but took %dms",
			length/time.Millisecond)
	}
}
