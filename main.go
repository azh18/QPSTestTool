package main

import (
	"os"
	"github.com/QPSTestTool/worker"
	"fmt"
	"time"
	"sync"
	task_package "github.com/QPSTestTool/task"
	"strconv"
)

func main(){
	//task_package.NewFabricTask()
	//return
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
	TestFabric(paraNum, batchNum)
	return
}

func TestFabric(paraNum int, batchNum int) error {
	workers := make([]*worker.Worker, paraNum)
	results := make(chan *worker.Result, paraNum)
	task, _ := task_package.NewFabricTask()
	for i:=0;i<paraNum;i++{
		workers[i] = worker.NewWorker(i, batchNum, task, &results)
	}
	begin := time.Now()
	wg := sync.WaitGroup{}
	for idx, w := range workers{
		fmt.Printf("Start worker %d.\n", idx)
		wg.Add(1)
		go func(group sync.WaitGroup) {
			w.DoTest()
			}(wg)
	}
	wg.Wait()
	duration := time.Since(begin)
	nTotalReq := batchNum * paraNum
	totalLatency := time.Duration(0.0)
	nSuccessReq := 0
	var qps float64
	for i:=0;i<paraNum;i++{
		res := <- results
		nSuccessReq += res.SuccessNum
		for _, l := range res.Latency{
			totalLatency += *l
		}
	}
	qps = float64(nSuccessReq) / float64(duration.Nanoseconds())
	avgLatency := float64(totalLatency) / float64(nSuccessReq)
	fmt.Printf("Test finished. QPS=%v, Average Latency=%v ms, success resp = %v, fail resp = %v",
		qps, avgLatency, nSuccessReq, nTotalReq-nSuccessReq)
	return nil
}
