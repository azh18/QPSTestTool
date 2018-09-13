package task

type Task interface{
	Do() error
	GetTaskName() string
}


