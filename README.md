# QPSTestTool

## Introduction

My first try to construct a universal framework for pressure test.

## Feature

- easy to add your own task

You can directly use this framework to run pressure test on any type of tasks, such as HTTP service, RPC calls, and reading data from database. The only thing that you need to do is to write a task description in task module, in which you need to implement some methods, and then add a simple wrapper function in ```main.go```.

- support tasks with different content

### Methods should be implemented

|Method|Description|
|---|---|
|```Do(ctx context.Context)```|The main procedure of executing the task|
|```Done() chan interface{}```|When a task is finished, you should send an item to the returned chan to notice the worker in the ```Do()``` method.|
|```GetTaskName() string```|You should return the name of your task.|
|```GetTimeout(ctx time.Duration)```|You should return the timeout you set in the task. When doing pressure test, tasks exceed the timeout will be seen as a failure.|
|```IsSuccessful() bool```|You should return whether this task is finished successfully at this time.|

### Actions should be added into ```main.go```

You should write a helper function that generate a list of tasks, and then push them into the parameters of the ```PressureTest``` function. Finally you should add your helper function in the ```main()``` function.

An example is the HttpTask in the code.

## Build

```go build main.go```

## Usage

```./main <parallelism> <batch number>```

## TODO

