package pipes

import (
	"context"
	"sync"
)

type Processor func(data []interface{}) interface{}

func Pipe(ctx context.Context, input chan interface{}, p Processor) chan interface{} {
	output := make(chan interface{})
	go func() {
		defer close(output)
		for val := range input {
			select {
			case <-ctx.Done():
				return
			case output <- p([]interface{}{val}):
			}
		}
	}()
	return output
}

func Source(ctx context.Context, data []interface{}) chan interface{} {
	retval := make(chan interface{})
	go func() {
		defer close(retval)
		for _, val := range data {
			select {
			case <-ctx.Done():
				return
			case retval <- val:
			}
		}
	}()
	return retval
}

func Sink(ctx context.Context, chans []chan interface{}) chan interface{} {
	output := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(len(chans))

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case output <- i:
			}
		}
	}

	for _, c := range chans {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

func FilterSink(ctx context.Context, chans []chan interface{}, p Processor) chan interface{} {
	output := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(len(chans))

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case output <- p([]interface{}{i}):
			}
		}
	}

	for _, c := range chans {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}
