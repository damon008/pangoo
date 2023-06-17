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
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or impliepj.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package infra

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"pangoo/katalyst/biz/model/analyse"
	"pangoo/pkg/conf"
	"pangoo/pkg/util/singleton"
	"strings"
	"sync"
)

/*
* @author Damon
* @date   2023/5/15 10:10
 */

func GetProjectInfo(appName, product string) *analyse.Project {
	var wg sync.WaitGroup
	//resultChan := make(chan *analyse.Project)
	//doneChan := make(chan struct{})
	var result *analyse.Project
	var p []*analyse.Project
	//http://192.168.25.116:8108/api/v4/search?scope=projects&search=quec-ota&private_token=txCUZ6GZQo5YTEGbrHhW
	url := fmt.Sprintf("%s/search?scope=projects&search=%s&private_token=%s", conf.EnvConfig.Config.GitConfig.ProjectUrl, appName, conf.EnvConfig.Config.GitConfig.Token)
	sp, err := singleton.HttpDo(url, "GET")
	hlog.Debugf("search projects resp body: %s", string(sp.Body()))
	if err = sonic.Unmarshal(sp.Body(), &p); err != nil {
		hlog.Error(err)
		return &analyse.Project{}
	}

	for _, d := range p {
		wg.Add(1)
		go func(pj *analyse.Project) {
			defer wg.Done()
			ur:= pj.WebURL
			name :=pj.Name
			if (strings.Contains(ur, strings.ToLower(product)) || strings.Contains(ur, strings.ToUpper(product))) {
				if name == appName {
					hlog.Debug(name)
					/*resultChan <- &analyse.Project{
						ID:            pj.ID,
						Name:          pj.Name,
						HTTPURLToRepo: pj.HTTPURLToRepo,
						SSHURLToRepo:  pj.SSHURLToRepo,
						WebURL: pj.WebURL,
					}*/
					result = &analyse.Project{
						ID:            pj.ID,
						Name:          name,
						HTTPURLToRepo: pj.HTTPURLToRepo,
						SSHURLToRepo:  pj.SSHURLToRepo,
						WebURL: pj.WebURL,
					}
					return
				}
			}
		}(d)
	}

	// 启动 Goroutine 监听结果通道和完成通道
	/*go func() {
		for {
			select {
			case data := <-resultChan:
				// 找到相同的数据，设置标志和结果，并发送通知到完成通道
				result = data
				hlog.Debug(result)
				doneChan <- struct{}{}
			case <-doneChan:
				// 任务完成，退出循环
				return
			}
		}
	}()*/

	wg.Wait()
	//close(resultChan)
	return result
}
