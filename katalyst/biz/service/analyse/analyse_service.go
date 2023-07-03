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

package analyse

import (
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"net/http"
	"os"
	"pangoo/katalyst/biz/model/analyse"
	"pangoo/pkg/util/exec"

	//"os/exec"
	"pangoo/katalyst/biz/infra"
	"pangoo/pkg/conf"
	"pangoo/pkg/response"
	"sync"

	"pangoo/pkg/util/singleton"
	"strings"
)

/*
* @author Damon
* @date   2023/5/12 15:15
 */

func CreateBranch(info *analyse.ProjectInfo) (*response.BaseResp,error) {
	var err error
	var res *protocol.Response
	project := infra.GetProjectInfo(info.AppName, info.Product)
	if project ==nil {
		return &response.BaseResp{
			Code: -3,
			Msg:  "not exist",
			Data: nil,
		}, err
	}
	url := fmt.Sprintf("%s/projects/%d/repository/branches?branch=%s&ref=%s&private_token=%s", conf.EnvConfig.Config.GitConfig.ProjectUrl, project.ID, info.Branch, info.Ref, conf.EnvConfig.Config.GitConfig.Token)
	if res, err = singleton.HttpDo(url, "POST"); err == nil {
		if res.StatusCode() == http.StatusOK ||
			res.StatusCode() == http.StatusCreated {
			return &response.BaseResp{
				Code: 0,
				Msg:  "success",
				Data: info.Branch,
			}, nil
		}
		return &response.BaseResp{
			Code: -1,
			Msg:  "failed",
			Data: string(res.Body()),
		}, err
	}
	hlog.Error("CreateBranch: ", err)
	return &response.BaseResp{
		Code: -1,
		Msg:  "failed",
		Data: "",
	}, err
}

func CreateTag(info *analyse.ProjectInfo) (*response.BaseResp,error) {
	var err error
	var res *protocol.Response
	project := infra.GetProjectInfo(info.AppName, info.Product)
	if project ==nil {
		return &response.BaseResp{
			Code: -3,
			Msg:  "not exist",
			Data: nil,
		}, err
	}
	url := fmt.Sprintf("%s/projects/%d/repository/tags?ref=%s&tag_name=%s&private_token=%s", conf.EnvConfig.Config.GitConfig.ProjectUrl, project.ID, info.Branch, info.Tag, conf.EnvConfig.Config.GitConfig.Token)
	if res, err = singleton.HttpDo(url, "POST"); err == nil {
		if res.StatusCode() == http.StatusOK ||
			res.StatusCode() == http.StatusCreated {
			return &response.BaseResp{
				Code: 0,
				Msg:  "success",
				Data: info.Tag,
			}, nil
		}
		return &response.BaseResp{
			Code: -1,
			Msg:  "failed",
			Data: string(res.Body()),
		}, nil
	}
	hlog.Error("CreateTag: ", err)
	return &response.BaseResp{
		Code: -1,
		Msg:  "failed",
		Data: "",
	}, err
}

func CompareInfo(appName, product, from, to string) (*response.BaseResp, error) {
	project := infra.GetProjectInfo(appName, product)
	if project ==nil {
		return &response.BaseResp{
			Code: -3,
			Msg:  "not exist",
			Data: nil,
		}, nil
	}
	hlog.Debug("project: ", project)
	if (project != nil) {
		url := fmt.Sprintf("%s/projects/%d/repository/commits/%s?private_token=%s", conf.EnvConfig.Config.GitConfig.ProjectUrl, project.ID, to, conf.EnvConfig.Config.GitConfig.Token)
		res, err := singleton.HttpDo(url, "GET")
		hlog.Debug("master data: ", string(res.Body()), "err: ", err)
		var commitID string
		meta := struct {
			Id string `json:"id"`
			ShortId string `json:"short_id"`
		}{}
		if err = sonic.Unmarshal(res.Body(), &meta); err ==nil && strings.Contains(string(res.Body()), "id") {
			commitID = meta.Id
		} else {
			hlog.Error(err)
			return &response.BaseResp{
				Code: -1,
				Msg:  "failed",
				Data: string(res.Body()),
			}, err
		}
		url = fmt.Sprintf("%s/projects/%d/repository/commits?ref_name=%s&private_token=%s", conf.EnvConfig.Config.GitConfig.ProjectUrl, project.ID, from, conf.EnvConfig.Config.GitConfig.Token)
		res, err = singleton.HttpDo(url, "GET")
		hlog.Debug("branch commits data: ", string(res.Body()), "err: ", err, "commitID: ", commitID)
		if strings.Contains(string(res.Body()), commitID) {
			return &response.BaseResp{
				Code: 0,
				Msg:  "success",
				Data: "包括该版本",
			}, nil
		}
		return &response.BaseResp{
			Code: -1,
			Msg:  "failed",
			Data: string(res.Body()),
		}, err
	}
	return &response.BaseResp{
		Code: -1,
		Msg:  "failed",
		Data: "未找到app项目",
	}, nil
}

