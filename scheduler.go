package scheduler

import (
	"runtime"

	"github.com/caelifer/scheduler/job"
	"github.com/caelifer/scheduler/worker"
)

// Scheduler is an interface type to provide an abstraction around load-balancing operation scheduling.
type Scheduler interface {
	Schedule(job.Interface)
	Shutdown()
}

// New builds a new Scheduler object. It starts its internal scheduling process
// in the background. This scheduling process makes sure that all available workers
// are always running in the background waiting for the work unit to come in. The
// work units are managed in a separate buffered channel of job.Interfaces. New
// takes two paramters: nworkers - a number of background workers, and njobs -
// a number of queued jobs, before scheduler would block.
func New(nworkers, njobs int) Scheduler {
	s := new(simpleScheduler) // Heap allocation
	s.wpool = make(chan worker.Interface, nworkers)
	s.jpool = make(chan job.Interface, njobs)
	s.quit = make(chan struct{})

	// Populate our pool of workers
	for i := 0; i < nworkers; i++ {
		s.wpool <- worker.New(s.wpool, s.quit)
	}
	// Start our scheduler on the background
	go func() {
		// Run until quit signal is received
		for {
			select {
			case w := <-s.wpool:
				// Run next job with next available worker. This will block when either
				// there are no available workers or jobs queue is empty
				if j, ok := <-s.jpool; ok {
					w.Run(j)
				}
			case <-s.quit:
				return
			}
		}
	}()
	return s
}

// simpleScheduler is an value type that implements Scheduler interface.
type simpleScheduler struct {
	wpool chan worker.Interface // Buffered channel of workers
	jpool chan job.Interface    // Buffered channel of pending work units
	quit  chan struct{}         // Quit channel for orderly shutdown
}

// Schedule is an implementation of Schedule interface for simpleSchedule value type.
func (s *simpleScheduler) Schedule(j job.Interface) {
	s.jpool <- j // Could block if jobs buffer is full
	runtime.Gosched()
}

func (s *simpleScheduler) Shutdown() {
	close(s.quit)
	close(s.jpool)
	close(s.wpool)
	s = nil // release memory
}
