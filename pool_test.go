package pool

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	wg := sync.WaitGroup{}
	p := New(2)
	wg.Add(3)
	counter := 0
	for i := 0; i < 3; i++ {
		p.Exec(func() {
			time.Sleep(2 * time.Second)
			counter++
			wg.Done()
		})
	}
	time.Sleep(3 * time.Second)
	assert.Equal(t, counter, 2)
	wg.Wait()
	assert.Equal(t, counter, 3)
}

func BenchmarkPool(b *testing.B) {
	b.ReportAllocs()
	p := New(10)
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(2000)
		for j := 0; j < 2000; j++ {
			p.Exec(func() {
				wg.Done()
			})
		}
		wg.Wait()
	}
}

func TestQueue(t *testing.T) {
	q := queue{}
	k := 0
	for i := 0; i < 5; i++ {
		ii := i
		q.PushBack(func() {
			k = ii
		})
	}
	for i := 0; i < 5; i++ {
		q.PopFront()()
		assert.Equal(t, i, k)
	}
}

func BenchmarkQueue(b *testing.B) {
	b.ReportAllocs()
	q := queue{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 2000; j++ {
			q.PushBack(func() {})
		}
		for j := 0; j < 2000; j++ {
			q.PopFront()
		}
	}
}
