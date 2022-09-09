package conveyor

import (
	"sync"
)

func RunConveyor(tasksFlow ...task) {

	wg := sync.WaitGroup{}
	var chOut chan interface{}
	for i, t := range tasksFlow {

		wg.Add(1)

		switch i {
		case 0:
			chOut = make(chan interface{})
			go func(t task, out chan interface{}) {
				t(nil, out)
				close(out)
				wg.Done()
			}(t, chOut)

		case len(tasksFlow) - 1:
			chIn := fanIn(chOut)
			go func(t task, in chan interface{}) {
				t(in, nil)
				wg.Done()
			}(t, chIn)

		default:
			chIn := fanIn(chOut)
			chOut = make(chan interface{})
			go func(t task, in, out chan interface{}) {
				t(in, out)
				close(out)
				wg.Done()
			}(t, chIn, chOut)
		}
	}
	wg.Wait()
}

func fanIn(chIn chan interface{}) chan interface{} {
	ch := make(chan interface{})
	go func() {
		for v := range chIn {
			ch <- v
		}
		close(ch)
	}()
	return ch
}
