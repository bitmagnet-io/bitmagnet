package semaphore

import (
	"context"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func checkLimit(t *testing.T, sem Semaphore, expected int) {
	limit := sem.GetLimit()
	if limit != expected {
		t.Error("semaphore must have limit = ", expected, ", but has ", limit)
	}
}

func checkCount(t *testing.T, sem Semaphore, expected int) {
	count := sem.GetCount()
	if count != expected {
		t.Error("semaphore must have count = ", expected, ", but has ", count)
	}
}

func checkLimitAndCount(t *testing.T, sem Semaphore, expectedLimit, expectedCount int) {
	checkLimit(t, sem, expectedLimit)
	checkCount(t, sem, expectedCount)
}

func TestNew(t *testing.T) {
	tests := []func(){
		func() {
			sem := New(1)
			checkLimitAndCount(t, sem, 1, 0)
		},
		func() {
			sem := New(0)
			checkLimitAndCount(t, sem, 0, 0)
		},
		func() {
			defer func() {
				if recover() == nil {
					t.Error("Panic expected")
				}
			}()
			_ = New(-1)
		},
	}
	for _, test := range tests {
		test()
	}
}

func TestSemaphore_Acquire(t *testing.T) {
	sem := New(1)

	err := sem.Acquire(nil, 1)

	if err != nil {
		t.Error("Error returned:", err.Error())
	}
	checkLimitAndCount(t, sem, 1, 1)
}

func TestSemaphore_Acquire_zero_panic_expected(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Panic expected")
		}
	}()
	sem := New(1)

	// acquire zero should panic
	sem.Acquire(nil, 0)
}

func TestSemaphore_Acquire_with_ctx(t *testing.T) {
	sem := New(1)

	err := sem.Acquire(context.Background(), 1)

	if err != nil {
		t.Error("Error returned:", err.Error())
	}
	checkLimitAndCount(t, sem, 1, 1)
}

func TestSemaphore_Acquire_ctx_done(t *testing.T) {
	sem := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	cancel() // make ctx.Done()

	err := sem.Acquire(ctx, 1)

	if err != context.Canceled {
		t.Error("Error is not context.Canceled")
	}
	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_TryAcquire(t *testing.T) {
	sem := New(1)

	if !(sem.TryAcquire(2) == false) {
		t.Fail()
	}
	checkLimitAndCount(t, sem, 1, 0)

	if !(sem.TryAcquire(1) == true) {
		t.Fail()
	}
	checkLimitAndCount(t, sem, 1, 1)

	if !(sem.TryAcquire(1) == false) {
		t.Fail()
	}
	checkLimitAndCount(t, sem, 1, 1)
}

func TestSemaphore_TryAcquire_panic_expected(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Panic expected")
		}
	}()
	sem := New(1)
	sem.TryAcquire(0)
}

func TestSemaphore_TryAcquire_contention(t *testing.T) {
	sem := New(5)

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			<-c
			for j := 0; j < 10000; j++ {
				if sem.TryAcquire(1) {
					runtime.Gosched()
					sem.Release(1)
				}
			}
			wg.Done()
		}()
	}

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 5, 0)
}

func TestSemaphore_Release(t *testing.T) {
	sem := New(1)

	sem.Acquire(nil, 1)
	oldCnt := sem.Release(1)

	if oldCnt != 1 {
		t.Error("semaphore must have old count = ", 1, ", but has ", oldCnt)
	}
	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_Release_zero_panic_expected(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Panic expected")
		}
	}()
	sem := New(1)
	sem.Acquire(nil, 1)
	sem.Release(0)
}

