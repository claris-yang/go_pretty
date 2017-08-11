package dispatcher

import (
	"fmt"
	"runtime"
	"time"
)

/**
 * 小米
 *
 * @author yangtao
 * @version V1.0 创建时间：2016-12-21
 *
 */

type IPayload interface {
	Handle()
	ToString() int
}

const (
	MaxWorder = 5
	MaxQueue  = 5
)

type Job struct {
	IPayload
	Req  int
	resp map[int]int
}

func (this *Job) ToString() int {
	return this.Req
}

func (this *Job) Handle() {
	//resp := make(map[int]int)
	var result int
	for i := 0; i < 1000000; i++ {
		this.resp[i] = i
		result = i
	}
	fmt.Printf("No %d, first data is %d, last data is %d\n", this.Req, this.resp[0], this.resp[result])
}

func NewJob(r int) *Job {
	return &Job{
		Req:  r,
		resp: make(map[int]int),
	}
}

type Worker struct {
	WorkerPool chan chan IPayload
	JobChannel chan IPayload
	quit       chan bool
}

func newWorker(workerPool chan chan IPayload) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan IPayload),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		w.WorkerPool <- w.JobChannel

		for {
			select {
			case job := <-w.JobChannel:
				job.Handle()
				w.WorkerPool <- w.JobChannel
			case <-time.After(time.Second):
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan IPayload, maxWorkers)
	job := make(chan IPayload, maxWorkers)

	return &Dispatcher{workerPool: pool,
		maxWorkers: maxWorkers,
		JobQueue:   job,
		quit:       make(chan bool),
		closeFlag:  false,
		Count:      0,
	}
}

func (d *Dispatcher) GOMAXPROCES() {
	if d.maxWorkers > runtime.NumCPU() {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(d.maxWorkers)
	}
}
func (d *Dispatcher) Run() {

	d.GOMAXPROCES()

	for i := 0; i < d.maxWorkers; i++ {
		go func() {
			worker := newWorker(d.workerPool)
			worker.Start()
			d.workerGroup = append(d.workerGroup, &worker)
		}()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func(job IPayload) {
				jobChannel := <-d.workerPool
				jobChannel <- job
				//fmt.Println(" job over _ ", job.ToString())
				fmt.Println(" @ ")
				d.Count++
			}(job)
		case quit := <-d.quit:
			if quit {
				break
			}
		}
	}
}

type Dispatcher struct {
	workerPool  chan chan IPayload
	maxWorkers  int
	JobQueue    chan IPayload
	quit        chan bool
	workerGroup []*Worker
	closeFlag   bool
	Count       int
}

func (d *Dispatcher) Stop() {
	for _, worker := range d.workerGroup {
		worker.Stop()
	}
	d.workerGroup = d.workerGroup[0:0]
	d.closeFlag = true
	d.quit <- true
}

func (d *Dispatcher) AddJob(p IPayload) {
	if !d.closeFlag {
		//work := Job{IPayload: p}
		go func(p IPayload) {
			d.JobQueue <- p
		}(p)
	}
}
