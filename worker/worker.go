package worker

import (
	"context"
	"github.com/prometheus/common/log"
	"github.com/zbw0046/QPSTestTool/task"
	"math/rand"
	"time"
)

type Result struct {
	ReqNum     int
	SuccessNum int
	Latency    []*time.Duration
}

type Worker struct {
	Id       int
	BatchNum int
	Tasks    []task.Task
	Result   chan *Result
}

func NewWorker(id int, batchNum int, taskPool []task.Task, resChan chan *Result) *Worker {
	return &Worker{
		Id:       id,
		BatchNum: batchNum,
		Tasks:    taskPool, // random choose task
		Result:   resChan,
	}
}

func NewResult(reqNum int) *Result {
	return &Result{
		ReqNum: reqNum,
	}
}

func (r *Result) AddDuration(t *time.Duration) {
	r.SuccessNum += 1
	r.Latency = append(r.Latency, t)
}

func (w *Worker) DoTest(ctx context.Context) {
	res := NewResult(w.BatchNum)
	for i := 0; i < w.BatchNum; i++ {
		begin := time.Now()
		taskToExecute := w.Tasks[rand.Intn(len(w.Tasks))]
		thisCtx, _ := context.WithTimeout(ctx, taskToExecute.GetTimeout())
		go taskToExecute.Do(thisCtx)
		select {
		case <-thisCtx.Done():
			log.Info("Task Timeout.\n")
			break
		case <-taskToExecute.Done():
			if status := taskToExecute.IsSuccessful(); status {
				latency := time.Since(begin)
				res.AddDuration(&latency)
			}
			break
		}
	}
	w.Result <- res
}
