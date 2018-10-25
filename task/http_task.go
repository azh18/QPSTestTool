package task

import (
	"context"
	"net/http"
	"time"
)

type HttpTask struct {
	taskName string
	url string
	timeout time.Duration
	status int
	retChan chan interface{}
}

func NewHttpTasks(timeout time.Duration) []*HttpTask{
	urlPool := []string{
		"http://www.baidu.com",
		"http://www.163.com",
		"http://www.qq.com",
		"http://www.taobao.com",
	}
	ret := make([]*HttpTask, 0)
	for _, url := range urlPool{
		httpTask := &HttpTask{
			taskName: "fetch:" + url,
			url: url,
			timeout: timeout,
			status: -1,
			retChan: make(chan interface{}, 1),
		}
		ret = append(ret, httpTask)
	}
	return ret
}

func (t *HttpTask) Do(ctx context.Context)  {
	url := t.url
	resp, err := http.Get(url)
	if err != nil{
		t.status = -1
	} else {
		t.status = resp.StatusCode
	}
	t.retChan <- 1
}

func (t *HttpTask) GetTaskName() string{
	return t.taskName
}

func (t *HttpTask) IsSuccessful() bool {
	if t.status == 200{
		return true
	} else {
		return false
	}
}

func (t *HttpTask) GetTimeout() time.Duration{
	return t.timeout
}

func (t *HttpTask) Done() chan interface{}{
	return t.retChan
}
