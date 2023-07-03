namespace go analyse

#include "base.thrift"
struct AnalyseReq {
    1: i16 projectId(api.path = 'projectId')
}

struct BaseResp {
	1: i8 code
    2: string msg
	3: string data
}

struct Ns {
    1: i32 id
    2: string name
    3: string path
    4: string kind
    5: string full_path
    6: i32 parent_id
    7: string avatar_url
    8: string web_url
}
//建表关联
struct ProjectInfo {
    1: i64 id
    2: string department
    3: string product
    4: string appName
    5: i16 appId
    6: string branch
    7: string ref
    8: string tag
    9: string refTag
}

struct Project {
    1: i16 id
    2: string description
    3: string name
    4: string http_url_to_repo
    5: string ssh_url_to_repo
    6: string web_url
}

struct Projects {
    1: i16 id
    2: string description
    3: string name
    4: string name_with_namespace
    5: string path
    6: string path_with_namespace
    7: string created_at
    8: string default_branch
    9: list<string> tag_list
    10: list<string> topics
    11: string ssh_url_to_repo
    12: string http_url_to_repo
    13: string web_url
    14: string avatar_url
    15: i32 star_count
    16: string last_activity_at
    17: Ns ns
}

service AnalyseApi {
    BaseResp CreateBranch(1: ProjectInfo projectInfo)(api.post="/api/v1/createBranch")
    BaseResp CreateTag(1: ProjectInfo projectInfo)(api.post="/api/v1/createTag")
    BaseResp GetAll()(api.get="/api/v1/getProjects")
    BaseResp CloneCode(1: list<string> gitUrls)(api.post="/api/v1/cloneCode")
    #BaseResp GetProject(1: i16 projectId)(api.get="/api/v1/getProject/:projectId")
    BaseResp CompareInfo(1: ProjectInfo projectInfo)(api.post="/api/v1/compareInfo")
    //BaseResp MergeReq(1: string projectName, 2: string branchName, 3: string targetBranch)(api.post="/api/v1/mergeReq/:projectName/:branchName/:targetBranch")
}
#hz new -module katalyst --idl=../idl/hello.thrift -t=template=slim --thrift-plugins=validator
#hz update --idl=../idl/analyse.thrift -t=template=slim