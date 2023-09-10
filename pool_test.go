package pool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	p := New(2)
	var counter int32
	for i := 0; i < 3; i++ {
		p.Exec(func() {
			time.Sleep(2 * time.Second)
			atomic.AddInt32(&counter, 1)
		})
	}
	time.Sleep(3 * time.Second)
	assert.Equal(t, int32(2), atomic.LoadInt32(&counter))
	time.Sleep(3 * time.Second)
	assert.Equal(t, int32(3), atomic.LoadInt32(&counter))
}

func BenchmarkPool(b *testing.B) {
	b.ReportAllocs()
	p := New(10)
	wg := sync.WaitGroup{}
	wg.Add(b.N)
	for i := 0; i < b.N; i++ {
		p.Exec(func() {
			wg.Done()
		})
	}
	wg.Wait()
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
		q.PushBack(func() {})
		q.PopFront()
	}
}
