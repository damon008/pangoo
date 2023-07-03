namespace go ops

struct JenkinsJobConfig {
    1: string jobName
    2: string env
    3: string configTplName  //任务xml模板名称
    4: string configTplCtx  //任务xml模板内容
    5: string configCtx  //任务xml实际内容
    6: string appName  //任务对应的应用名称
    7: string codeUrl  //任务对应的app代码地址
    8: optional string aliPipelineUrl  //阿里云霄流水线地址
    9: optional string aliPipelineShellName  //阿里云霄流水线执行shellName
    10: optional map<string,string> configTplMap  //任务模板文件替换map
}

struct BaseReq {
    1: string department   //app对应的大模块saas、paas
    2: string product   //app对应的产品业务线
}

struct AppInfo {
    1: string appName
    2: string gitUrl  //git code url
    3: string branch     //app的tag名或分支名
    4: string language //java  front erlang
    5: double CPU
    6: i16 memory //Mi
    7: i8 replicas //副本
    8: i8 minReplicas //最小副本数
    9: i8 maxReplicas //最大副本数
    10: list<map<string, string>> portMap
    11: i64 buildNo
    12: i64 appId // 上游平台传递的应用id，用于回调传参获取任务成功状态
}

//department、product、env、configBranch、应用名称、branch、git地址、语言、LimitCPU、LimitMem、replicas、map[string]string
struct AppDeployReq {
    1: string env //uat-cn
    2: list<AppInfo> appInfo
    3: string department
    4: string product
    5: string configBranch
}

struct BaseResp {
	1: i16 code
    2: string msg
	3: map<string,string> data
}

struct JobReq {
    1: string department   //app对应的大模块saas、paas
    2: string product   //app对应的产品业务线
    3: string appName
    4: string env
    5: i64 buildNo
}

service JenkinsJobApi {
    BaseResp CreateView(1: BaseReq baseReq)(api.post="/api/v1/createView")
    //BaseResp DeleteView(1: BaseReq baseReq)(api.delete="/api/v1/deleteView")
    BaseResp DeployApp(1: AppDeployReq deployReq)(api.post="/api/v1/deployApp")
    BaseResp UpdateApp(1: AppDeployReq updateReq)(api.post="/api/v1/updateApp")
    BaseResp RestartApp(1: AppDeployReq restartReq)(api.post="/api/v1/restartApp")
    BaseResp RollbackApp(1: AppDeployReq rollbackReq)(api.post="/api/v1/rollbackApp")

    //BaseResp UpdateJob(1: AppDeployReq deployReq)(api.put="/api/v1/updateJob")

    //BaseResp DeleteJob(1: JobReq jobReq)(api.delete="/api/v1/deleteJob")
    //BaseResp BuildLog(1: JobReq jobReq)(api.post="/api/v1/buildLog")
}
