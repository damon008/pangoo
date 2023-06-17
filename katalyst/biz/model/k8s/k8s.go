// Code generated by thriftgo (0.2.9). DO NOT EDIT.

package k8s

import (
	"context"
	"fmt"
)

type DeployReq struct {
	AppName    string  `thrift:"appName,1" form:"appName" json:"appName" query:"appName"`
	Tag        string  `thrift:"tag,2" form:"tag" json:"tag" query:"tag"`
	Env        string  `thrift:"env,3" form:"env" json:"env" query:"env"`
	Department string  `thrift:"department,4" form:"department" json:"department" query:"department"`
	Product    string  `thrift:"product,5" form:"product" json:"product" query:"product"`
	Port       *int16  `thrift:"port,6,optional" form:"port" json:"port,omitempty" query:"port"`
	Cmd        *string `thrift:"cmd,7,optional" form:"cmd" json:"cmd,omitempty" query:"cmd"`
}

func NewDeployReq() *DeployReq {
	return &DeployReq{}
}

func (p *DeployReq) GetAppName() (v string) {
	return p.AppName
}

func (p *DeployReq) GetTag() (v string) {
	return p.Tag
}

func (p *DeployReq) GetEnv() (v string) {
	return p.Env
}

func (p *DeployReq) GetDepartment() (v string) {
	return p.Department
}

func (p *DeployReq) GetProduct() (v string) {
	return p.Product
}

var DeployReq_Port_DEFAULT int16

func (p *DeployReq) GetPort() (v int16) {
	if !p.IsSetPort() {
		return DeployReq_Port_DEFAULT
	}
	return *p.Port
}

var DeployReq_Cmd_DEFAULT string

func (p *DeployReq) GetCmd() (v string) {
	if !p.IsSetCmd() {
		return DeployReq_Cmd_DEFAULT
	}
	return *p.Cmd
}

func (p *DeployReq) IsSetPort() bool {
	return p.Port != nil
}

func (p *DeployReq) IsSetCmd() bool {
	return p.Cmd != nil
}

func (p *DeployReq) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DeployReq(%+v)", *p)
}

type BaseResp struct {
	Code int16  `thrift:"code,1" form:"code" json:"code" query:"code"`
	Msg  string `thrift:"msg,2" form:"msg" json:"msg" query:"msg"`
	Data string `thrift:"data,3" form:"data" json:"data" query:"data"`
}

func NewBaseResp() *BaseResp {
	return &BaseResp{}
}

func (p *BaseResp) GetCode() (v int16) {
	return p.Code
}

func (p *BaseResp) GetMsg() (v string) {
	return p.Msg
}

func (p *BaseResp) GetData() (v string) {
	return p.Data
}

func (p *BaseResp) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaseResp(%+v)", *p)
}

type K8sApi interface {
	CreateDeployment(ctx context.Context, deployReq *DeployReq) (r *BaseResp, err error)

	GetDeploymentList(ctx context.Context, namespace string) (r *BaseResp, err error)

	GetDeploymentByName(ctx context.Context, namespace string, name string) (r *BaseResp, err error)

	UpdateDeployment(ctx context.Context, deployReq *DeployReq) (r *BaseResp, err error)
	//重启pod
	DeletePod(ctx context.Context, deployReq *DeployReq) (r *BaseResp, err error)

	RollBack(ctx context.Context, deployReq *DeployReq) (r *BaseResp, err error)
}
