scheduler
=========

`scheduler` is a simple load-balancer across n-workers. It maintains a buffered channel of work-uinits (github.com/caelifer/scheduler.jobs.Interface)

# Installation
```
go get github.com/caelifer/schduler
```

# Usage
```
import "github.com/caelifer/scheduler"

const (
  NWORKERS = 50
  NJOBS = 100
)

...
sched := scheduler.New(NWORKERS, NJOBS)
out := make(chan Result)

...
sched.Schedule(func(){
  out <- longRunningTask()
})

...
res := <-out
```
