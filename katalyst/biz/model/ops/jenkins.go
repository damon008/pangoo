// Code generated by thriftgo (0.2.11). DO NOT EDIT.

package ops

import (
	"context"
	"fmt"
)

type JenkinsJobConfig struct {
	JobName              string            `thrift:"jobName,1" form:"jobName" json:"jobName" query:"jobName"`
	Env                  string            `thrift:"env,2" form:"env" json:"env" query:"env"`
	ConfigTplName        string            `thrift:"configTplName,3" form:"configTplName" json:"configTplName" query:"configTplName"`
	ConfigTplCtx         string            `thrift:"configTplCtx,4" form:"configTplCtx" json:"configTplCtx" query:"configTplCtx"`
	ConfigCtx            string            `thrift:"configCtx,5" form:"configCtx" json:"configCtx" query:"configCtx"`
	AppName              string            `thrift:"appName,6" form:"appName" json:"appName" query:"appName"`
	CodeUrl              string            `thrift:"codeUrl,7" form:"codeUrl" json:"codeUrl" query:"codeUrl"`
	AliPipelineUrl       *string           `thrift:"aliPipelineUrl,8,optional" form:"aliPipelineUrl" json:"aliPipelineUrl,omitempty" query:"aliPipelineUrl"`
	AliPipelineShellName *string           `thrift:"aliPipelineShellName,9,optional" form:"aliPipelineShellName" json:"aliPipelineShellName,omitempty" query:"aliPipelineShellName"`
	ConfigTplMap         map[string]string `thrift:"configTplMap,10,optional" form:"configTplMap" json:"configTplMap,omitempty" query:"configTplMap"`
}

func NewJenkinsJobConfig() *JenkinsJobConfig {
	return &JenkinsJobConfig{}
}

func (p *JenkinsJobConfig) GetJobName() (v string) {
	return p.JobName
}

func (p *JenkinsJobConfig) GetEnv() (v string) {
	return p.Env
}

func (p *JenkinsJobConfig) GetConfigTplName() (v string) {
	return p.ConfigTplName
}

func (p *JenkinsJobConfig) GetConfigTplCtx() (v string) {
	return p.ConfigTplCtx
}

func (p *JenkinsJobConfig) GetConfigCtx() (v string) {
	return p.ConfigCtx
}

func (p *JenkinsJobConfig) GetAppName() (v string) {
	return p.AppName
}

func (p *JenkinsJobConfig) GetCodeUrl() (v string) {
	return p.CodeUrl
}

var JenkinsJobConfig_AliPipelineUrl_DEFAULT string

func (p *JenkinsJobConfig) GetAliPipelineUrl() (v string) {
	if !p.IsSetAliPipelineUrl() {
		return JenkinsJobConfig_AliPipelineUrl_DEFAULT
	}
	return *p.AliPipelineUrl
}

var JenkinsJobConfig_AliPipelineShellName_DEFAULT string

func (p *JenkinsJobConfig) GetAliPipelineShellName() (v string) {
	if !p.IsSetAliPipelineShellName() {
		return JenkinsJobConfig_AliPipelineShellName_DEFAULT
	}
	return *p.AliPipelineShellName
}

var JenkinsJobConfig_ConfigTplMap_DEFAULT map[string]string

func (p *JenkinsJobConfig) GetConfigTplMap() (v map[string]string) {
	if !p.IsSetConfigTplMap() {
		return JenkinsJobConfig_ConfigTplMap_DEFAULT
	}
	return p.ConfigTplMap
}

func (p *JenkinsJobConfig) IsSetAliPipelineUrl() bool {
	return p.AliPipelineUrl != nil
}

func (p *JenkinsJobConfig) IsSetAliPipelineShellName() bool {
	return p.AliPipelineShellName != nil
}

func (p *JenkinsJobConfig) IsSetConfigTplMap() bool {
	return p.ConfigTplMap != nil
}

func (p *JenkinsJobConfig) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("JenkinsJobConfig(%+v)", *p)
}

type BaseReq struct {
	Department string `thrift:"department,1" form:"department" json:"department" query:"department"`
	Product    string `thrift:"product,2" form:"product" json:"product" query:"product"`
}

func NewBaseReq() *BaseReq {
	return &BaseReq{}
}

func (p *BaseReq) GetDepartment() (v string) {
	return p.Department
}

func (p *BaseReq) GetProduct() (v string) {
	return p.Product
}

func (p *BaseReq) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaseReq(%+v)", *p)
}

type AppInfo struct {
	AppName     string              `thrift:"appName,1" form:"appName" json:"appName" query:"appName"`
	GitUrl      string              `thrift:"gitUrl,2" form:"gitUrl" json:"gitUrl" query:"gitUrl"`
	Branch      string              `thrift:"branch,3" form:"branch" json:"branch" query:"branch"`
	Language    string              `thrift:"language,4" form:"language" json:"language" query:"language"`
	CPU         float64             `thrift:"CPU,5" form:"CPU" json:"CPU" query:"CPU"`
	Memory      int16               `thrift:"memory,6" form:"memory" json:"memory" query:"memory"`
	Replicas    int8                `thrift:"replicas,7" form:"replicas" json:"replicas" query:"replicas"`
	MinReplicas int8                `thrift:"minReplicas,8" form:"minReplicas" json:"minReplicas" query:"minReplicas"`
	MaxReplicas int8                `thrift:"maxReplicas,9" form:"maxReplicas" json:"maxReplicas" query:"maxReplicas"`
	PortMap     []map[string]string `thrift:"portMap,10" form:"portMap" json:"portMap" query:"portMap"`
	BuildNo     int64               `thrift:"buildNo,11" form:"buildNo" json:"buildNo" query:"buildNo"`
	AppId       int64               `thrift:"appId,12" form:"appId" json:"appId" query:"appId"`
}

