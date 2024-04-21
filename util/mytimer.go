package util

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"go_msg/workers"
	"time"
)

var Ch chan int

type MyJob struct {
	W *workers.Worker
	c *cron.Cron
}

func (m MyJob) Run() {
	workers.DelWorker(*m.W)
	m.c.Stop()
	fmt.Println(time.Now().Second(), "s, MyJob timeout, Stop it, pid:", m.W.Pid)
}
func MyTimer(job MyJob) {
	job.c = cron.New(cron.WithSeconds())
	// 添加定时任务
	ID, err := job.c.AddJob("@every 10s", job)
	if err != nil {
		fmt.Println("添加定时任务失败：", err)
		return
	} else {
		fmt.Println("entryID:", ID)
	}
	ending := 0
	for ending != 1 {
		select {
		case value := <-job.W.ChTimer:
			if value == 1 {
				fmt.Println("value:", value)
				job.c.Remove(ID)
				ID, err = job.c.AddJob("@every 10s", job)
				if err != nil {
					fmt.Println("addFunc err:", err)
					return
				} else {
					fmt.Println(time.Now().Second(), "pid:", job.W.Pid, "resign entry:", ID)
				}
				job.c.Start()
			} else if value == 0 {
				fmt.Println("timer start, value:", value)
				job.c.Start()
			} else {
				fmt.Println("timer stop, value", value)
				job.c.Stop()
				ending = 1
			}
		}
	}
	fmt.Println("timer ending")
}
