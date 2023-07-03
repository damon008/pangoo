/*
 * Copyright 2023-present by Damon All Rights Reserved
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpool

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/linclin/grpool"
)

/*
* @author Damon
* @date   2023/7/3 9:32
 */

func main() {
	a := make(chan int, 1)
	b := make(chan int, 1)
	a = b
	a <- 1
	fmt.Println(<-b)
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	// number of workers, size of job queue,job timeout
	pool := grpool.NewPool(10, 100, 3 * time.Second)
	defer pool.Release()

	// how many jobs we should wait
	pool.WaitCount(100)

	// submit one or more jobs to pool
	for i := 0; i < 100; i++ {
		count := i
		pool.JobQueue <- grpool.Job{
			Jobid: count,
			Jobfunc: func() (interface{}, error) {
				// say that job is done, so we can know how many jobs are finished
				SleepRandomDuration()
				fmt.Printf("hello %d\n", count)
				return count, nil
			},
		}
	}

	// wait until we call JobDone for all jobs
	pool.WaitAll()
	for res := range pool.Jobresult {
		fmt.Println("res", res)
	}

}
func SleepRandomDuration() {
	ns := int64(10) * 1000000000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	d := time.Duration(r.Int63n(ns)) * time.Nanosecond
	time.Sleep(d)
}
