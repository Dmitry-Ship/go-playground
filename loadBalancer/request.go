package loadBalancer

import (
	"container/heap"
	"fmt"
	"math/rand"
	"time"
)

const nRequester = 100
const nWorker = 10

// Simulation of some tasksChan: just sleep for a while and report how long.
func op() int {
	n := rand.Int63n(int64(time.Second))
	time.Sleep(time.Duration(nWorker * n))
	return int(n)
}

type Request struct {
	process    func() int
	resultChan chan int
}

func makeRequest(tasksChan chan Request) {
	c := make(chan int)
	for {
		time.Sleep(time.Duration(rand.Int63n(int64(nWorker * 2 * time.Second))))
		tasksChan <- Request{op, c}
		<-c
	}
}

type Worker struct {
	i        int
	requests chan Request
	pending  int
}

func (w *Worker) run(done chan *Worker) {
	for {
		req := <-w.requests
		req.resultChan <- req.process()
		done <- w
	}
}

type Pool []*Worker

func (p Pool) Len() int { return len(p) }

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].i = i
	a[j].i = j
}

func (p *Pool) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	w := x.(*Worker)
	a[n] = w
	w.i = n
	*p = a
}

func (p *Pool) Pop() interface{} {
	a := *p
	*p = a[0 : len(a)-1]
	w := a[len(a)-1]
	w.i = -1 // for safety
	return w
}

type Balancer struct {
	pool     Pool
	doneChan chan *Worker
	i        int
}

func NewBalancer() *Balancer {
	done := make(chan *Worker, nWorker)
	b := &Balancer{make(Pool, 0, nWorker), done, 0}
	for i := 0; i < nWorker; i++ {
		w := &Worker{requests: make(chan Request, nRequester)}
		heap.Push(&b.pool, w)
		go w.run(b.doneChan)
	}
	return b
}

func (b *Balancer) balance(tasksChan chan Request) {
	for {
		select {
		case req := <-tasksChan:
			b.dispatch(req)
		case w := <-b.doneChan:
			b.completed(w)
		}
		b.print()
	}
}

func (b *Balancer) print() {
	sum := 0
	sumsq := 0
	for _, w := range b.pool {
		fmt.Printf("%d ", w.pending)
		sum += w.pending
		sumsq += w.pending * w.pending
	}
	avg := float64(sum) / float64(len(b.pool))
	variance := float64(sumsq)/float64(len(b.pool)) - avg*avg
	fmt.Printf(" %.2f %.2f\n", avg, variance)
}

func (b *Balancer) dispatch(req Request) {
	w := heap.Pop(&b.pool).(*Worker)
	w.requests <- req
	w.pending++
	//	fmt.Printf("started %p; now %d\n", w, w.pending)
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
	w.pending--
	//	fmt.Printf("finished %p; now %d\n", w, w.pending)
	heap.Remove(&b.pool, w.i)
	heap.Push(&b.pool, w)
}

func TestBalancer() {
	tasksChan := make(chan Request)
	for i := 0; i < nRequester; i++ {
		go makeRequest(tasksChan)
	}
	NewBalancer().balance(tasksChan)
}
