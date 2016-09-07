package worker

import (
	"log"

	"github.com/caelifer/scheduler/job"
)

// Worker is an interface type required for Scheduler to schedule its work load.
type Interface interface {
	Run(job.Interface)
}

// private implementation for simple worker
type simpleWorker struct {
	jobs chan job.Interface
	done chan<- Interface
}

// New constructs new Worker. It takes one parameter - a done channel to
// signal to Scheduler that this worker is done with its work.
func New(done chan<- Interface) Interface {
	w := new(simpleWorker) // Heap
	w.done = done
	w.jobs = make(chan job.Interface, 1)

	go func() {
		for {
			// Execute
			(<-w.jobs)()
		}
	}()

	return w
}

// Run method Implements Worker interface
func (w *simpleWorker) Run(job job.Interface) {
	w.jobs <- func() {
		defer func() {
			// Handle job's panic.
			if r := recover(); r != nil {
				log.Println("job failed:", r)
			}
			w.done <- w
		}()

		// Execute payload
		job()
	}
}
