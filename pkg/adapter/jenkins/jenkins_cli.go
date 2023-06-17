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

package jenkins

/*
* @author Damon
* @date   2023/5/25 19:09
 */

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/bndr/gojenkins"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// var httpCli *http.Client
var once sync.Once
var instance atomic.Value

var j *gojenkins.Jenkins

/*type JenkinsObj struct {
	ClientSet *gojenkins.Jenkins
}*/

func NewJenkinsCli(baseUrl, user, token string) *gojenkins.Jenkins {
	once.Do(func() {
		// httpCli = &http.Client{Timeout: 3 * time.Second}
		j = gojenkins.CreateJenkins(nil, baseUrl, user, token)
		// instance.Store(&j)
		hlog.Debug("cli: ", j.Server)
	})
	return j // instance.Load().(*gojenkins.Jenkins)
}

func CreateView(ctx context.Context, name string) (*gojenkins.View, error) {
	return j.CreateView(ctx, name, gojenkins.LIST_VIEW)
}

/*func  DeleteView(ctx context.Context, name string) (*gojenkins.View, error) {
	return j.Delete(ctx, name, gojenkins.LIST_VIEW)
}*/

func CreateJob(ctx context.Context, config string, options interface{}) (*gojenkins.Job, error) {
	return j.CreateJob(ctx, config, options)
}

func DeleteJob(ctx context.Context, name string) (bool, error) {
	return j.DeleteJob(ctx, name)
}

func BuildJob(ctx context.Context, name string, params map[string]string) (*gojenkins.Build, error) {
	queueId, err := j.BuildJob(ctx, name, params)
	if err != nil {
		return nil, err
	}
	job, _ := j.GetBuildFromQueueID(ctx, queueId)
	return job, nil
}

func GetJob(ctx context.Context, name string) (*gojenkins.Job, error) {
	job, err := j.GetJob(ctx, name)
	if err != nil {
		hlog.Error(err)
		return nil, err
	}
	// build, err:= job.GetLastBuild(ctx)
	// content,err := build.GetConsoleOutputFromIndex(ctx, 0)
	return job, nil
}

// UpdateJob Update a job. If a job is existed, update its config
func UpdateJob(ctx context.Context, job string, config string) *gojenkins.Job {
	return j.UpdateJob(ctx, job, config)
}
