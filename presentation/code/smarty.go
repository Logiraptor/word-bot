package main

import (
	"runtime"
	"sync"
)

type job struct {
	//...
	results chan result
	wg *sync.WaitGroup
}

type result struct {
	//...
}

func consumeJobs(jobs chan job) {
	for job := range jobs {
		job.results <- perform(job)
		job.wg.Done()
	}
}

func enqueueJobs(jobChan chan job, resultChan chan result) {
	wg := new(sync.WaitGroup)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			wg.Add(1)
			jobChan <- job{wg: wg, results: resultChan /*...*/}
		}
	}
	wg.Wait()
	close(resultChan)
}

func workerPool() {
	jobChan := make(chan job)
	resultChan := make(chan result)
	for i := 0; i < runtime.NumCPU(); i++ {
		go consumeJobs(jobChan)
	}
	go enqueueJobs(jobChan, resultChan)
	for result := range resultChan {
		// do something with result
		_ = result
	}
}

func perform(j job) result {
	return result{}
}
