package scheduler

import (
	"github.com/caelifer/scheduler/job"
	"github.com/caelifer/scheduler/worker"
)

// Public interfaces

// Scheduler is an interface type to provide an abstraction around load-balancing operation scheduling.
type Scheduler interface {
	Schedule(job.Interface)
}

// New builds a new Scheduler object. It starts its internal scheduling process
// in the background. This scheduling process makes sure that all available workers
// are always running in the background waiting for the work unit to come in. The
// work units are managed in a separate buffered channel of job.Interfaces. New
// takes two paramters: nworkers - a number of background workers, and njobs -
// a number of queued jobs, before scheduler would block.
func New(nworkers, njobs int) Scheduler {
	s := new(simpleScheduler) // Heap
	s.workPool = make(chan worker.Interface, nworkers)
	s.jobs = make(chan job.Interface, njobs)

	// Populate our pool of workers
	for i := 0; i < nworkers; i++ {
		s.workPool <- worker.New(s.workPool)
	}

	// Start our scheduler on the background
	go func() {
		// Run until the main program finishes
		for {
			// Get new worker from the worker pool
			w := <-s.workPool

			// Run next job with this worker
			w.Run(<-s.jobs) // should not block
		}
	}()

	return s
}

// Private concrete implementation

// simpleScheduler is an value type that implements Scheduler interface.
type simpleScheduler struct {
	workPool chan worker.Interface // Buffered channel of workers
	jobs     chan job.Interface    // Buffered channel of pending work units
}

// Schedule is an implementation of Schedule interface for simpleSchedule value type.
func (s *simpleScheduler) Schedule(j job.Interface) {
	s.jobs <- j // Could block if jobs buffer is full
}
