package worker

import (
	"runtime"
	"sync"
	"testing"
)

func TestWorker(t *testing.T) {
	done := make(chan Interface, 1)
	quit := make(chan struct{})
	w := New(done, quit)

	if _, ok := w.(*simpleWorker); !ok {
		t.Fatal("New worker is not of *simpleWorker type")
	}

	// Testing job's panic handling
	func() {
		var wg sync.WaitGroup
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Job's panic was caught: %q", r)
			}
		}()
		wg.Add(1)
		w.Run(func() { defer wg.Done(); panic("test panic") })
		wg.Wait()

		<-done // remove done worker
	}()

	// Testing worker shutdown
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r.(interface {
					Error() string
				}).Error() != "send on closed channel" {
					t.Fatalf("Expected %q panic, got: %q", "send on closed channel", r)
				}
			} else {
				t.Fatal("Expected panic after submitting job to finished worker")
			}
		}()

		// Simulate the shutdown
		close(quit)
		runtime.Gosched()

		// Submit new job
		w.Run(func() { return })
	}()
}
