package util

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type emptyTask struct {
	value int
}

func (e *emptyTask) Execute() {
	fmt.Println(e.value)
}

func TestDoubleQueue(t *testing.T) {
	queue := NewDoubleQueue(context.Background(), 3)
	rander := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tmpValue int
	go queue.Start(true)
	for i:=0; i<100; i++ {
		tmpValue = rander.Intn(2)
		if tmpValue == 0 {
			queue.AddLow(&emptyTask{value: tmpValue})
		} else {
			queue.AddHigh(&emptyTask{value: tmpValue})
		}
	}
	time.Sleep(50*time.Millisecond)
}
