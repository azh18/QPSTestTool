# QPSTestTool

## Introduction

My first try to construct a universal framework for pressure test.

## Feature

You can directly use this framework to run pressure test on any type of tasks, such as HTTP service, RPC calls, and reading data from database. The only thing that you need to do is to write a task description in task module, in which you need to implement a ```Do()``` method and a ```GetTaskName()``` method, and then add a simple wrapper function in ```main.go```.

## Build

```go build main.go```

## Usage

```./main <parallelism> <batch number>```

## TODO

- support tasks with different request content

- simplify the process of defining new task type
