package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDemoBasic(t *testing.T) {
	ch := make(chan int)
	go func() {
		ch <- 42
	}()
	v := <-ch
	fmt.Println("basic:", v)
}

func TestDemoBuffered(t *testing.T) {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	close(ch)
	for v := range ch {
		fmt.Println("buffered:", v)
	}
}

func TestDemoSelectTimeout(t *testing.T) {
	ch := make(chan int)
	select {
	case v := <-ch:
		fmt.Println("select:", v)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("select:", "timeout")
	}
}

func TestDemoWaitGroup(t *testing.T) {
	jobs := []int{1, 2, 3, 4}
	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, j := range jobs {
		go func(x int) {
			defer wg.Done()
			fmt.Println("wg:", x)
		}(j)
	}
	wg.Wait()
}

func TestDemoMutex(t *testing.T) {
	var mu sync.Mutex
	var n int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			mu.Lock()
			n++
			mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			mu.Lock()
			n++
			mu.Unlock()
		}
	}()
	wg.Wait()
	fmt.Println("mutex:", n)
}

func TestDemoAtomic(t *testing.T) {
	var n int64
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			atomic.AddInt64(&n, 1)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			atomic.AddInt64(&n, 1)
		}
	}()
	wg.Wait()
	fmt.Println("atomic:", n)
}

func TestDemoContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	ch := make(chan string)
	go func() {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ctx.Done():
			return
		case ch <- "ok":
		}
	}()
	select {
	case s := <-ch:
		fmt.Println("context:", s)
	case <-ctx.Done():
		fmt.Println("context:", "cancel")
	}
}

func TestDemoWorkerPool(t *testing.T) {
	jobCh := make(chan int)
	resCh := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				resCh <- j * j
			}
		}()
	}
	go func() {
		for i := 1; i <= 5; i++ {
			jobCh <- i
		}
		close(jobCh)
	}()
	go func() {
		wg.Wait()
		close(resCh)
	}()
	for r := range resCh {
		fmt.Println("pool:", r)
	}
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func TestDemoPipeline(t *testing.T) {
	for v := range square(gen(1, 2, 3, 4, 5)) {
		fmt.Println("pipe:", v)
	}
}

func TestDemoBounded(t *testing.T) {
	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(n int) {
			sem <- struct{}{}
			defer func() {
				<-sem
				wg.Done()
			}()
			time.Sleep(10 * time.Millisecond)
			fmt.Println("bounded:", n)
		}(i)
	}
	wg.Wait()
}

func TestProducerConsumer(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	jobs := make(chan string, 2)
	var wg sync.WaitGroup

	consumers := 2
	for id := 1; id <= consumers; id++ {
		wg.Add(1)
		go func(cid int) {
			defer wg.Done()
			for it := range jobs {
				fmt.Println("consume:", cid, it)
			}
		}(id)
	}

	go func() {
		for _, it := range items {
			fmt.Println("produce:", it)
			jobs <- it
		}
		close(jobs)
	}()

	wg.Wait()
}