func TestSemaphore_Release_with_ctx(t *testing.T) {
	sem := New(1)

	sem.Acquire(context.Background(), 1)
	oldCnt := sem.Release(1)

	if oldCnt != 1 {
		t.Error("semaphore must have old count = ", 1, ", but has ", oldCnt)
	}
	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_Release_without_Acquire_panic_expected(t *testing.T) {
	sem := New(1)

	defer func() {
		if recover() == nil {
			t.Error("Panic expected")
		}
	}()
	sem.Release(1)
}

func TestSemaphore_Release_more_than_Acquired_panic_expected(t *testing.T) {
	sem := New(1)

	defer func() {
		if recover() == nil {
			t.Error("Panic expected")
		}
	}()
	sem.Acquire(context.Background(), 1)
	sem.Release(2)
}

func TestSemaphore_Acquire_2_times_Release_2_times(t *testing.T) {
	sem := New(2)
	checkLimitAndCount(t, sem, 2, 0)

	sem.Acquire(nil, 1)
	checkLimitAndCount(t, sem, 2, 1)

	sem.Acquire(nil, 1)
	checkLimitAndCount(t, sem, 2, 2)

	oldCnt := sem.Release(1)
	if oldCnt != 2 {
		t.Error("semaphore must have old count = ", 2, ", but has ", oldCnt)
	}
	checkLimitAndCount(t, sem, 2, 1)

	oldCnt = sem.Release(1)
	if oldCnt != 1 {
		t.Error("semaphore must have old count = ", 1, ", but has ", oldCnt)
	}
	checkLimitAndCount(t, sem, 2, 0)
}

func TestSemaphore_Acquire_Release_under_limit(t *testing.T) {
	sem := New(100)

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			<-c
			err := sem.Acquire(nil, 1)
			if err != nil {
				panic(err)
			}
			sem.Release(1)
			wg.Done()
		}()
	}

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 100, 0)
}

func TestSemaphore_Acquire_Release_under_limit_ctx_done(t *testing.T) {
	sem := New(10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			<-c
			for {
				err := sem.Acquire(ctx, 1)
				if err != nil {
					if err == context.DeadlineExceeded {
						break
					}
					panic(err)
				}
				sem.Release(1)
			}
			wg.Done()
		}()
	}

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 10, 0)
}

func TestSemaphore_Acquire_Release_over_limit(t *testing.T) {
	sem := New(1)

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			<-c
			for j := 0; j < 1000; j++ {
				err := sem.Acquire(nil, 1)
				if err != nil {
					panic(err)
				}
				sem.Release(1)
			}
			wg.Done()
		}()
	}

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_Acquire_Release_over_limit_ctx_done(t *testing.T) {
	sem := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			<-c
			for {
				err := sem.Acquire(ctx, 1)
				if err != nil {
					if err == context.DeadlineExceeded {
						break
					}
					panic(err)
				}
				sem.Release(1)
			}
			wg.Done()
		}()
	}

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_SetLimit(t *testing.T) {
	sem := New(1)
	checkLimitAndCount(t, sem, 1, 0)

	sem.SetLimit(2)
	checkLimitAndCount(t, sem, 2, 0)

	sem.SetLimit(1)
	checkLimitAndCount(t, sem, 1, 0)

	sem.SetLimit(0)
	checkLimitAndCount(t, sem, 0, 0)
}

func TestSemaphore_SetLimit_negative_limit_panic_expected(t *testing.T) {
	sem := New(1)
	checkLimitAndCount(t, sem, 1, 0)

	defer func() {
		if recover() == nil {
			t.Error("Panic expected")
		}
	}()
	sem.SetLimit(-1)
}

func TestSemaphore_SetLimit_increase_limit(t *testing.T) {
	sem := New(1)
	checkLimitAndCount(t, sem, 1, 0)

	sem.Acquire(nil, 1)
	checkLimitAndCount(t, sem, 1, 1)

	sem.SetLimit(2)
	checkLimitAndCount(t, sem, 2, 1)

	sem.Acquire(nil, 1)
	checkLimitAndCount(t, sem, 2, 2)

	sem.Release(1)
	checkLimitAndCount(t, sem, 2, 1)

	sem.Release(1)
	checkLimitAndCount(t, sem, 2, 0)
}

