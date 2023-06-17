/*
 * Copyright 2023-2100 by Damon All Rights Reserved
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

package ops

import (
	"context"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"

	jk "pangoo/katalyst/biz/model/ops"
	"pangoo/pkg/adapter/jenkins"
)

func CreateView(name string) (*jk.BaseResp, error) {

	if _, err := jenkins.CreateView(context.Background(), name); err != nil {
		return &jk.BaseResp{
			Code: -1,
			Msg:  "failed: " + err.Error(),
			Data: nil,
		}, err
	}
	return &jk.BaseResp{
		Code: 0,
		Msg:  "success",
		Data: nil,
	}, nil
}

/*func BuildLog(job *jk.JobReq) (*jk.BaseResp, error) {
	name := job.Department + "_" + job.Product + "_" + job.Env + "_" + job.AppName
	url := "http://jenkins-cloud.quectel.com/view/" + job.Department + "_" + job.Product +"/job/" + job.Department + "_" + job.Product + "_" + job.Env + "_" + job.AppName + "/"+ string(job.BuildNo) + "/consoleText"
}*/

func DeleteJob(job *jk.JobReq) (*jk.BaseResp, error) {
	name := job.Department + "_" + job.Product + "_" + job.Env + "_" + job.AppName
	if _, err := jenkins.DeleteJob(context.Background(), name); err != nil {
		return &jk.BaseResp{
			Code: -1,
			Msg:  err.Error(),
			Data: nil,
		}, err
	}
	return &jk.BaseResp{
		Code: 0,
		Msg:  "success",
		Data: nil,
	}, nil
}

// JobIsDeploy 判断任务是否成功发布过，用于 update、restart、rollback 等模式
func JobIsDeploy(ctx context.Context, name string) bool {
	job, err := jenkins.GetJob(ctx, name)
	if err != nil {
		return false
	}
	_, err = job.GetLastSuccessfulBuild(ctx)
	return err == nil
}

func JobIsExist(ctx context.Context, name string) bool {
	_, err := jenkins.GetJob(ctx, name)
	return err == nil
}

func JobUpdate(ctx context.Context, name string, config string) bool {
	job, _ := jenkins.GetJob(ctx, name)
	cConfig, _ := job.GetConfig(ctx)
	if cConfig == config {
		hlog.Debugf("[%s] 任务已存在, 且配置和当前配置一致, 无需变更!", name)
		return true
	}

	if jenkins.UpdateJob(ctx, name, config) == nil { // nil 则更新失败
		hlog.Errorf("%s 任务更新失败，未知错误", name)
		return false
	}
	hlog.Debugf("%s 任务更新成功", name)
	return true
}

func JobCreate(ctx context.Context, config string, name string) bool {
	job, err := jenkins.CreateJob(ctx, config, name)
	if err != nil || job == nil {
		hlog.Errorf("[%s] 任务创建失败!", name, err.Error())
		return false
	}
	hlog.Debugf("[%s] 任务创建成功", name)
	return true
}

// AddJob 新建任务，任务已存在则更新
func AddJob(ctx context.Context, name string, config string) bool {
	if !gstr.ContainsI(name, "common") && gstr.Contains(config, "RepBuild") { // 判断占位符是否全部被替换，没有则返回失败
		hlog.Errorf("%s 任务创建或更新失败，任务配置未初始化完成\n%s", name, config)
		return false
	}

	if JobIsExist(ctx, name) {
		return JobUpdate(ctx, name, config)
	}

	return JobCreate(ctx, config, name)
}

// BuildJob 构建任务，返回构建序号与成功与否
func BuildJob(ctx context.Context, name string, params map[string]string) (int64, bool) {
	buildNo := int64(0)
	build, err := jenkins.BuildJob(ctx, name, params)
	if err != nil {
		hlog.Errorf("任务 [%s] 参数 [%s] 构建失败: %s", name, params, err)
		return buildNo, false
	}

	buildNo = build.GetBuildNumber()
	// for build.IsRunning(ctx) {
	// 	time.Sleep(time.Second * 3)
	// 	_, _ = build.Poll(ctx)
	// }

	// if build.IsGood(ctx) {
	// 	hlog.Infof("任务 [%s] 参数 [%s] 构建成功，序号: %d", name, params, buildNo)
	// 	return buildNo, true
	// }
	// hlog.Errorf("任务 [%s] 参数 [%s] 构建失败，序号: %d", name, params, buildNo)
	return buildNo, true
}

func DeployApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	// 创建一个等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 设置并发任务数量
	concurrency := len(job.AppInfo)
	wg.Add(concurrency)

	// 使用信号量控制并发数量
	semaphore := make(chan struct{}, concurrency)

	// m 存储各任务的构建序号
	m := make(map[string]string)

	// 获取当前批次的构建主体信息
	baseJKSInfo := NewJobConfig()
	baseJKSInfo.ScanAppDeployReq(job)

	// 启动并发任务
	for i := 0; i < len(job.AppInfo); i++ {
		go func(index int) {
			defer wg.Done()
			// 获取信号量，控制并发数量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ctx := context.Background()
			// 获取当前批次的构建主体信息及具体任务信息
			thisJKSInfo := baseJKSInfo
			thisJKSInfo.ScanAppInfo(job.AppInfo[index])
			jobName := job.Department + "_" + job.Product + "_" + job.Env + "_" + thisJKSInfo.BuildAppName
			buildNo := int64(0)

			// 新增对应任务，失败则跳过
			if !AddJob(ctx, jobName, GetJobActualConfig(thisJKSInfo)) {
				m[thisJKSInfo.BuildAppName] = gconv.String(buildNo)
				hlog.Debugf("任务 [%s] 创建或更新失败", jobName)
				return
			}

			// 将结构体转换为构建所需参数，使用 jks 的tag进行转换
			paramMap := StructToMapWithTag(*thisJKSInfo, "param")
			hlog.Debugf("任务 [%s] 构建参数 [%s]", jobName, paramMap)
			buildNo, _ = BuildJob(ctx, jobName, paramMap)
			m[thisJKSInfo.BuildAppName] = gconv.String(buildNo)
		}(i)
	}
	// 等待所有任务完成
	wg.Wait()

	return &jk.BaseResp{
		Code: 0,
		Msg:  "当前批次构建完成",
		Data: m,
	}, nil
}

func UpdateApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	// 创建一个等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 设置并发任务数量
	concurrency := len(job.AppInfo)
	wg.Add(concurrency)

	// 使用信号量控制并发数量
	semaphore := make(chan struct{}, concurrency)

	// m 存储各任务的构建序号
	m := make(map[string]string)

	// 获取当前批次的构建主体信息
	baseJKSInfo := NewJobConfig("update")
	baseJKSInfo.ScanAppDeployReq(job)

	// 启动并发任务
	for i := 0; i < len(job.AppInfo); i++ {
		go func(index int) {
			defer wg.Done()
			// 获取信号量，控制并发数量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ctx := context.Background()
			// 获取当前批次的构建主体信息及具体任务信息
			thisJKSInfo := baseJKSInfo
			thisJKSInfo.ScanAppInfo(job.AppInfo[index])
			jobName := job.Department + "_" + job.Product + "_" + job.Env + "_" + thisJKSInfo.BuildAppName
			buildNo := int64(0)

			// 更新对应任务
			if !AddJob(ctx, jobName, GetJobActualConfig(thisJKSInfo)) {
				m[thisJKSInfo.BuildAppName] = gconv.String(buildNo)
				return
			}

			// 判断是否已部署，已部署则对集群进行更新
			if JobIsDeploy(ctx, jobName) {
				buildNo, _ = BuildJob(ctx, jobName, StructToMapWithTag(*thisJKSInfo, "param"))
			}
			m[thisJKSInfo.BuildAppName] = gconv.String(buildNo)
		}(i)
	}
	// 等待所有任务完成
	wg.Wait()

	return &jk.BaseResp{
		Code: 0,
		Msg:  "ok",
		Data: m,
	}, nil
}

func RestartApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	// 创建一个等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 设置并发任务数量
	concurrency := len(job.AppInfo)
	wg.Add(concurrency)

	// 使用信号量控制并发数量
	semaphore := make(chan struct{}, concurrency)

	// m 存储各任务的构建序号
	m := make(map[string]string)

	// 获取当前批次的构建主体信息
	baseJKSInfo := NewJobConfig("restart") // restart 构建模式
	baseJKSInfo.ScanAppDeployReq(job)

	// 启动并发任务
	for i := 0; i < len(job.AppInfo); i++ {
		go func(index int) {
			defer wg.Done()
			// 获取信号量，控制并发数量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ctx := context.Background()
			// 获取当前批次的构建主体信息及具体任务信息
			thisJKSInfo := baseJKSInfo
			thisJKSInfo.ScanAppInfo(job.AppInfo[index])
			jobName := job.Department + "_" + job.Product + "_" + job.Env + "_" + thisJKSInfo.BuildAppName

			// 判断是否已部署，已部署则对集群进行更新
			if JobIsDeploy(ctx, jobName) {
				buildNo, _ := BuildJob(ctx, jobName, StructToMapWithTag(*thisJKSInfo, "param"))
				m[thisJKSInfo.BuildAppName] = gconv.String(buildNo)
			}
		}(i)
	}
	// 等待所有任务完成
	wg.Wait()

	return &jk.BaseResp{
		Code: 0,
		Msg:  "ok",
		Data: m,
	}, nil
}

func RollbackApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	// 创建一个等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 设置并发任务数量
	concurrency := len(job.AppInfo)
	wg.Add(concurrency)

	// 使用信号量控制并发数量
	semaphore := make(chan struct{}, concurrency)

	// m 存储各任务的构建序号
	m := make(map[string]string)

	// 获取当前批次的构建主体信息
	baseJKSInfo := NewJobConfig("rollback") // rollback 构建模式
	baseJKSInfo.ScanAppDeployReq(job)

	// 启动并发任务
	for i := 0; i < len(job.AppInfo); i++ {
		go func(index int) {
			defer wg.Done()
			// 获取信号量，控制并发数量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ctx := context.Background()
			// 获取当前批次的构建主体信息及具体任务信息
			thisJKSInfo := baseJKSInfo
			thisJKSInfo.ScanAppInfo(job.AppInfo[index])
			jobName := job.Department + "_" + job.Product + "_" + job.Env + "_" + thisJKSInfo.BuildAppName
			buildNo := int64(0)

			// 判断是否已部署，已部署则对集群进行更新
			if JobIsDeploy(ctx, jobName) {
				buildNo, _ = BuildJob(ctx, jobName, StructToMapWithTag(*thisJKSInfo, "param"))
			}
			m[thisJKSInfo.BuildAppName] = gconv.String(buildNo)
		}(i)
	}
	// 等待所有任务完成
	wg.Wait()

	return &jk.BaseResp{
		Code: 0,
		Msg:  "ok",
		Data: m,
	}, nil
}
