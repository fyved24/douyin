package utils

import "fmt"

func InitTime() {
	tm := NewTimerTask()
	_, err := tm.AddTaskByFunc("refreshCount", "@every 1m", RefreshRedisToDB)
	if err != nil {
		fmt.Println("定时任务启动失败")
	}
}
