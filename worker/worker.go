package worker

import "github.com/caelifer/scheduler/job"

// Worker is an interface type required for Scheduler to schedule its work load.
type Interface interface {
	Run(job.Interface)
}

// private implementation for simple worker
type simpleWorker struct {
	done chan<- Interface
}

// New constructs new Worker. It takes one parameter - a done channel to
// signal to Scheduler that this worker is done with its work.
func New(done chan<- Interface) Interface {
	w := new(simpleWorker) // Heap
	w.done = done
	return w
}

// Run method Implements Worker interface
func (w *simpleWorker) Run(j job.Interface) {
	defer func() { w.done <- w }() // TODO handle job's panic
	j()
}
