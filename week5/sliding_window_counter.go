package main

import (
	"fmt"
	"time"
)

type LimitError struct {
	code int
	msg  string
}

func (p LimitError) Error() string {
	return p.msg
}

var (
	ErrOverRequests = &LimitError{code: 400, msg: "requests number over threshold"}
)

type FixedWindowFlowLimit struct {
	timestamp         int64
	window            int64
	rateFlowThreshold int64
	counter           int64
}

func NewFixedWindowFlowLimit(windowSize, threshold int64) *FixedWindowFlowLimit {
	return &FixedWindowFlowLimit{
		timestamp:         time.Now().UnixNano() / 1e6,
		window:            windowSize,
		rateFlowThreshold: threshold,
		counter:           0,
	}
}

func (p *FixedWindowFlowLimit) Acquire(command string, callback func() error) *LimitError {
	now := time.Now().UnixNano() / 1e6
	window := now - p.timestamp
	if window < p.window {
		if p.counter > p.rateFlowThreshold {
			return ErrOverRequests
		}

		p.counter += 1
		err := callback()
		if err != nil {
			return &LimitError{
				code: 500,
				msg:  err.Error(),
			}
		}

	} else {
		p.counter = 0
		p.timestamp = p.timestamp + p.window
		return p.Acquire(command, callback)
	}
	return nil
}

type SlidingWindowCounter struct {
	timestamp         int64
	window            int64
	windowSize        int64
	box               int64
	rateFlowThreshold int64
	counter           []int64
}

func NewSlidingWindowCounter(window, windowSize, threshold int64) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		timestamp:         time.Now().UnixNano() / 1e6,
		window:            window,
		windowSize:        windowSize,
		box:               window / windowSize,
		rateFlowThreshold: threshold,
		counter:           make([]int64, windowSize),
	}
}

func (p SlidingWindowCounter) sum() int64 {
	var sum int64 = 0
	for i := 0; i < int(p.windowSize); i++ {
		sum += p.counter[i]
	}
	return sum
}

func (p SlidingWindowCounter) clear() {
	for i := 0; i < int(p.windowSize); i++ {
		p.counter[i] = 0
	}
}

func (p SlidingWindowCounter) Debug() {
	for i := 0; i < int(p.windowSize); i++ {
		fmt.Printf("%d ", p.counter[i])
	}
	fmt.Printf("%d\n", p.sum())
}

func (p SlidingWindowCounter) index(window int64) int64 {
	var leave int64
	if window > p.window {
		leave = window - p.window
	} else {
		leave = window
	}
	index := leave / p.box
	if leave%p.box == 0 {
		index -= 1
	}
	return index
}

func (p *SlidingWindowCounter) exchange(index int64) {
	counter := make([]int64, p.windowSize)
	from := 0
	for i := index; i < p.windowSize; i++ {
		counter[from] = p.counter[i]
		from += 1
	}
}

func (p *SlidingWindowCounter) Acquire(command string, callback func() error) *LimitError {
	window := time.Now().UnixNano()/1e6 - p.timestamp
	if window > 2*p.window {
		p.timestamp = time.Now().UnixNano() / 1e6
		p.clear()
		return p.Acquire(command, callback)
	} else if window > p.window {
		index := p.index(window)
		p.exchange(index)
		if p.sum() >= p.rateFlowThreshold {
			return ErrOverRequests
		}
		p.counter[index] += 1
		// p.debug()
		err := callback()
		if err != nil {
			return &LimitError{
				code: 500,
				msg:  err.Error(),
			}
		}
	} else if window > p.box {
		if p.sum() >= p.rateFlowThreshold {
			return ErrOverRequests
		}
		index := p.index(window)
		p.counter[index] += 1
		// p.debug()
		err := callback()
		if err != nil {
			return &LimitError{
				code: 500,
				msg:  err.Error(),
			}
		}
	}

	return nil
}

// func test1() {
// 	obj := NewFixedWindowFlowLimit(1000, 30)
// 	for i := 0; i < 35; i++ {
// 		err := obj.Acquire("test", func() error {
// 			println("Running...")
// 			return nil
// 		})
// 		if err != nil {
// 			fmt.Printf("Error %+v\n", err)
// 		}
// 	}
// 	time.Sleep(time.Millisecond * 1001)
// 	for i := 0; i < 29; i++ {
// 		err := obj.Acquire("test", func() error {
// 			println("Running...")
// 			return nil
// 		})
// 		if err != nil {
// 			fmt.Printf("Error %+v\n", err)
// 		}
// 	}
// 	time.Sleep(time.Millisecond * 1001)
// 	for i := 0; i < 40; i++ {
// 		err := obj.Acquire("test", func() error {
// 			println("Running...")
// 			return nil
// 		})
// 		if err != nil {
// 			fmt.Printf("Error %+v\n", err)
// 		}
// 	}
// }

// func test2() {
// 	return
// }

func main() {
	// test1()
	// test2()
	obj := NewSlidingWindowCounter(100, 10, 10)
	go func() {
		for {
			err := obj.Acquire("1 ", func() error {
				println()
				// println("1 running")
				return nil
			})
			if err != nil {
				obj.Debug()
				fmt.Printf("   %+v\n", err)
			}
			time.Sleep(time.Millisecond * 3)
		}
	}()
	// go func() {
	// 	for {
	// 		err := obj.Acquire("2 ", func() error {
	// 			// println("2 running")
	// 			println()
	// 			return nil
	// 		})
	// 		if err != nil {
	// 			obj.Debug()
	// 			fmt.Printf("   %+v\n", err)
	// 		}
	// 		time.Sleep(time.Millisecond * 7)
	// 	}
	// }()

	select {}
}
