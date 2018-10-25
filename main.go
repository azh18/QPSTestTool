package main

import (
	"context"
	"os"
	"github.com/zbw0046/QPSTestTool/worker"
	"fmt"
	"time"
	"sync"
	task_package "github.com/zbw0046/QPSTestTool/task"
	"strconv"
	"github.com/prometheus/common/log"
)

func main(){
	paraNum, err := strconv.Atoi(os.Args[1]) // 并发度
	if err != nil{
		fmt.Printf("The first arg (parallelism) is not an integer!\n")
		return
	}
	batchNum, err := strconv.Atoi(os.Args[2]) // 每次请求量
	if err != nil{
		fmt.Printf("The second arg (requests num) is not an integer!\n")
		return
	}
	TestHttp(paraNum, batchNum)
	return
}

// random select task from taskObjs, to achieve the effect of random fetch
func PressureTest(taskObjs []task_package.Task, paraNum int, batchNum int) (error){
	workers := make([]*worker.Worker, paraNum)
	results := make(chan *worker.Result, paraNum)
	for i:=0;i<paraNum;i++{
		workers[i] = worker.NewWorker(i, batchNum, taskObjs, results)
	}
	begin := time.Now()
	wg := sync.WaitGroup{}
	parentContext := context.Background()
	for idx, w := range workers{
		fmt.Printf("Start worker %d.\n", idx)
		wg.Add(1)
		ctx, _ := context.WithTimeout(parentContext, time.Second * 30)
		go func(ctx context.Context, group *sync.WaitGroup) {
			w.DoTest(ctx)
			group.Done()
		}(ctx, &wg)
	}
	wg.Wait()

	var qps float64
	duration := time.Since(begin)

	nTotalReq := batchNum * paraNum
	totalLatency := time.Duration(0.0)
	nSuccessReq := 0
	for i:=0;i<paraNum;i++{
		res := <- results
		fmt.Printf("get one result.\n")
		nSuccessReq += res.SuccessNum
		for _, l := range res.Latency{
			totalLatency += *l
		}
	}

	qps = float64(nSuccessReq) / float64(duration.Seconds())
	avgLatency := float64(totalLatency.Nanoseconds()/(10E6)) / float64(nSuccessReq)

	fmt.Printf("Test finished. QPS=%v, Average Latency=%v ms, success resp = %v, fail resp = %v\n",
		qps, avgLatency, nSuccessReq, nTotalReq-nSuccessReq)
	return nil
}


func TestHttp(paraNum int, batchNum int) {
	if newTasks := task_package.NewHttpTasks(time.Millisecond * 150); len(newTasks) == 0{
		log.Fatal("create new http task error.\n")
	} else {
		// TODO: add a helper function to automatically transform the slice of user-defined task to the slice of task interface
		httpTasks := make([]task_package.Task, 0)
		for _, t := range newTasks {
			httpTasks = append(httpTasks, t)
		}
		PressureTest(httpTasks, paraNum, batchNum)
	}
}