func TestSemaphore_SetLimit_decrease_limit(t *testing.T) {
	sem := New(2)
	checkLimitAndCount(t, sem, 2, 0)

	sem.Acquire(nil, 1)
	checkLimitAndCount(t, sem, 2, 1)

	sem.Acquire(nil, 1)
	checkLimitAndCount(t, sem, 2, 2)

	sem.SetLimit(1)
	checkLimitAndCount(t, sem, 1, 2)

	sem.Release(1)
	checkLimitAndCount(t, sem, 1, 1)

	sem.Release(1)
	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_New0_SetLimit1_Acquire_Release_SetLimit0(t *testing.T) {
	sem := New(0)
	checkLimitAndCount(t, sem, 0, 0)

	sem.SetLimit(1)
	checkLimitAndCount(t, sem, 1, 0)

	sem.Acquire(nil, 1)
	checkLimitAndCount(t, sem, 1, 1)

	oldCnt := sem.Release(1)
	if oldCnt != 1 {
		t.Error("semaphore must have old count = ", 1, ", but has ", oldCnt)
	}
	checkLimitAndCount(t, sem, 1, 0)

	sem.SetLimit(0)
	checkLimitAndCount(t, sem, 0, 0)
}

func TestSemaphore_SetLimit_increase_broadcast(t *testing.T) {
	getWGs := func(cnt int) []*sync.WaitGroup {
		wgs := make([]*sync.WaitGroup, cnt)
		for i := range wgs {
			wgs[i] = &sync.WaitGroup{}
			wgs[i].Add(1)
		}
		return wgs
	}

	sem := New(1)
	sem.Acquire(nil, 1)

	innerWGs := getWGs(2)
	outerWGs := getWGs(2)

	go func() {
		outerWGs[0].Wait()

		// here we a trying to acquire over limit
		checkLimitAndCount(t, sem, 1, 1)
		sem.Acquire(nil, 1)

		innerWGs[0].Done()
		outerWGs[1].Wait()

		sem.Release(1)

		innerWGs[1].Done()
	}()

	checkLimitAndCount(t, sem, 1, 1)
	outerWGs[0].Done()

	time.Sleep(100 * time.Millisecond)

	// increase limit so inner goroutine can acquire semaphore
	sem.SetLimit(2)
	checkLimitAndCount(t, sem, 2, 1)

	innerWGs[0].Wait()

	checkLimitAndCount(t, sem, 2, 2)

	outerWGs[1].Done()
	innerWGs[1].Wait()

	checkLimitAndCount(t, sem, 2, 1)

	sem.Release(1)
	checkLimitAndCount(t, sem, 2, 0)
}

func TestSemaphore_Acquire_Release_SetLimit_under_limit(t *testing.T) {
	sem := New(100)

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			<-c
			for j := 0; j < 10000; j++ {
				err := sem.Acquire(nil, 1)
				if err != nil {
					panic(err)
				}
				runtime.Gosched()
				sem.Release(1)
				runtime.Gosched()
			}
			wg.Done()
		}()
	}

	c2 := make(chan struct{})
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		<-c
		for {
			select {
			case <-c2:
				sem.SetLimit(100)
				wg2.Done()
				return
			default:
			}
			newLimit := rand.Intn(50) + 50 // range [50, 99]
			sem.SetLimit(newLimit)
			runtime.Gosched()
		}

	}()

	close(c) // start
	wg.Wait()

	close(c2) // stop 'set limit' goroutine
	wg2.Wait()

	checkLimitAndCount(t, sem, 100, 0)
}

func TestSemaphore_Acquire_Release_SetLimit_under_limit_ctx_done(t *testing.T) {
	sem := New(100)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			<-c
			for {
				err := sem.Acquire(ctx, 1)
				if err != nil {
					if err == context.DeadlineExceeded {
						break
					}
					panic(err)
				}
				runtime.Gosched()
				sem.Release(1)
				runtime.Gosched()
			}
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		<-c
		for {
			select {
			case <-ctx.Done():
				sem.SetLimit(100)
				wg.Done()
				return
			default:
			}
			newLimit := rand.Intn(50) + 50 // range [50, 99]
			sem.SetLimit(newLimit)
			runtime.Gosched()
		}

	}()

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 100, 0)
}