func MergeReq(info *analyse.ProjectInfo) (*response.BaseResp, error) {
	project := infra.GetProjectInfo(info.AppName, info.Product)
	if project ==nil {
		return &response.BaseResp{
			Code: -3,
			Msg:  "not exist",
			Data: nil,
		}, nil
	}
	url := fmt.Sprintf("%s/projects/%d/merge_requests?private_token=%s", conf.EnvConfig.Config.GitConfig.ProjectUrl, project.ID, conf.EnvConfig.Config.GitConfig.Token)
	//hlog.Info(url)
	data := struct {
		Id int16 `json:"id"`
		SourceBranch string `json:"source_branch"`
		TargetBranch string `json:"target_branch"`
		Title string `json:"title"`
		RemoveSourceBranch bool `json:"remove_source_branch"`
	}{
		Id:                 project.ID,
		SourceBranch:       info.Branch,
		TargetBranch:       info.Ref,
		Title:              "Merge Request",
		RemoveSourceBranch: false,
	}
	res, err := singleton.Json4Post(url, &data)
	hlog.Error("merge req: ", err)
	if res.StatusCode() != http.StatusCreated {
		return &response.BaseResp{
			Code: -1,
			Msg:  "failed",
			Data: string(res.Body()),
		}, err//errno.NewErrNo()
	}
	meta := struct {
		Iid int `json:"iid"`
	}{}
	if err = sonic.Unmarshal(res.Body(), &meta); err ==nil {
		//合并合并请求，无法合并到受保护的分支
		url = fmt.Sprintf("%s/projects/%d/merge_requests/%d/merge?private_token=%s", conf.EnvConfig.Config.GitConfig.ProjectUrl, project.ID, meta.Iid, conf.EnvConfig.Config.GitConfig.Token)
		resp, err := singleton.HttpDo(url, "PUT")
		hlog.Info(resp.StatusCode())
		if err !=nil || resp.StatusCode() != http.StatusOK || resp.StatusCode() != http.StatusCreated {
			hlog.Error(resp.StatusCode())
			hlog.Error(err)
			return &response.BaseResp{
				Code: -1,
				Msg:  "failed",
				Data: string(resp.Body()),
			}, err
		}
		return &response.BaseResp{
			Code: 0,
			Msg:  "success",
			Data: string(resp.Body()),
		}, nil
	}
	return &response.BaseResp{
		Code: -1,
		Msg:  "failed",
		Data: string(res.Body()),
	}, err//errno.NewErrNo()
}

func CloneCode(urls []string)  {
	var wg sync.WaitGroup
	//for i := 0;i < len(urls); i++ {
	for _,url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			/*project := infra.GetProjectInfo(url)
			url = fmt.Sprintf("%s/projects/%d/repository/files/%s?private_token=%s&ref=2.32.2-dev", conf.EnvConfig.Config.GitConfig.ProjectUrl, project.ID, "src%2F", conf.EnvConfig.Config.GitConfig.Token)
			resp, err := singleton.HttpDo(url, "GET")
			if err !=nil {
				hlog.Error(err)
			}
			hlog.Info(resp.StatusCode())
			hlog.Info(string(resp.Body()))*/
			url = strings.ReplaceAll(url, "http://", fmt.Sprintf("http://oauth2:%s@", conf.EnvConfig.Config.GitConfig.Token))
			/*err := exec.Cmd("","apt-get", "update")
			if err != nil {
				hlog.Error(err)
				return
			}

			err = exec.Cmd("","apt-get", "install", "-y", "git")
			if err != nil {
				hlog.Error(err)
				return
			}*/

			if err := os.MkdirAll("E:/data/projects/code", 0o777); err != nil {
				hlog.Error(err)
			}
			err := exec.Cmd("E:/data/projects/code","git", "clone", "-b", "2.32.2-dev", url)
			if err!=nil {
				hlog.Error(err)
			}
		}()
	}
	wg.Wait()
}
