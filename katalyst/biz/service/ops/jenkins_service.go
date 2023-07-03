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
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
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

// AddJob 新建任务，任务已存在则更新，mode支持 create(仅创建任务，已存在则跳过), add(创建任务，已存在则更新), update(仅更新任务)，默认add
func AddJob(ctx context.Context, name string, config string, mode ...string) bool {
	jobMode := "add"
	if len(mode) > 0 {
		jobMode = mode[0]
	}

	isExist := JobIsExist(ctx, name)

	if (jobMode == "create" && isExist) || (jobMode == "update" && !isExist) { // create已存在则跳过，update未存在通过
		return false
	}

	if jobMode == "add" && !gstr.ContainsI(name, "common") && gstr.Contains(config, "RepBuild") { // 判断占位符是否全部被替换，没有则返回失败
		hlog.Errorf("%s 任务创建或更新失败，任务配置未初始化完成\n%s", name, config)
		return false
	}

	if isExist { // 已存在则更新，未存在则创建
		return JobUpdate(ctx, name, config)
	} else {
		return JobCreate(ctx, config, name)
	}
}

type AppInfo struct {
	ID      int64  `json:"id"`      // 应用ID，回调devops平台需要传递
	AppName string `json:"appName"` // 应用名
	Status  int    `json:"status"`  // 1 任务构建成功，0 任务构建失败
}

type BuildStatusCallback struct {
	mu      sync.Mutex // 互斥锁，用于保护共享变量
	AppInfo []AppInfo  `json:"appInfo"`
}

func (s *BuildStatusCallback) AppendAppInfo(info AppInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.AppInfo = append(s.AppInfo, info)
}

func BuildJob(ctx context.Context, name string, param map[string]string) int64 {
	buildObj, err := jenkins.BuildJob(ctx, name, param)
	if err != nil {
		hlog.Errorf("任务 [%s] 参数 [%s] 构建失败: %s", name, param, err)
		return 0
	}
	hlog.Debugf("任务 [%s] 参数 [%s] 已触发, 序号 [%d]", name, param, buildObj.GetBuildNumber())
	return buildObj.GetBuildNumber()
}

// GetBuildStatus 返回 1 代表任务成功，0代表失败
func GetBuildStatus(ctx context.Context, name string, buildNo int64) int {
	buildObj, err := jenkins.GetBuild(ctx, name, buildNo)
	if err != nil {
		hlog.Errorf("任务 [%s] 序号 [%d] 获取构建对象失败, %s", name, buildNo, err.Error())
		return 0
	}

	// 堵塞直至任务完成
	for buildObj.IsRunning(ctx) {
		time.Sleep(time.Second * 3)
		_, _ = buildObj.Poll(ctx)
	}

	if buildObj.IsGood(ctx) {
		hlog.Infof("任务 [%s] 序号 [%d] 构建成功", name, buildNo)
		return 1
	}
	hlog.Errorf("任务 [%s] 序号 [%d] 构建失败, 状态 [%s]", name, buildNo, buildObj.GetResult())
	return 0
}

// APIBuildMode jenkins 应用操作模式，用于不同接口触发不同的逻辑
type APIBuildMode int

// 定义构建模式对应的枚举值
const (
	APIBuildModeDeployApp   APIBuildMode = iota + 1 // 部署应用接口模式
	APIBuildModeUpdateApp                           // 更新应用接口模式
	APIBuildModeRestartApp                          // 重启应用接口模式
	APIBuildModeRollbackApp                         // 回滚应用接口模式
)

// 获取应用操作模式枚举值对应的字符串值
func (s APIBuildMode) String() string {
	strArr, iS := []string{"DeployApp", "UpdateApp", "RestartApp", "RollbackApp"}, int(s)
	if iS < 1 || iS > len(strArr) {
		return ""
	}
	return strArr[iS-1]
}

