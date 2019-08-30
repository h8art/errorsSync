package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

const errorsCount = 5

func main() {
	var funcs []func() error

	for i := 0; i < 100; i++ {
		funcs = append(funcs, func() error {
			time.Sleep(1 * time.Second)
			fmt.Println("done")
			return errors.New("hi")
		})
	}

	var c = make(chan error, len(funcs))
	wg := sync.WaitGroup{}

	for _, f := range funcs {
		wg.Add(1)
		go func(wgrp *sync.WaitGroup, allErr chan<- error) {
			err := f()
			if err != nil {
				allErr <- err
				if len(allErr) == errorsCount {
					log.Fatal("too many errors")
				}
			}
			wgrp.Done()
		}(&wg, c)
	}
	wg.Wait()
	fmt.Println("Funcs done")
}
