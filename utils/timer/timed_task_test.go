package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var job = mockJob{}

type mockJob struct{}

func (job mockJob) Run() {
	mockFunc()
}

const LAYOUT = "2006-01-02 15:04:05"

func mockFunc() {
	//time.Sleep(time.Second)
	fmt.Println(time.Now().Format(LAYOUT))
}

func TestNewTimerTask(t *testing.T) {
	tm := NewTimerTask()
	_tm := tm.(*timer)

	_, err := tm.AddTaskByFunc("func", "@every 1m", mockFunc)
	assert.Nil(t, err)
	_, ok := _tm.taskList["func"]
	if !ok {
		t.Error("no find func")
	}

	//{
	//	_, err := tm.AddTaskByJob("job", "@every 1s", job)
	//	assert.Nil(t, err)
	//	_, ok := _tm.taskList["job"]
	//	if !ok {
	//		t.Error("no find job")
	//	}
	//}
	//
	//{
	//	_, ok := tm.FindCron("func")
	//	if !ok {
	//		t.Error("no find func")
	//	}
	//	_, ok = tm.FindCron("job")
	//	if !ok {
	//		t.Error("no find job")
	//	}
	//	_, ok = tm.FindCron("none")
	//	if ok {
	//		t.Error("find none")
	//	}
	//}
	//{
	//	tm.Clear("func")
	//	_, ok := tm.FindCron("func")
	//	if ok {
	//		t.Error("find func")
	//	}
	//}
	for {

	}
}
