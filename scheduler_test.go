package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	var tests = []struct {
		workers int
		results []int
	}{
		{1 << 0, make([]int, 100000)},
		{1 << 1, make([]int, 100000)},
		{1 << 3, make([]int, 100000)},
		{1 << 4, make([]int, 100000)},
		{1 << 5, make([]int, 100000)},
		{1 << 6, make([]int, 100000)},
		{1 << 7, make([]int, 100000)},
		{1 << 8, make([]int, 100000)},
		{1 << 9, make([]int, 100000)},
	}

	for _, tc := range tests {
		var wg sync.WaitGroup
		N := len(tc.results)

		t0 := time.Now()
		sch := New(tc.workers)
		makeJob := func(max, i int, res []int) func() {
			return func() {
				defer wg.Done()
				// Make sure number of running workers is less or eaqual to
				sch := sch.(*simpleScheduler)
				if n := len(sch.wpool) + 1; n > max {
					t.Fatalf("Number of runing workers is greater then max specified: %d > %d", n, max)
				}
				res[i] = i
			}
		}

		// Schedule and execute all pending work
		wg.Add(N)
		for i := 0; i < N; i++ {
			sch.Schedule(makeJob(tc.workers, i, tc.results))
		}
		wg.Wait()

		// Check work consistency
		for i, p := 1, tc.results[0]; i < len(tc.results); i++ {
			if tc.results[i]-p != 1 {
				t.Errorf("Wrong number at index %d; wanted %d, got %d", i, p+1, tc.results[i])
			}
			p = tc.results[i]
		}
		t.Logf("Scheduling %d jobs to completion with %3d workers took %v", N, tc.workers, time.Since(t0))

		// Shutdown
		sch.Shutdown()
	}
}

func BenchmarkScheduler1Worker(b *testing.B) {
	benchmarkScheduler(b, New(1))
}

func BenchmarkScheduler10Workers(b *testing.B) {
	benchmarkScheduler(b, New(10))
}

func BenchmarkScheduler50Workers(b *testing.B) {
	benchmarkScheduler(b, New(50))
}

func BenchmarkScheduler100Workers(b *testing.B) {
	benchmarkScheduler(b, New(100))
}

func benchmarkScheduler(b *testing.B, sch Scheduler) {
	var wg sync.WaitGroup
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		sch.Schedule(func() {
			defer wg.Done()
			// NOOP
		})
	}
	wg.Wait()
	sch.Shutdown()
}

// vim: :ts=4:sw=4:autoindent:noexpandtab
