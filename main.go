package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const ERRORS_COUNT = 5

func main() {
	var funcs []func() error

	for i := 0; i < 100; i++ {
		funcs = append(funcs, func() error {
			time.Sleep(1 * time.Second)
			fmt.Println("done")
			return errors.New("hi")
		})
	}

	var c = make(chan error, ERRORS_COUNT)
	wg := sync.WaitGroup{}

	for _, f := range funcs {
		wg.Add(1)
		go func(wgrp *sync.WaitGroup, allErr chan<- error) {
			err := f()
			if err != nil {
				allErr <- err
				if len(allErr) == cap(allErr) {
					panic("too many errors")
				}
			}
			wgrp.Done()
		}(&wg, c)
	}
	wg.Wait()
	fmt.Println("Funcs done")
}
