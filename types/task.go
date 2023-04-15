package types

type Task struct {
	ID         string
	Exec       Execute
	ExecParams []string
	Retry      int
	TResult    TaskResult
}

type TaskResult struct {
	ResultValue interface{}
	Err         Error
}

type Error struct {
	Code ErrorCode
	Name string
	Err  error
	Skip bool
}

type ErrorCode int

type Execute func(params ...string) TaskResult

func (t *Task) Do() TaskResult {
	result := t.Exec(t.ExecParams...)
	t.TResult = result
	return result
}
