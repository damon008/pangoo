namespace go image


//struct ImageAgent {
//    1: string docker
//    2: string authStr
//}
struct AppInfo {
    1: string app
    2: string tag
    3: optional string gitUrl  //git code 地址
    4: optional string dockerFile //dockerFile路径
    5: optional string confPath //配置中心地址
    6: optional string baseRepo //镜像base推送地址
}

struct ImageInfo {
    1: list<AppInfo> apps
}

struct ImageCommitInfo {
    1: string repo //镜像推送地址
    2: string callBack //服务回调地址
    3: i8  sizeLimit //镜像大小限制(GB)
    4: i8  layerSizeLimit //镜像提交的层的大小限制(保留参数)(GB)
    5: i8  layerLimit //镜像层数限制(保留参数)
    6: i32  imageId   //镜像ID
    7: string podName   //pod名称
}

struct BaseResp {
	1: i8 code
    2: string msg
	3: string data
}

service ImageApi {
    BaseResp ImageBuild(1: ImageInfo imageInfo)(api.post="/api/v1/imageBuild")
    BaseResp ImagePush(1: string repo)(api.post="/api/v1/imagePush")
    BaseResp ImagePull(1: string repo)(api.get="/api/v1/imagePull")
    BaseResp GetImageInfo(1: string repo)(api.get="/api/v1/getImageInfo")
}

