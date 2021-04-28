package workerpool

import (
	"math/rand"
	"strconv"
	"time"
)

type Task struct {
	Err         error
	Id          int
	executeTask func() (error, string)
	resultsChan chan *Result
}

type Result struct {
	Value  string
	TaskId int
}

func performLongWork() (error, string) {
	randomNumber := rand.Intn(10000)

	time.Sleep(time.Duration(randomNumber) * time.Millisecond)
	return nil, strconv.Itoa(randomNumber)
}

func (t *Task) Run() {
	err, result := t.executeTask()
	t.Err = err

	formattedResult := Result{
		Value:  result,
		TaskId: t.Id,
	}

	t.resultsChan <- &formattedResult
}

func NewTask(id int, resultsChan chan *Result) *Task {
	return &Task{
		executeTask: performLongWork,
		Id:          id,
		resultsChan: resultsChan,
	}
}
