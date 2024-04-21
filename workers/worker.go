package workers

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Worker struct {
	Conn     net.Conn
	WorkerId int
	Pid      int32
	Role     string
	ChTimer  chan int
}

type workerManger struct {
	workers      []*Worker
	workerLen    int
	MasterId     int
	MasterUpdate chan bool
}

var globaWokerID int = 1
var Manager = workerManger{workers: make([]*Worker, 0), MasterId: 0, workerLen: 0, MasterUpdate: make(chan bool)}
var WorkerWg sync.WaitGroup

func WorkerInit(conn net.Conn, pid int32) *Worker {
	var w Worker
	w.Pid = pid
	w.Conn = conn
	WorkerRoleSign(&w)
	WorkIdSign(&w)
	w.ChTimer = make(chan int, 10)
	return &w
}

func CheckMaster() bool {
	if Manager.MasterId == 0 {
		return false
	}
	return true

}

func PrintWorkerInfo() {
	for {
		fmt.Printf("---------------------------\n")
		fmt.Printf("workesLen:%d\r\n", len(Manager.workers))
		fmt.Printf("masterId:%d\r\n", Manager.MasterId)
		for i, w := range Manager.workers {
			fmt.Printf("Worker %d:\r\n", i+1)
			fmt.Printf("  WorkerId:%d\r\n", w.WorkerId)
			fmt.Printf("  Pid:%d\r\n", w.Pid)
			fmt.Printf("  Role:%s\r\n", w.Role)
			fmt.Printf("-----------------------\n")
			//fmt.Printf("ChTimer:%d\r\n", w.ChTimer)
		}
		fmt.Printf("------------------------\n")
		time.Sleep(5 * time.Second)
	}
}
func FindMaster() *Worker {
	if Manager.MasterId == 0 {
		return nil
	}
	for _, w := range Manager.workers {
		if w.WorkerId == Manager.MasterId {
			return w
		}
	}
	return nil
}
func WorkIdSign(worker *Worker) {
	worker.WorkerId = globaWokerID
	globaWokerID += 1
}
func AddWorker(worker Worker) bool {
	if Manager.workerLen >= 3 {
		return false
	}
	//WorkIdSign(&worker)
	WorkerRoleSign(&worker)
	Manager.workers = append(Manager.workers, &worker)
	Manager.workerLen++
	return true
}

func ResignMaster() bool {
	if Manager.workerLen == 0 {
		Manager.MasterId = 0
		return false
	}

	Manager.workers[0].Role = "Master"
	Manager.MasterId = Manager.workers[0].WorkerId
	Manager.MasterUpdate <- true
	return true
}

func DelWorker(worker Worker) bool {
	for i, w := range Manager.workers {
		if worker.WorkerId == w.WorkerId {
			Manager.workers = append(Manager.workers[:i], Manager.workers[i+1:]...)
			Manager.workerLen--
			break
		}
	}
	if worker.WorkerId == Manager.MasterId {
		ResignMaster()
	}
	return true
}
func WorkerRoleSign(worker *Worker) {
	if CheckMaster() {
		worker.Role = "Slave"
	} else {
		worker.Role = "Master"
		Manager.MasterId = worker.WorkerId
	}
}

func updateMaster() {
	for {
		select {
		case value := <-Manager.MasterUpdate:
			fmt.Println("UpdateMaster, value:", value)
			var worker = Worker{
				Conn:     nil,
				WorkerId: 0,
				Pid:      111,
				Role:     "Slave",
				ChTimer:  nil,
			}
			AddWorker(worker)
			fmt.Println(Manager)
			WorkerWg.Done()
		}
	}
}

//func main() {
//
//	WorkerWg.Add(1)
//	go updateMaster()
//	var worker1 = Worker{
//		Conn:     nil,
//		WorkerId: 0,
//		Pid:      111,
//		Role:     "Slave",
//		ChTimer:  nil,
//	}
//	var worker2 = Worker{
//		Conn:     nil,
//		WorkerId: 0,
//		Pid:      222,
//		Role:     "Slave",
//		ChTimer:  nil,
//	}
//	var worker3 *Worker
//	worker3 = WorkerInit(nil, 333)
//	AddWorker(worker1)
//	AddWorker(worker2)
//	AddWorker(*worker3)
//	fmt.Println(len(Manager.workers))
//	fmt.Println(worker1.Role)
//	fmt.Println(Manager)
//	fmt.Println(*FindMaster())
//	DelWorker(worker1)
//	fmt.Println(Manager)
//	time.Sleep(5 * time.Second)
//	WorkerWg.Wait()
//}
