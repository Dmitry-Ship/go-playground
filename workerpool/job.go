package workerpool

import (
	"math/rand"
	"strconv"
	"time"
)

type Result struct {
	Value string
	JobId int
}

type Job struct {
	Err         error
	Id          int
	executeTask func() (error, string)
	resultChan  chan Result
}

func performLongWork() (error, string) {
	randomNumber := rand.Intn(1000)

	time.Sleep(time.Duration(randomNumber) * time.Millisecond)
	return nil, strconv.Itoa(randomNumber)
}

func (j *Job) Run() {
	err, result := j.executeTask()

	j.Err = err
	j.resultChan <- Result{
		Value: result,
		JobId: j.Id,
	}
}

func NewJob(id int, resultChan chan Result) *Job {
	return &Job{
		executeTask: performLongWork,
		Id:          id,
		resultChan:  resultChan,
	}
}
