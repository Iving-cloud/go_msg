package main

import (
	"fmt"
	"go_msg/util"
	"sync"
	"time"
)

var timers int = 0

//	type Job1 struct {
//		JobName string
//	}
//
//	func (t Job1) Run() {
//		fmt.Println(time.Now(), "I'm Job1", t.JobName)
//	}
func sayHello() {
	fmt.Println(time.Now().Second(), "hello Time")
}

func say() {
	for {
		timers += 1
		fmt.Println("Hello World")
		time.Sleep(1 * time.Second)
	}
}

func getInput() {
	for {
		var in int
		_, err := fmt.Scan(&in)
		if err != nil {
			break
		}
		util.Ch <- in
	}
}

func main() {
	var wg sync.WaitGroup

	//Job := Job1{JobName: "job1"}
	wg.Add(1)
	go util.MyTimer(sayHello)
	wg.Add(1)
	go say()
	wg.Add(1)
	go getInput()
	wg.Wait()
}
