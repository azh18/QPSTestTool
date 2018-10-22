package worker

import (
	"time"
	"github.com/zbw0046/QPSTestTool/task"
	"fmt"
)

type Result struct{
	ReqNum int
	SuccessNum int
	Latency []*time.Duration
}

type Worker struct{
	Id int
	BatchNum int
	Task task.Task
	Result *chan *Result
}

func NewWorker(id int, batchNum int, task task.Task, resChan *chan *Result) *Worker {
	return &Worker{
		Id:id,
		BatchNum: batchNum,
		Task: task,
		Result: resChan,
	}
}

func NewResult(reqNum int) *Result{
	return &Result{
		ReqNum: reqNum,
	}
}

func (r *Result) AddDuration(t *time.Duration){
	r.SuccessNum += 1
	r.Latency = append(r.Latency, t)
}

func (w *Worker) DoTest(){
	res := NewResult(w.BatchNum)
	for i:=0;i<w.BatchNum;i++{
		begin := time.Now()
		err := w.Task.Do()
		if err != nil{
			fmt.Printf("An error %+=v occured while run task.\n", err)
			continue
		}
		latency := time.Since(begin)
		res.AddDuration(&latency)
	}
	*w.Result <- res
}

