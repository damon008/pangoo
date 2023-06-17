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

package deployment

import (
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	v1 "k8s.io/api/apps/v1"
	"pangoo/katalyst/biz/build"
	"pangoo/pkg/adapter/k8s"
)

/*
* @author Damon
* @date   2023/5/10 14:13
 */

//var TaskServe *TaskService

func NewTaskService() *TaskService {
	return &TaskService{
		Deployment:  k8s.NewDeployment(),
		Service:     k8s.NewService(),
		HPA:         k8s.NewHPA(),
		NameSpace:   k8s.NewNamespce(),
	}
	/*TaskServe = &TaskService{
		Deployment:  k8s.NewDeployment(),
		Service:     k8s.NewService(),
		HPA:         k8s.NewHPA(),
		NameSpace:   k8s.NewNamespce(),
	}
	return TaskServe*/
}

type TaskService struct {
	Deployment *k8s.Deployment
	Service    *k8s.Service
	HPA        *k8s.HPA
	NameSpace  *k8s.Namespce

	DeploymentBuilder build.DeploymentBuilder
	ServiceBuilder    build.ServiceBuilder
	HPABuilder        build.HPABuilder

	//IDbService  db.IDbService
	//CheckerRepo IChecker
}

func (p *TaskService) CreateDeployment() interface{} {
	var list *v1.DeploymentList
	var err error
	if list, err = p.Deployment.List("kube-system"); err !=nil {
		hlog.Error("failed: ", err)
		return nil
	}
	l,_ := sonic.MarshalString(list)
	hlog.Info("list: ", l)
	return list

}

func (p *TaskService) GetDeploymentList(ns string) *v1.DeploymentList {
	var list *v1.DeploymentList
	var err error
	if list, err = p.Deployment.List(ns); err !=nil {
		hlog.Error("failed: ", err)
		return nil
	}
	l,_ := sonic.MarshalString(list)
	hlog.Info("list: ", l)
	return list
}

func (p *TaskService) GetDeployment(ns string, name string) (*v1.Deployment, error) {
	return p.Deployment.Get(ns, name)
}
