package util

import (
	"fmt"
	"sync"
	"time"
)

var w  = &sync.WaitGroup{}

func ConcurrentHandler(i int)  {
	w.Add(1)
	go func(x int) {
		defer w.Done()
		do(i)
	}(i)
	/*for i := 1; i <= 10; i++ {
		w.Add(1)
		go func(x int) {
			defer w.Done()
			do()
		}(i)
	}*/
	w.Wait()
}

func do(i int) {
	fmt.Println(i, " " + time.Now().String())
}
