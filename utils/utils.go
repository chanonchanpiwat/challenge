package utils

import (
	"sync"
)

func FanIn[T any](
	done <-chan interface{}, channels ...<-chan T,
) <-chan T {
	var wg sync.WaitGroup
	multiplexedStream := make(chan T)
	multiplex := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()
	return multiplexedStream
}

func Take[T any](done <-chan interface{}, ch <-chan T, count int) <-chan T {
	takeChannel := make(chan T)
	go func() {
		defer close(takeChannel)
		for i := count; i > 0 || i == -1; {
			if i != -1 {
				i--
			}
			select {
			case <-done:
				return
			case takeChannel <- <-ch:
			}
		}
	}()

	return takeChannel
}