func TestSemaphore_Acquire_Release_SetLimit_over_limit(t *testing.T) {
	sem := New(1)

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			<-c
			for j := 0; j < 10000; j++ {
				err := sem.Acquire(nil, 1)
				if err != nil {
					panic(err)
				}
				runtime.Gosched()
				sem.Release(1)
				runtime.Gosched()
			}
			wg.Done()
		}()
	}

	c2 := make(chan struct{})
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		<-c
		for {
			select {
			case <-c2:
				sem.SetLimit(1)
				wg2.Done()
				return
			default:
			}
			newLimit := rand.Intn(50) + 1 // range [1, 50]
			sem.SetLimit(newLimit)
			runtime.Gosched()
		}

	}()

	close(c) // start
	wg.Wait()

	close(c2) // stop 'set limit' goroutine
	wg2.Wait()

	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_Acquire_Release_SetLimit_over_limit_ctx_done(t *testing.T) {
	sem := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			<-c
			for {
				err := sem.Acquire(ctx, 1)
				if err != nil {
					if err == context.DeadlineExceeded {
						break
					}
					panic(err)
				}
				runtime.Gosched()
				sem.Release(1)
				runtime.Gosched()
			}
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		<-c
		for {
			select {
			case <-ctx.Done():
				sem.SetLimit(1)
				wg.Done()
				return
			default:
			}
			newLimit := rand.Intn(50) + 1 // range [1, 50]
			sem.SetLimit(newLimit)
			runtime.Gosched()
		}

	}()

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_Acquire_Release_SetLimit_random_limit(t *testing.T) {
	sem := New(1)

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			<-c
			for j := 0; j < 10000; j++ {
				err := sem.Acquire(nil, 1)
				if err != nil {
					panic(err)
				}
				runtime.Gosched()
				sem.Release(1)
				runtime.Gosched()
			}
			wg.Done()
		}()
	}

	c2 := make(chan struct{})
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		<-c
		for {
			select {
			case <-c2:
				sem.SetLimit(1)
				wg2.Done()
				return
			default:
			}
			newLimit := rand.Intn(200) + 1 // range [1, 200]
			sem.SetLimit(newLimit)
			runtime.Gosched()
		}

	}()

	close(c) // start
	wg.Wait()

	close(c2) // stop 'set limit' goroutine
	wg2.Wait()

	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_Acquire_Release_SetLimit_random_limit_ctx_done(t *testing.T) {
	sem := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan struct{})
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			<-c
			for {
				err := sem.Acquire(ctx, 1)
				if err != nil {
					if err == context.DeadlineExceeded {
						break
					}
					panic(err)
				}
				runtime.Gosched()
				sem.Release(1)
				runtime.Gosched()
			}
			wg.Done()
		}()
	}

	wg.Add(1)
	go func() {
		<-c
		for {
			select {
			case <-ctx.Done():
				sem.SetLimit(1)
				wg.Done()
				return
			default:
			}
			newLimit := rand.Intn(200) + 1 // range [1, 200]
			sem.SetLimit(newLimit)
			runtime.Gosched()
		}

	}()

	close(c) // start
	wg.Wait()

	checkLimitAndCount(t, sem, 1, 0)
}

func TestSemaphore_broadcast_channel_race(t *testing.T) {
	threads := 4
	acquiresPerRun := 5

	// runTest method creates a short-lived semaphore with contention over only
	// a few attempts per thread. The condition being tested is a regression
	// in which the last thread to call acquire in the group will hang forever.
	runTest := func(done chan struct{}) {
		sem := New(1)
		wg := sync.WaitGroup{}
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < acquiresPerRun; j++ {
					runtime.Gosched()
					if err := sem.Acquire(context.Background(), 1); err != nil {
						t.Fatal(err)
					}
					sem.Release(1)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(done)
	}

	// Run several iterations of the runTest method, which hopefully will result
	// in one the iterations hanging.
	for run := 0; run < 1000; run++ {
		done := make(chan struct{})
		go runTest(done)
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("single run took more than ten seconds to finish")
		}
	}
}

func TestSemaphore_weighted_acquire_gt_release(t *testing.T) {
	runTest := func(done chan struct{}) {
		weight := 1
		biggerWeight := weight * 10
		limit := weight * 100
		releaseSignalChan := make(chan struct{})
		semp := New(limit)

		for i := 0; i < limit; i++ {
			_ = semp.Acquire(context.TODO(), weight)
			go func() {
				<-releaseSignalChan
				semp.Release(weight)
			}()
		}

		// broadcast to release
		go close(releaseSignalChan)

		// Try to acquire while releasing
		_ = semp.Acquire(context.TODO(), biggerWeight)
		close(done)
	}

	// Run several iterations of the runTest method, which hopefully will result
	// in one the iterations hanging.
	for run := 0; run < 100000; run++ {
		done := make(chan struct{})
		go runTest(done)
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("single run took more than ten seconds to finish")
		}
	}
}
