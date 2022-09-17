package linksSearcher

import "sync"

type ExecutionFn func(arg string) []string

type Job struct {
	Id     int
	Arg    string
	ExecFn ExecutionFn
}

func (j *Job) exec() Result {
	val := j.ExecFn(j.Arg)
	return Result{
		Id:    j.Id,
		Value: val,
	}
}

type Result struct {
	Id    int
	Value []string
}

type WorkerPool struct {
	workerCount int
	jobs        chan Job
	results     chan Result
	Done        chan struct{}
}

func NewWP(wcount int) *WorkerPool {
	return &WorkerPool{
		workerCount: wcount,
		jobs:        make(chan Job, wcount),
		results:     make(chan Result, wcount),
		Done:        make(chan struct{}),
	}
}

func (wp *WorkerPool) Run() {
	var wg sync.WaitGroup
	for i := 0; i < wp.workerCount; i++ {
		wg.Add(1)
		go worker(&wg, wp.jobs, wp.results)
	}
	wg.Wait()
	close(wp.results)
	close(wp.Done)
}

func (wp *WorkerPool) Result() chan Result {
	return wp.results
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) EndJob() {
	close(wp.jobs)
}

func worker(wg *sync.WaitGroup, jobs chan Job, results chan Result) {
	defer wg.Done()
	for job := range jobs {
		results <- job.exec()
	}
}
