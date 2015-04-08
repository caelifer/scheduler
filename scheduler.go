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
// on the background. This scheduling process makes sure that all available workers
// are always running on the background waiting for the work unit to come in. The
// work units are managed in a separate buffered channel of job.Interfaces. New
// takes two paramters.  nworkers - a number of background workers, and njobs -
// a number of queued jobs, before scheduler would block.
func New(nworkers, njobs int) Scheduler {
	s := new(simpleScheduler) // Heap
	s.workPool = make(chan worker.Interface, nworkers)
	s.done = make(chan worker.Interface)
	s.jobs = make(chan job.Interface, njobs)

	// Populate our pool of workers
	for i := 0; i < nworkers; i++ {
		s.workPool <- worker.New(s.done)
	}

	// Start our scheduler on the background
	go func() {
		// Run until the main program finishes
		for {
			select {
			case w := <-s.workPool:
				// Schedule new worker on the background "thread"
				go func() { w.Run(<-s.jobs) }()
			case w := <-s.done:
				s.workPool <- w // Never blocks
			}
		}
	}()

	return s
}

// Private concreat implementation

// simpleScheduler is an value type that implementes Scheduler interface.
type simpleScheduler struct {
	workPool chan worker.Interface // Buffered channel of workers
	done     chan worker.Interface // Unbuffered channel for workers to signal they are done
	jobs     chan job.Interface    // Buffered channel of pending work units
}

// Schedule is an implementation of Schedlue interface for simpleSchedule value type.
func (s *simpleScheduler) Schedule(j job.Interface) {
	s.jobs <- j // Could block if jobs buffer is full
}
