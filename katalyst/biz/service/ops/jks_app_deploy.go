/*
 * @Author: Chambers
 * @Date: 2023/6/2 9:42:09
 * @File: pangoo/katalyst/biz/service/ops/jks_app_deploy.go
 * @IDE:  GoLand
 */

package ops

import (
	"reflect"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"

	jk "pangoo/katalyst/biz/model/ops"
)

// JKSDeployMode jenkins 任务部署模式，pipeline根据此模式对应的字段进行发布操作
type JKSDeployMode int

// 定义构建模式对应的枚举值
const (
	JKSDeployModeCI        JKSDeployMode = iota + 1 // 仅构建镜像
	JKSDeployModeCD                                 // 仅部署镜像
	JKSDeployModeCICD                               // 构建镜像并部署镜像
	JKSDeployModeUpdate                             // 更新，包括修改pod limit资源、伸缩pod
	JKSDeployModeRestart                            // 重启deployment
	JKSDeployModeRollback                           // 回滚到上个版本
	JKSDeployModeConfigmap                          // 仅更新配置，发布common
)

// JKSDeployModeStr 定义构建模式枚举值对应的字符串值
var JKSDeployModeStr = []string{"ci", "cd", "cicd", "update", "restart", "rollback", "configmap"}

const (
	defaultTplJob        = "conf/config-tpl.xml"                                         // 定义jenkins模板任务配置文件
	defaultDevOpsOpenApi = "http://10.66.38.120:18081/product/productAppDeploy/callback" // 构建状态需要回调devops平台的接口
)

// JobConfig 流水线构建发布核心结构体，conf 标签用于生成任务配置，param 标签用于生成发布参数
type JobConfig struct {
	BuildHistoryNum uint16 `json:"BuildHistoryNum" conf:"RepBuildHistoryNum"` // 构建历史保留数量

	BuildAppEnv     string `json:"BuildAppEnv" conf:"RepBuildAppEnv"`         // 应用部署环境
	BuildDepartment string `json:"BuildDepartment" conf:"RepBuildDepartment"` // 应用所属产品线
	BuildProduct    string `json:"BuildProduct" conf:"RepBuildProduct"`       // 应用所属项目

	BuildAppID       int64               `json:"BuildAppID"`                                  // 应用ID，上游平台传递的应用id，用于回调传参获取任务成功状态
	BuildAppName     string              `json:"BuildAppName" conf:"RepBuildAppName"`         // 应用名称
	BuildAppGitUrl   string              `json:"BuildAppGitUrl" conf:"RepBuildAppGitUrl"`     // 应用代码仓库地址
	BuildAppLang     string              `json:"BuildAppLang" conf:"RepBuildAppLang"`         // 应用语言：java erlang front
	BuildAppCPU      float64             `json:"BuildAppCPU" conf:"RepBuildAppCPU"`           // 应用 cpu limit 值，单位 C
	BuildAppMemory   int16               `json:"BuildAppMemory" conf:"RepBuildAppMemory"`     // 应用内存 limit值，单位 M
	BuildAppReplicas int8                `json:"BuildAppReplicas" conf:"RepBuildAppReplicas"` // 应用副本数
	BuildAppPortMap  []map[string]string `json:"BuildAppPortMap" conf:"RepBuildAppPort"`      // 应用端口信息

	BuildAppBranch    string        `json:"BuildAppBranch" param:"app_branch"`       // 应用发布分支或标签，默认master
	BuildConfigBranch string        `json:"BuildConfigBranch" param:"config_branch"` // 应用发布配置分支或标签，默认master
	BuildDeployMode   JKSDeployMode `json:"BuildDeployMode" param:"app_deploy_mode"` // 任务部署模式，默认cicd

	BuildPrefix      string            `json:"BuildPrefix"`      // 当前jenkins任务前缀
	BuildJobFullName string            `json:"BuildJobFullName"` // 生成的jenkins 任务全名
	BuildParam       map[string]string `json:"BuildParam"`       // 构建所需参数
	BuildJobRepMap   map[string]string `json:"BuildJobRepMap"`   // 生成替换参数，用于创建更新任务时替换模板
}

// 获取构建模式枚举值对应的字符串值
func (s JKSDeployMode) String() string {
	strArr, iS := JKSDeployModeValues(), int(s)
	if iS < 1 || iS > len(strArr) {
		return ""
	}
	return strArr[iS-1]
}

// Int 获取构建模式枚举值对应的数值
func (s JKSDeployMode) Int() int {
	if gstr.InArray(JKSDeployModeValues(), s.String()) {
		return int(s)
	}
	return 0
}

// JKSDeployModeValues 获取构建模式所有的字符串智
func JKSDeployModeValues() []string {
	return JKSDeployModeStr
}

// JKSDeployModeExistOf 判断字符串值是否属于构建模式对应的字符串值
func JKSDeployModeExistOf(str string) bool {
	for _, v := range JKSDeployModeValues() {
		if v == str {
			return true
		}
	}
	return false
}

