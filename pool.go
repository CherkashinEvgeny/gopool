package pool

import (
	"errors"
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
		panic(illegalPoolSizeError)
	}
	return &Pool{
		size:        0,
		concurrency: concurrency,
		queue:       queue{},
		lock:        sync.Mutex{},
	}
}

var illegalPoolSizeError = errors.New("pool size must be > 0")

func (p *Pool) Exec(task Task) {
	if task == nil {
		panic(nilTaskError)
	}
	size := p.schedule(task)
	if size <= p.concurrency {
		go p.worker()
	}
}

var nilTaskError = errors.New("task must not be nil")

func (p *Pool) schedule(task Task) (size int) {
	p.lock.Lock()
	p.queue.PushBack(task)
	p.size++
	size = p.size
	p.lock.Unlock()
	return size
}

func (p *Pool) worker() {
	for p.exec() {
	}
}

func (p *Pool) exec() (found bool) {
	p.lock.Lock()
	task := p.queue.PopFront()
	p.lock.Unlock()
	if task == nil {
		return false
	}
	defer func() {
		p.lock.Lock()
		p.size--
		p.lock.Unlock()
	}()
	task()
	return true
}

type queue struct {
	front       *queueElement
	back        *queueElement
	elementPool queueElementPool
}

func (q *queue) PushBack(task Task) {
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

func (q *queue) PopFront() (task Task) {
	el := q.front
	if el == nil {
		return nil
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
	return task
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
