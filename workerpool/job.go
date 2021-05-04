package workerpool

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

func (j *Job) Run() {
	err, result := j.executeTask()

	j.Err = err
	j.resultChan <- Result{
		Value: result,
		JobId: j.Id,
	}
}

func NewJob(id int, resultChan chan Result, fn func() (error, string)) Job {
	return Job{
		executeTask: fn,
		Id:          id,
		resultChan:  resultChan,
	}
}
