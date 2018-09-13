package main

import (
	"os"
	"github.com/QPSTestTool/worker"
	"fmt"
	"time"
	"sync"
	task_package "github.com/QPSTestTool/task"
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
	TestFabric(paraNum, batchNum)
	return
}

func PressureTest(taskObj task_package.Task, paraNum int, batchNum int) (error){
	workers := make([]*worker.Worker, paraNum)
	results := make(chan *worker.Result, paraNum)
	for i:=0;i<paraNum;i++{
		workers[i] = worker.NewWorker(i, batchNum, taskObj, &results)
	}
	begin := time.Now()
	wg := sync.WaitGroup{}
	for idx, w := range workers{
		fmt.Printf("Start worker %d.\n", idx)
		wg.Add(1)
		go func(group *sync.WaitGroup) {
			w.DoTest()
			group.Done()
		}(&wg)
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


func TestFabric(paraNum int, batchNum int){
	if newTask, err := task_package.NewFabricTask(); err != nil{
		log.Fatal("create new task error. err=%+=v\n", err)
	} else {
		PressureTest(newTask, paraNum, batchNum)
	}
}