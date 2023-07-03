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

package bench

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

/*
* @author Damon
* @date   2023/7/3 9:25
 */

//work1

type W1 struct {
	WgSend       *sync.WaitGroup
	Wg           *sync.WaitGroup
	MaxNum       int
	Ch           chan string
	DispatchStop chan struct{}
}

func (w *W1) Dispatch(job string) {
	w.WgSend.Add(10 * w.MaxNum)
	for i := 0; i < 10*w.MaxNum; i++ {
		go func(i int) {
			defer w.WgSend.Done()

			select {
			case w.Ch <- fmt.Sprintf("%d", i):
				return
			case <-w.DispatchStop:
				logrus.Debugln("退出发送 job: ", fmt.Sprintf("%d", i))
				return
			}
		}(i)
	}
}

func (w *W1) StartPool() {
	if w.Ch == nil {
		w.Ch = make(chan string, w.MaxNum)
	}

	if w.WgSend == nil {
		w.WgSend = &sync.WaitGroup{}
	}

	if w.Wg == nil {
		w.Wg = &sync.WaitGroup{}
	}

	if w.DispatchStop == nil {
		w.DispatchStop = make(chan struct{})
	}

	w.Wg.Add(w.MaxNum)
	for i := 0; i < w.MaxNum; i++ {
		go func() {
			defer w.Wg.Done()
			for v := range w.Ch {
				logrus.Debugf("完成工作: %s \n", v)
			}
		}()
	}
}

func (w *W1) Stop() {
	close(w.DispatchStop)
	w.WgSend.Wait()

	close(w.Ch)
	w.Wg.Wait()
}

func DealW1(max int) {
	w := NewWorker(w1, max)
	w.StartPool()
	w.Dispatch("")

	w.Stop()
}



//work2

type SubWorker struct {
	JobChan chan string
}

func (sw *SubWorker) Run(wg *sync.WaitGroup, poolCh chan chan string, quitCh chan struct{}) {
	if sw.JobChan == nil {
		sw.JobChan = make(chan string)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			poolCh <- sw.JobChan

			select {
			case res := <-sw.JobChan:
				logrus.Debugf("完成工作: %s \n", res)

			case <-quitCh:
				logrus.Debugf("消费者结束...... \n")
				return

			}
		}
	}()
}

type W2 struct {
	SubWorkers []SubWorker
	Wg         *sync.WaitGroup
	MaxNum     int
	ChPool     chan chan string
	QuitChan   chan struct{}
}

func (w *W2) Dispatch(job string) {
	jobChan := <-w.ChPool

	select {
	case jobChan <- job:
		logrus.Debugf("发送任务 : %s 完成 \n", job)
		return

	case <-w.QuitChan:
		logrus.Debugf("发送者（%s）结束 \n", job)
		return

	}
}

func (w *W2) StartPool() {
	if w.ChPool == nil {
		w.ChPool = make(chan chan string, w.MaxNum)
	}

	if w.SubWorkers == nil {
		w.SubWorkers = make([]SubWorker, w.MaxNum)
	}

	if w.Wg == nil {
		w.Wg = &sync.WaitGroup{}
	}

	for i := 0; i < len(w.SubWorkers); i++ {
		w.SubWorkers[i].Run(w.Wg, w.ChPool, w.QuitChan)
	}
}

func (w *W2) Stop() {
	close(w.QuitChan)
	w.Wg.Wait()

	close(w.ChPool)
}

func DealW2(max int) {
	w := NewWorker(w2, max)
	w.StartPool()
	for i := 0; i < 10*max; i++ {
		go w.Dispatch(fmt.Sprintf("%d", i))
	}

	w.Stop()
}


//

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
}

type Work interface {
	Dispatch(job string)
	StartPool()
	Stop()
}

const (
	w1 = "w1"
	w2 = "w2"
)

// NewWorker doc
func NewWorker(t string, max int) Work {
	switch t {
	case w1:
		return &W1{
			MaxNum: max,
			Wg:     &sync.WaitGroup{},
			Ch:     make(chan string),
		}

	case w2:
		return &W2{
			Wg:       &sync.WaitGroup{},
			MaxNum:   max,
			ChPool:   make(chan chan string, max),
			QuitChan: make(chan struct{}),
		}

	}

	return nil

}
