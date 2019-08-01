package main

import (
	"net"
	"runtime"
	"sync"
)

type Task func(*net.UDPConn) error // this function should not contain dead loop
type FuncQueue chan Task

type WorkerPool struct {
	TaskQueue    FuncQueue
	NumWorkers   int
	BufferSize   int
	Mu           sync.Mutex
	SuccessCount uint
	CloseFlag    bool
	wg           sync.WaitGroup
}

func NewWorkerPool(routineNumber int, bufferSize int) WorkerPool {
	var wp WorkerPool
	wp.Init(routineNumber, bufferSize)

	return wp
}
func (p *WorkerPool) Init(routineNumber int, bufferSize int) {
	p.TaskQueue = make(FuncQueue, bufferSize)

	if availableCPU := runtime.NumCPU(); routineNumber > availableCPU {
		p.NumWorkers = availableCPU
	} else {
		p.NumWorkers = routineNumber
	}

	p.BufferSize = bufferSize
	p.Mu = sync.Mutex{}
	p.SuccessCount = uint(0)
	p.CloseFlag = false
}

func (p *WorkerPool) worker(taskQueue FuncQueue) {

	defer p.wg.Done()
	
	var conn *net.UDPConn
	var err error
	for retry := 0; retry < 5; retry++ {
		conn, err = connectCollector(CollectorAddr)
		if err != nil {
			log.Warnf("Failed to connect remote collector, current retry number %d\n", retry+1)
			continue
		}
		break
	}

	if conn == nil {
		log.Error("Failed to connect remote collector\n")
	}
	log.Info("Successfully launched worker, start to process task\n")

	var task Task
	for {
		if p.CloseFlag {
			return
		}
		
		task = <-taskQueue
		
		if task(conn) == nil {
			p.Mu.Lock()
			p.SuccessCount++
			p.Mu.Unlock()
		}
	}
}

func (p *WorkerPool) Start() {
	for i := 0; i < p.NumWorkers; i++ {
		go p.worker(p.TaskQueue)
	}
}

func (p *WorkerPool) Stop() {
	p.CloseFlag = true
	
	p.wg.Wait()
	close(p.TaskQueue)
}

func (p *WorkerPool) AddTask(task Task) {
	p.TaskQueue <- task
}

func (p *WorkerPool) AddMultipleTask(task Task, num int) {
	for i := 0; i < num; i++ {
		p.TaskQueue <- task
	}
}

func (p *WorkerPool) GetSuccessCount() uint {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	return p.SuccessCount
}
