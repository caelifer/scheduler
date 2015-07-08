scheduler
=========

`scheduler` is a simple load-balancer across n-workers. It maintains a buffered channel of work-uinits (github.com/caelifer/scheduler/jobs.Interface)

# Installation
```
go get github.com/caelifer/schduler
```

# Usage
```
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/caelifer/scheduler"
)

const (
	NWORKERS = 10
	NJOBS    = 5
	NSAMPS   = 100
)

func main() {
	wg := new(sync.WaitGroup)
	sch := scheduler.New(NWORKERS, NJOBS)

	t0 := time.Now()
	rand.Seed(t0.UnixNano())

	for i := 0; i < NSAMPS; i++ {
		wg.Add(1)
		sch.Schedule(func() {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		})
	}

	wg.Wait()
	fmt.Printf("Ran %d samples in %s\n", NSAMPS, time.Since(t0))
}
```
This code should produce output similar to
```
$ go run test.go
Ran 100 samples in 546.21049ms
```
