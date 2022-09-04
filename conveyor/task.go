package conveyor

type task func(in, out chan interface{})
