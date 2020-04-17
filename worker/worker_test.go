package worker

import (
	"sync"
	"testing"
)

func TestWorker(t *testing.T) {
	done := make(chan Interface, 1)
	quit := make(chan struct{})
	wrk := New(done, quit)

	// Testing job's panic handling
	func() {
		var wg sync.WaitGroup
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Job's panic was caught: %q", r)
			}
		}()
		wg.Add(1)
		wrk.Run(func() { defer wg.Done(); panic("test panic") })
		wg.Wait()

		<-done // consume done signal
	}()

	// Testing worker shutdown
	const ExpectedPanicMessage = "send on closed channel"
	func() {
		defer func() {
			if r := recover(); r != nil {
				if err := r.(interface{ Error() string }).Error(); err != ExpectedPanicMessage {
					t.Fatalf("Expected %q panic, got: %q", ExpectedPanicMessage, err)
				}
				return
			}
			t.Fatalf("Expected panic after submitting job to finished worker")
		}()
		// Simulate the shutdown
		close(quit)
		<-done
		// Submit new job
		wrk.Run(func() { panic("should never happen") })
	}()
}

// vim: :ts=4:sw=4:ai:noexpandtab
