package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	var tests = []struct {
		workers, jobs int
		results       []int
	}{
		{1, 1, make([]int, 10000)},
		{10, 1, make([]int, 10000)},
		{10, 10, make([]int, 10000)},
		{100, 10, make([]int, 10000)},
		{100, 100, make([]int, 10000)},
		{100, 500, make([]int, 10000)},
	}

	for _, tc := range tests {
		var wg sync.WaitGroup
		N := len(tc.results)

		t0 := time.Now()
		sch := New(tc.workers, tc.jobs)
		makeJob := func(max, i int, res []int) func() {
			return func() {
				defer wg.Done()
				// Make sure number of running workers is less or eaqual to
				sch := sch.(*simpleScheduler)
				if n := len(sch.workPool) + 1; n > max {
					t.Fatalf("Number of runing workers is greater then max specified: %d > %d", n, max)
				}
				res[i] = i
			}
		}

		wg.Add(N)
		for i := 0; i < N; i++ {
			sch.Schedule(makeJob(tc.workers, i, tc.results))
		}
		wg.Wait()

		for i, p := 1, tc.results[0]; i < len(tc.results); i++ {
			if tc.results[i]-p != 1 {
				t.Errorf("Wrong number at index %d; wanted %d, got %d", i, p+1, tc.results[i])
			}
			p = tc.results[i]
		}
		t.Logf("Scheduling and completion of %d jobs with %3d workers and %3d job queue took %v", N, tc.workers, tc.jobs, time.Since(t0))
	}
}

func BenchmarkScheduler1over10(b *testing.B) {
	benchmarkScheduler(b, New(1, 10))
}

func BenchmarkScheduler5over10(b *testing.B) {
	benchmarkScheduler(b, New(5, 10))
}

func BenchmarkScheduler10over10(b *testing.B) {
	benchmarkScheduler(b, New(10, 10))
}

func BenchmarkScheduler1over100(b *testing.B) {
	benchmarkScheduler(b, New(1, 100))
}

func BenchmarkScheduler50over100(b *testing.B) {
	benchmarkScheduler(b, New(50, 100))
}

func BenchmarkScheduler100over10(b *testing.B) {
	benchmarkScheduler(b, New(100, 10))
}

func BenchmarkScheduler100over100(b *testing.B) {
	benchmarkScheduler(b, New(100, 100))
}

func benchmarkScheduler(b *testing.B, sch Scheduler) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		sch.Schedule(func() {
			defer wg.Done()
			time.Sleep(10 * time.Microsecond)
		})
	}
	wg.Wait()
}
