package util

import (
	"context"
	"fmt"
)

type Task interface {
	Execute()
}

// Double queue is a priority queue with two priority.
// The tasks within low-priority queue would be done when there is no task in high-priority queue
type DoubleQueue struct {
	highChan chan Task
	lowChan chan Task
	cacheSize int

	ctx context.Context
}

func NewDoubleQueue(ctx context.Context, cacheSize int) *DoubleQueue {
	return &DoubleQueue{
		highChan:  make(chan Task, cacheSize),
		lowChan:   make(chan Task, cacheSize),
		cacheSize: cacheSize,
		ctx:       ctx,
	}
}

func (q *DoubleQueue) AddHigh(t Task) {
	q.highChan <- t
}

func (q *DoubleQueue) AddLow(t Task) {
	q.lowChan <- t
}

func (q *DoubleQueue) Start(withWarning bool) {
	var warnSize int = q.cacheSize/2
	var task Task
	for {
		if withWarning {
			if len(q.highChan) > warnSize {
				fmt.Println("WARN: DoubleQueue highChan half full.")
			}
			if len(q.lowChan) > warnSize {
				fmt.Println("WARN: DoubleQueue lowChan half full.")
			}
		}
		select {
		case <- q.ctx.Done():
			fmt.Println("double priority queue context done: ", q.ctx.Err())
			return

		case task = <- q.highChan:
			task.Execute()
		default:
			select{
				case task = <- q.highChan:
					task.Execute()
				case task = <- q.lowChan:
					task.Execute()
			}
		}
	}
}
