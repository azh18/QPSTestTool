package task

import (
	"context"
	"time"
)

type Task interface{
	Do(ctx context.Context)
	Done() chan interface{} // in Do(), if task is finished, should send a msg to this channel
	GetTaskName() string
	GetTimeout() time.Duration
	IsSuccessful() bool // to mark whether the task is finished successfully.
}