func BuildJobMulti(job *jk.AppDeployReq, buildMode APIBuildMode) (*jk.BaseResp, error) {
	// 创建一个等待组，用于等待所有任务完成
	var wg sync.WaitGroup

	// 设置并发任务数量
	appCount := len(job.AppInfo)
	wg.Add(appCount)

	// 使用信号量控制并发数量
	semaphore := make(chan struct{}, appCount)

	// m 存储各任务的构建序号、buildStatus 存储各个任务是否构建生产
	m := gmap.NewStrStrMap(true) // 并发安全
	buildStatus := &BuildStatusCallback{}

	// 根据接口模式获取当前批次的构建主体信息
	baseJKSInfo := &JobConfig{}
	switch buildMode {
	default:
		hlog.Errorf("不支持的接口模式: %s", buildMode.String())
		return nil, nil
	case APIBuildModeDeployApp:
		baseJKSInfo = NewJobConfig()
	case APIBuildModeUpdateApp:
		baseJKSInfo = NewJobConfig(JKSDeployModeUpdate)
	case APIBuildModeRestartApp:
		baseJKSInfo = NewJobConfig(JKSDeployModeRestart)
	case APIBuildModeRollbackApp:
		baseJKSInfo = NewJobConfig(JKSDeployModeRollback)
	}
	baseJKSInfo.ScanAppDeployReq(job)

	// 启动并发任务
	for _, app := range job.AppInfo {
		go func(app *jk.AppInfo) {
			defer wg.Done()
			// 获取信号量，控制并发数量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			ctx := context.Background()
			// 获取当前批次的构建主体信息及具体任务信息，使用字面量而非指针，防止并发数据异常
			tJKS := *baseJKSInfo
			// erlang 服务只有 cd 部署模式
			tJKS.ScanAppInfo(app)
			if strings.ToLower(tJKS.BuildAppLang) == "erlang" {
				tJKS.ScanAppInfo(app, JKSDeployModeCD)
			} else {
				tJKS.ScanAppInfo(app)
			}
			tJKS.GenerateBuildInfo()
			buildNo := int64(0)

			msg, ok, isDeployed := "", true, true
			switch buildMode {
			case APIBuildModeDeployApp: // 需判断任务新增成功
				ok = AddJob(ctx, tJKS.BuildJobFullName, GetJobActualConfigFromTpl(&tJKS))
				msg = "创建或更新失败，跳过"
			case APIBuildModeUpdateApp: // 需判断任务更新成功，再判断是否部署
				ok = AddJob(ctx, tJKS.BuildJobFullName, GetJobActualConfigFromTpl(&tJKS), "update")
				msg = "更新失败，跳过"
				if ok {
					isDeployed = JobIsDeploy(ctx, tJKS.BuildJobFullName)
					msg = "尚未部署，跳过"
				}
			case APIBuildModeRestartApp, APIBuildModeRollbackApp: // 需判断是否部署
				isDeployed = JobIsDeploy(ctx, tJKS.BuildJobFullName)
				msg = "尚未部署，跳过"
			}

			if !ok {
				hlog.Errorf("任务 [%s] %s", tJKS.BuildAppName, msg)
				m.Set(tJKS.BuildAppName, gconv.String(buildNo))
				return
			}
			if !isDeployed {
				hlog.Warnf("任务 [%s] %s", tJKS.BuildAppName, msg)
				m.Set(tJKS.BuildAppName, gconv.String(buildNo))
				return
			}

			buildNo = BuildJob(ctx, tJKS.BuildJobFullName, tJKS.BuildParam)
			m.Set(tJKS.BuildAppName, gconv.String(buildNo))
		}(app)
	}
	// 等待并发组全部完成
	wg.Wait()

	// 使用单个 goroutine 循环获取构建结果
	go func() {
		for _, v := range job.GetAppInfo() {
			ctx := context.Background()
			if m.Get(v.AppName) == "0" { // 如果序号为0，则不需要回调
				continue
			}
			status := GetBuildStatus(ctx, baseJKSInfo.BuildPrefix+"_"+v.AppName, gconv.Int64(m.Get(v.AppName)))
			buildStatus.AppendAppInfo(AppInfo{
				ID:      v.AppId,
				AppName: v.AppName,
				Status:  status,
			})
		}

		// 存在触发的构建，则回调 devops 平台，发送构建结果
		if len(buildStatus.AppInfo) != 0 {
			if err := CallBackBuildStatus(context.Background(), buildStatus); err != nil {
				hlog.Errorf("回调 devops 平台失败：%s", err.Error())
			}
		}
	}()

	return &jk.BaseResp{
		Code: 0,
		Msg:  "ok",
		Data: gconv.MapStrStr(m),
	}, nil
}

func DeployApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	return BuildJobMulti(job, APIBuildModeDeployApp)
}

func UpdateApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	return BuildJobMulti(job, APIBuildModeUpdateApp)
}

func RestartApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	return BuildJobMulti(job, APIBuildModeRestartApp)
}

func RollbackApp(job *jk.AppDeployReq) (*jk.BaseResp, error) {
	return BuildJobMulti(job, APIBuildModeRollbackApp)
}

func CallBackBuildStatus(ctx context.Context, data *BuildStatusCallback) error {
	jsonStr, err := sonic.MarshalString(data)
	if err != nil {
		hlog.Errorf("json解析失败: %s", err.Error())
		return err
	}
	hlog.Debugf("当前批次任务回调参数: %s", jsonStr)
	req, err := g.Client().Post(ctx, defaultDevOpsOpenApi, jsonStr)
	if err != nil {
		hlog.Errorf("回调请求失败: %s", err.Error())
		return err
	}
	defer func(req *gclient.Response) {
		_ = req.Close()
	}(req)

	if statusCode := req.StatusCode; req.StatusCode != 200 {
		msg := fmt.Sprintf("http请求状态码非200, 当前状态码: %d", statusCode)
		hlog.Error(msg)
		return errors.New(msg)
	}

	resp, err := io.ReadAll(req.Response.Body)
	if err != nil {
		hlog.Errorf("读取构建结果回调响应数据失败", err.Error())
		return err
	}
	hlog.Infof("当前批次任务回调结果, %s", string(resp))
	return nil
}
