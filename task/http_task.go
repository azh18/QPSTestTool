package task

import "net/http"

type HttpTask struct {
	taskName string
	cli *http.Client
}

func NewHttpTask() *HttpTask{
	return nil
}

func (t *HttpTask) Do() error {
	t.cli.Get("http://www.baidu.com")
	return nil
}

func (t *HttpTask) GetTaskName() string{
	return t.taskName
}

