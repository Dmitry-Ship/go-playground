package workerpool

type Result struct {
	Value string
	JobId int
}

type Job struct {
	Err         error
	Id          int
	executeTask func() (string, error)
	resultChan  chan Result
}

func NewJob(id int, resultChan chan Result, fn func() (string, error)) Job {
	return Job{
		executeTask: fn,
		Id:          id,
		resultChan:  resultChan,
	}
}

func (j *Job) Run() {
	result, err := j.executeTask()

	j.Err = err
	j.resultChan <- Result{
		Value: result,
		JobId: j.Id,
	}
}