func NewAppInfo() *AppInfo {
	return &AppInfo{}
}

func (p *AppInfo) GetAppName() (v string) {
	return p.AppName
}

func (p *AppInfo) GetGitUrl() (v string) {
	return p.GitUrl
}

func (p *AppInfo) GetBranch() (v string) {
	return p.Branch
}

func (p *AppInfo) GetLanguage() (v string) {
	return p.Language
}

func (p *AppInfo) GetCPU() (v float64) {
	return p.CPU
}

func (p *AppInfo) GetMemory() (v int16) {
	return p.Memory
}

func (p *AppInfo) GetReplicas() (v int8) {
	return p.Replicas
}

func (p *AppInfo) GetMinReplicas() (v int8) {
	return p.MinReplicas
}

func (p *AppInfo) GetMaxReplicas() (v int8) {
	return p.MaxReplicas
}

func (p *AppInfo) GetPortMap() (v []map[string]string) {
	return p.PortMap
}

func (p *AppInfo) GetBuildNo() (v int64) {
	return p.BuildNo
}

func (p *AppInfo) GetAppId() (v int64) {
	return p.AppId
}

func (p *AppInfo) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AppInfo(%+v)", *p)
}

// department、product、env、configBranch、应用名称、branch、git地址、语言、LimitCPU、LimitMem、replicas、map[string]string
type AppDeployReq struct {
	Env          string     `thrift:"env,1" form:"env" json:"env" query:"env"`
	AppInfo      []*AppInfo `thrift:"appInfo,2" form:"appInfo" json:"appInfo" query:"appInfo"`
	Department   string     `thrift:"department,3" form:"department" json:"department" query:"department"`
	Product      string     `thrift:"product,4" form:"product" json:"product" query:"product"`
	ConfigBranch string     `thrift:"configBranch,5" form:"configBranch" json:"configBranch" query:"configBranch"`
}

func NewAppDeployReq() *AppDeployReq {
	return &AppDeployReq{}
}

func (p *AppDeployReq) GetEnv() (v string) {
	return p.Env
}

func (p *AppDeployReq) GetAppInfo() (v []*AppInfo) {
	return p.AppInfo
}

func (p *AppDeployReq) GetDepartment() (v string) {
	return p.Department
}

func (p *AppDeployReq) GetProduct() (v string) {
	return p.Product
}

func (p *AppDeployReq) GetConfigBranch() (v string) {
	return p.ConfigBranch
}

func (p *AppDeployReq) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AppDeployReq(%+v)", *p)
}

type BaseResp struct {
	Code int16             `thrift:"code,1" form:"code" json:"code" query:"code"`
	Msg  string            `thrift:"msg,2" form:"msg" json:"msg" query:"msg"`
	Data map[string]string `thrift:"data,3" form:"data" json:"data" query:"data"`
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

func (p *BaseResp) GetData() (v map[string]string) {
	return p.Data
}

func (p *BaseResp) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("BaseResp(%+v)", *p)
}

type JobReq struct {
	Department string `thrift:"department,1" form:"department" json:"department" query:"department"`
	Product    string `thrift:"product,2" form:"product" json:"product" query:"product"`
	AppName    string `thrift:"appName,3" form:"appName" json:"appName" query:"appName"`
	Env        string `thrift:"env,4" form:"env" json:"env" query:"env"`
	BuildNo    int64  `thrift:"buildNo,5" form:"buildNo" json:"buildNo" query:"buildNo"`
}

func NewJobReq() *JobReq {
	return &JobReq{}
}

func (p *JobReq) GetDepartment() (v string) {
	return p.Department
}

func (p *JobReq) GetProduct() (v string) {
	return p.Product
}

func (p *JobReq) GetAppName() (v string) {
	return p.AppName
}

func (p *JobReq) GetEnv() (v string) {
	return p.Env
}

func (p *JobReq) GetBuildNo() (v int64) {
	return p.BuildNo
}

func (p *JobReq) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("JobReq(%+v)", *p)
}

type JenkinsJobApi interface {
	CreateView(ctx context.Context, baseReq *BaseReq) (r *BaseResp, err error)
	//BaseResp DeleteView(1: BaseReq baseReq)(api.delete="/api/v1/deleteView")
	DeployApp(ctx context.Context, deployReq *AppDeployReq) (r *BaseResp, err error)

	UpdateApp(ctx context.Context, updateReq *AppDeployReq) (r *BaseResp, err error)

	RestartApp(ctx context.Context, restartReq *AppDeployReq) (r *BaseResp, err error)

	RollbackApp(ctx context.Context, rollbackReq *AppDeployReq) (r *BaseResp, err error)
}
