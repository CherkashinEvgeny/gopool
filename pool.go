package pool

import (
	"sync"
)

type Task func()

type Pool struct {
	size        int
	concurrency int
	queue       queue
	lock        sync.Mutex
}

func New(concurrency int) *Pool {
	if concurrency <= 0 {
		panic("pool size must be > 0")
	}
	return &Pool{
		size:        0,
		concurrency: concurrency,
		queue:       queue{},
		lock:        sync.Mutex{},
	}
}

func (p *Pool) Exec(task Task) {
	if task == nil {
		panic("task must not be nil")
	}
	p.lock.Lock()
	p.queue.Enqueue(task)
	runNewWorker := p.size < p.concurrency
	if runNewWorker {
		p.size++
	}
	p.lock.Unlock()
	if runNewWorker {
		go p.worker()
	}
}

func (p *Pool) worker() {
	for {
		p.lock.Lock()
		task, found := p.queue.Dequeue()
		if !found {
			p.size--
		}
		p.lock.Unlock()
		if !found {
			break
		}
		task()
	}
}

type queue struct {
	front       *queueElement
	back        *queueElement
	elementPool queueElementPool
}

func (q *queue) Enqueue(task Task) {
	el := q.elementPool.Get()
	el.value = task
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

func (q *queue) Dequeue() (task Task, found bool) {
	el := q.front
	if el == nil {
		return nil, false
	}
	next := el.next

	q.front = next
	if next == nil {
		q.back = nil
	}

	task = el.value
	el.value = nil
	el.next = nil
	q.elementPool.Put(el)
	return task, true
}

type queueElement struct {
	value Task
	next  *queueElement
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
	return value
}

func (p *queueElementPool) Put(value *queueElement) {
	p.values = append(p.values, value)
}
