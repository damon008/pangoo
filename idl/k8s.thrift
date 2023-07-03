namespace go k8s


struct DeployReq {
    1: string appName
    2: string tag     //app的tag
    3: string env     //app部署对应的环境FAT-US
    4: string department
    5: string product
    6: optional i16 port  //启动port
    7: optional string cmd  //启动cmd
}

struct BaseResp {
	1: i16 code
    2: string msg
	3: string data
}

service K8sApi {
    BaseResp CreateDeployment(1: DeployReq deployReq)(api.post="/api/v1/createDeployment/:namespace")
    BaseResp GetDeploymentList(1: string namespace)(api.get="/api/v1/getDeploymentList/:namespace")
    BaseResp GetDeploymentByName(1: string namespace, 2: string name)(api.get="/api/v1/getDeployment/:namespace/:name")
    BaseResp UpdateDeployment(1: DeployReq deployReq)(api.put="/api/v1/updateDeployment/:namespace")
    //重启pod
    BaseResp DeletePod(1: DeployReq deployReq)(api.delete="/api/v1/deletePod")
    BaseResp RollBack(1: DeployReq deployReq)(api.put="/api/v1/rollBack")
}