// NewJobConfig 创建一个 JobConfig 对象，默认15个构建历史，默认 CICD 构建模式
func NewJobConfig(mode ...JKSDeployMode) *JobConfig {
	thisMode := JKSDeployModeCICD // 默认 CICD 模式
	if len(mode) > 0 {
		thisMode = mode[0]
	}
	return &JobConfig{
		BuildHistoryNum: 15, // 默认15 个构建历史
		BuildDeployMode: thisMode,
	}
}

// ScanAppDeployReq 映射 jk.AppDeployReq 结构体中的公共字段到 JobConfig
func (s *JobConfig) ScanAppDeployReq(in *jk.AppDeployReq) {
	s.BuildAppEnv = in.GetEnv()
	s.BuildDepartment = in.GetDepartment()
	s.BuildProduct = in.GetProduct()
	s.BuildConfigBranch = in.GetConfigBranch()
	s.BuildPrefix = gstr.Join([]string{in.GetDepartment(), in.GetProduct(), in.GetEnv()}, `_`) // 通过下划线 _ 拼接应用前缀信息
}

// ScanAppInfo 映射 jk.AppInfo 结构体中的应用信息字段到 JobConfig
func (s *JobConfig) ScanAppInfo(in *jk.AppInfo, mode ...JKSDeployMode) {
	thisMode := s.BuildDeployMode // 默认 CICD 模式
	if len(mode) > 0 {
		thisMode = mode[0]
	}
	s.BuildAppID = in.GetAppId()
	s.BuildAppName = in.GetAppName()
	s.BuildAppGitUrl = in.GetGitUrl()
	s.BuildAppLang = in.GetLanguage()
	s.BuildAppCPU = in.GetCPU()
	s.BuildAppMemory = in.GetMemory()
	s.BuildAppReplicas = in.GetReplicas()
	s.BuildAppPortMap = in.GetPortMap()
	s.BuildAppBranch = in.GetBranch()
	s.BuildDeployMode = thisMode
}

// GenerateBuildInfo 通过结构体现有信息，生成构建任务名及构建参数
func (s *JobConfig) GenerateBuildInfo() {
	s.BuildJobFullName = s.BuildPrefix + "_" + s.BuildAppName
	s.BuildParam = StructToMapWithTag(*s, "param")    // 使用 param 标签获取构建参数, 需使用字面量而非指针
	s.BuildJobRepMap = StructToMapWithTag(*s, "conf") // 使用 conf 标签获取任务占位符, 需使用字面量而非指针
}

// StructToMapWithTag 根据指定 tag 将结构体转换为 map，没有tag的字段不转换，仅支持字面量结构体，不支持指针
func StructToMapWithTag(obj interface{}, tag string) map[string]string {
	fields, values, result := reflect.TypeOf(obj), reflect.ValueOf(obj), make(map[string]string)
	if values.Kind() != reflect.Struct { // 非结构类型返回 nil
		return nil
	}

	for i := 0; i < fields.NumField(); i++ {
		field, value := fields.Field(i), values.Field(i)

		if tagValue := field.Tag.Get(tag); tagValue != "" { // 获取指定 tag 对应的字段
			kind, val := value.Kind(), ""
			switch {
			case kind == reflect.Uint16:
				val = gconv.String(value.Uint())
			case kind == reflect.String:
				val = value.String()
			case kind == reflect.Float64:
				val = gconv.String(value.Float())
			case kind == reflect.Int8, kind == reflect.Int16:
				val = gconv.String(value.Int())
			case kind == reflect.Int && field.Type != reflect.TypeOf(0): // 判断是基于 int 的新类型
				val = JKSDeployMode(value.Int()).String()
			case field.Type == reflect.TypeOf([]map[string]string{}): // 判断是map切片
				if value.IsNil() { // 判断当前 端口信息 入参为nil，提前退出switch
					break
				}
				if s, err := sonic.MarshalString(value.Interface()); err == nil {
					val = s
				}
			default:
				continue // 不支持类型，跳过
			}
			if len(val) != 0 && val != "0" { // 转换map时排除空值及0值
				result[tagValue] = val
			}
		}
	}
	return result
}

// GetJobTplConfig 获取模板任务的xml文件配置，缓存 60 s
func GetJobTplConfig() string {
	return gfile.GetContentsWithCache(defaultTplJob, time.Second*60)
}

// GetJobActualConfigFromTpl 根据模板任务xml文件配置及具体任务的结构体转map，生成实际任务配置
func GetJobActualConfigFromTpl(jobConfig *JobConfig) string {
	hlog.Debugf("任务模板替换map：%s", jobConfig.BuildJobRepMap)
	return gstr.ReplaceByMap(GetJobTplConfig(), jobConfig.BuildJobRepMap)
}

// ArrMapStrStrFilterEmpty 去除map数组中的空值，包括nil，空切片，空map等，返回存在有效值的map数组，否则返回nil
func ArrMapStrStrFilterEmpty(s []map[string]string) []map[string]string {
	if s == nil || len(s) == 0 { // s 是一个空切片或者未分配内存的切片
		return nil
	}

	res := make([]map[string]string, 0, len(s))
	for _, v := range s {
		if v == nil || len(v) == 0 { // v 是一个未分配内存的 map 或者空 map
			continue
		}
		res = append(res, v)
	}
	return res
}
