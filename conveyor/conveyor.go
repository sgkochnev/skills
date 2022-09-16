package conveyor

import (
	"sync"
)

func RunConveyor(tasksFlow ...task) {

	wg := &sync.WaitGroup{}
	var out chan interface{}

	for _, t := range tasksFlow {
		wg.Add(1)

		in := out
		out = make(chan interface{})

		go func(t task, in, out chan interface{}) {
			defer close(out)
			defer wg.Done()
			t(in, out)

		}(t, in, out)
	}

	wg.Wait()
}
