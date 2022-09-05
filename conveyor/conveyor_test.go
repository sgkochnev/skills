package conveyor

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const (
	test1Count      = 3
	test1Sleep      = 100 * time.Millisecond
	test1SleepDelta = 50 * time.Millisecond
	test1Result     = "gggooo!!!"
)

func TestConveyor1(t *testing.T) {
	var resultSentence atomic.Value
	resultSentence.Store("")
	tasksFlow := []task{
		task(func(in, out chan interface{}) {
			out <- "g"
			out <- "o"
			out <- "!"
		}),
		task(func(in, out chan interface{}) {
			for val := range in {
				out <- strings.Repeat(val.(string), test1Count)
				time.Sleep(test1Sleep)
			}
		}),
		task(func(in, out chan interface{}) {
			for val := range in {
				resultSentence.Store(resultSentence.Load().(string) + val.(string))
			}
		}),
	}

	start := time.Now()

	RunConveyor(tasksFlow...)

	end := time.Since(start)

	expectedTime := test1Sleep*test1Count + test1SleepDelta

	if end > expectedTime {
		t.Errorf("Execution took too long. Got: %s. Expected: < %s", end, expectedTime)
	}
	res := resultSentence.Load().(string)
	if res != test1Result {
		t.Errorf("Last task have not collected inputs. Got: %s. Expected: %s", res, test1Result)
	}
}

const (
	test2Sleep = 10 * time.Millisecond
)

func TestConveyor2(t *testing.T) {
	var result, resultInFirstTask uint32
	tasksFlow := []task{
		task(func(in, out chan interface{}) {
			out <- 1
			time.Sleep(test2Sleep)
			resultInFirstTask = atomic.LoadUint32(&result)
		}),
		task(func(in, out chan interface{}) {
			for range in {
				atomic.AddUint32(&result, 1)
			}
		}),
	}
	RunConveyor(tasksFlow...)
	if result == 0 || resultInFirstTask == 0 {
		t.Errorf("Conveyour isn't working properly, flow of values mustn't depend on tasks execution")
	}
}
