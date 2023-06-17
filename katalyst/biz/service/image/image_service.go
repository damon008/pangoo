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

package image

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"pangoo/katalyst/biz/model/image"
	"pangoo/pkg/adapter/docker"
	"pangoo/pkg/conf"
	"sync"
)

/*
* @author Damon
* @date   2023/5/20 10:23
 */

type ImageServe struct {
	Docker  *docker.DockerAgent `thrift:"docker,1" form:"docker" json:"docker" query:"docker"`
	//AuthStr string `thrift:"authStr,2" form:"authStr" json:"authStr" query:"authStr"`
}

type ImageCommitInfo struct {
	Repo           string `thrift:"repo,1" form:"repo" json:"repo" query:"repo"`
	CallBack       string `thrift:"callBack,2" form:"callBack" json:"callBack" query:"callBack"`
	SizeLimit      int8   `thrift:"sizeLimit,3" form:"sizeLimit" json:"sizeLimit" query:"sizeLimit"`
	LayerSizeLimit int8   `thrift:"layerSizeLimit,4" form:"layerSizeLimit" json:"layerSizeLimit" query:"layerSizeLimit"`
	LayerLimit     int8   `thrift:"layerLimit,5" form:"layerLimit" json:"layerLimit" query:"layerLimit"`
	ImageId        int32  `thrift:"imageId,6" form:"imageId" json:"imageId" query:"imageId"`
	PodName        string `thrift:"podName,7" form:"podName" json:"podName" query:"podName"`
}

func NewImageServe() *ImageServe {
	cli := docker.NewDockerAgent(conf.EnvConfig.Config.DockerConfig.Username, conf.EnvConfig.Config.DockerConfig.Password)
	return &ImageServe{
		Docker:  cli,
		//AuthStr: str,
	}
	/*if str,err := docker.Auth(conf.EnvConfig.DockerConfig.Username, conf.EnvConfig.DockerConfig.Password); err ==nil {
		return &ImageServe{
			Docker:  cli,
			//AuthStr: str,
		}
	}*/
	return nil
}

func (is *ImageServe) ImageBuild(ctx context.Context, req image.ImageInfo)  {
	// 使用WaitGroup来等待所有goroutine完成
	var wg sync.WaitGroup
	for _,body := range req.Apps {
		wg.Add(1)
		go func() {
			defer wg.Done()
			is.Docker.ImageBuild(*body.DockerFile, *body.BaseRepo + body.Tag)
			go is.Docker.PushImage("", is.Docker.AuthStr)
		}()
	}
	for i := 0;i < len(req.GetApps()); i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			body := req.GetApps()[i]
			hlog.Info("APP info: ", body)
			is.Docker.ImageBuild("", "")
			go is.Docker.PushImage("", is.Docker.AuthStr)
		}(i)
	}
	// 等待所有goroutine完成
	wg.Wait()
	hlog.Info("所有Docker镜像已完成编译")
}
