package types

type TaskBuilder interface {
	GetNextTask(previousTask *Task) *Task
	TaskCount() int
	Result() *JobResult
}

type Job struct {
	ID    string
	Retry int
	Tasks TaskBuilder
}

type JobResult struct {
	Finished bool
	Err      Error
	Value    string
}

func NewJob(id string, tasks TaskBuilder) *Job {
	return &Job{
		ID:    id,
		Tasks: tasks,
	}
}
