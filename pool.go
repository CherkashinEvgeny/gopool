package pool

import (
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
)

type Task func()

type Pool struct {
	size  int64
	max   int64
	queue queue
}

var illegalPoolSizeError = errors.New("pool size must be > 0")

func New(max int) *Pool {
	if max <= 0 {
		panic(illegalPoolSizeError)
	}
	return &Pool{
		size:  0,
		max:   int64(max),
		queue: queue{},
	}
}

var illegalTaskError = errors.New("task must not be nil")

func (p *Pool) Exec(task Task) {
	if task == nil {
		panic(illegalTaskError)
	}
	size := p.scheduleTask(task)
	if size <= p.max {
		p.startNewWorker()
	}
}

func (p *Pool) scheduleTask(task Task) (size int64) {
	p.queue.PushBack(task)
	size = atomic.AddInt64(&p.size, 1)
	return
}

func (p *Pool) startNewWorker() {
	go p.workerPayload()
}

func (p *Pool) workerPayload() {
	for p.execTask() {
	}
}

func (p *Pool) execTask() (found bool) {
	task := p.queue.PopFront()
	if task == nil {
		return false
	}
	defer func() {
		atomic.AddInt64(&p.size, -1)
	}()
	task()
	return true
}

type queueElement struct {
	Data Task
	next *queueElement
}

type queueElementPool struct {
	values []*queueElement
}

func (p *queueElementPool) Get() (value *queueElement) {
	if len(p.values) == 0 {
		value = &queueElement{}
	} else {
		value = p.values[len(p.values)-1]
		p.values = p.values[:len(p.values)-1]
	}
	return
}

func (p *queueElementPool) Put(value *queueElement) {
	p.values = append(p.values, value)
}

type queue struct {
	front       *queueElement
	back        *queueElement
	elementPool queueElementPool
	lock        sync.Mutex
}

func (q *queue) PushBack(task Task) {
	q.lock.Lock()
	defer q.lock.Unlock()

	el := q.elementPool.Get()
	el.Data = task
	el.next = nil

	prev := q.back
	if prev != nil {
		prev.next = el
	}

	if q.front == nil {
		q.front = el
	}
	q.back = el
}

func (q *queue) PopFront() (task Task) {
	q.lock.Lock()
	defer q.lock.Unlock()

	el := q.front
	if el == nil {
		return
	}
	next := el.next

	q.front = next
	if next == nil {
		q.back = nil
	}

	task = el.Data
	el.Data = nil
	el.next = nil
	q.elementPool.Put(el)
	return
}
