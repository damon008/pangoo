package license

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"primus/consumer-service/model/auth"
	"primus/pkg/constants"
	"primus/pkg/errno"
	"primus/pkg/res"
	"primus/pkg/util/singleton"
)

type IssuedContent struct {
	AppID         string   `json:"appId"`                                                                                     //证书签发来源标识
	IssuedTime    int64    `json:"issuedTime" `                                                                               //证书签发时间
	NotBefore     int64    `json:"notBefore"`                                                                                 //证书有效起始时间
	NotAfter      int64    `json:"notAfter"`                                                                                  //证书有效截止时间
	CustomerInfo  string   `json:"customerInfo"`                                                                              //证书适用方标识
	Authorization string   `json:"authorization" validate:"required,oneof=all training inference reosurce-manager log-proxy"` //授权范围,可选：all,training,inference,reosurce-manager,log-proxy
	MachineCodes  []string `json:"machineCodes"`                                                                              //机器码，⽤于机器签发限定
	NodeNum       int      `json:"nodeNum"`                                                                                   //节点数量，数量为0，代表不限制数量
}

func LicenceChecker() {
	hlog.Debug("软件授权系统自检")
	//var licenseFile = "E:\\study\\primus\\pkg\\license\\LICENSE"
	var licenseFile = "/Users/damon/program/study/go/primus/pkg/license/LICENSE"
	byteLicence, err := VerifyLicence(licenseFile)
	if err != nil {
		hlog.Error("verify the licence error: ", err)
		constants.LicenceState = errno.NewErrNo(-1, "false")
		//constants.LicenceState = errno.ConvertErr(err)
		return
	}
	licenceContent := IssuedContent{}
	err = sonic.Unmarshal(byteLicence, &licenceContent)
	if err != nil {
		hlog.Error("unmarshal the licenceContent error", err)
		constants.LicenceState = errno.NewErrNo(40303, "unmarshal the licenceContent error")
		return
	}
	// 获取当前时间
	now := GetTimeNow().Unix()
	if now < licenceContent.IssuedTime {
		hlog.Error("check the licence fail,the system time incorrect,now is %d,issuedTime is %d", now, licenceContent.IssuedTime)
		constants.LicenceState = errno.NewErrNo(40304, "check the licence fail,the system time incorrect")
		return
	}
	if now < licenceContent.NotBefore || now > licenceContent.NotAfter {
		hlog.Error("check the licence fail,certificate expired,now is %d,not before is %d , not after is %d", now, licenceContent.NotBefore, licenceContent.NotAfter)
		constants.LicenceState = errno.NewErrNo(40305, "check the licence fail,certificate expired")
		return
	}
	// 集群规模校验
	hlog.Debug("the content is ", licenceContent)
	/*if licenceContent.NodeNum != 0 { //node number 等于0 代表对节点数不做限制
		nodeNum, err := GetClusterNodes()
		if err != nil {
			hlog.Error("get the node num fail", err)
			constants.LicenceState = errno.NewErrNo(40307, "check the licence fail,get the node num fail")
			return
		}
		if nodeNum > licenceContent.NodeNum {
			hlog.Error("the the number of nodes is incorret,the node num is %d,the licence nodeNum is %d", nodeNum, licenceContent.NodeNum)
			constants.LicenceState = errno.NewErrNo(40308, "check the licence fail,the number of nodes is incorret")
		}
	}*/

	// 机器码校验 机器码的个数小于0，代表对机器码不做限制
	if len(licenceContent.MachineCodes) > 0 {
		machineCode, err := GetSystemSerialNum()
		if err != nil {
			hlog.Error("get the machineCode fail", err)
			constants.LicenceState = errno.NewErrNo(40311, "check the licence fail,get the machineCode fail")
			return
		}
		machineHashCode := md5.Sum([]byte(machineCode))
		machineHashCodeStr := fmt.Sprintf("%x", machineHashCode)
		if !containsInMap(convertStrSlice2Map(licenceContent.MachineCodes), machineHashCodeStr) {
			hlog.Error("check the machineCode fail,the machineCode is %s", machineCode)
			constants.LicenceState = errno.NewErrNo(40310, "check the licence fail,check the machineCode fail")
			return
		}
	}
}

func LicenceIssued() []app.HandlerFunc {
	return []app.HandlerFunc{
		func(ctx context.Context, c *app.RequestContext) {
			//resp, err := util.NewHTTPClient().Get("http://localhost:2999/v1/checkLicence")
			resp, err := singleton.Get("http://localhost:10000/v1/checkLicence", "")
			if err != nil {
				hlog.Error("request the licence check handler fail: ", err)
				c.Header("Copyright-Software", "unauthorized")
				c.AbortWithStatus(http.StatusBadGateway)
				c.JSON(http.StatusBadGateway, auth.Response{
					Msg:  &auth.Msg{
						Code: -1,
						Msg:  "unauthorized",
					},
					Data: "",
				})
				return
			}
			hlog.Info(string(resp.Body()))
			licenceContent := res.DefaultResponse{}
			err = sonic.Unmarshal(resp.Body(), &licenceContent)
			//拿到授权状态再进行判断
			//hlog.Info(licenceContent.Msg.Msg)
			//hlog.Info(licenceContent.Msg.Code)
			//hlog.Info(licenceContent.Data)
			if resp.StatusCode() != 200 || err != nil || licenceContent.Msg.Code != 0 {
				// Credentials doesn't match, we return 401 and abort handlers chain.
				c.Header("Copyright-Software", "unauthorized")
				c.AbortWithStatus(http.StatusBadGateway)
				c.JSON(http.StatusBadGateway, auth.Response{
					Msg:  &auth.Msg{
						Code: -1,
						Msg:  "unauthorized",
					},
					Data: "",
				})
				return
			} else {
				hlog.Debug(licenceContent)
			}
		},
	}
}

func convertStrSlice2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// ContainsInMap 判断字符串是否在 map 中
func containsInMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}

//func respLicence(ctx *context.Context, httpStatus int, msg string) {
//	ctx.Output.Status = httpStatus
//	resp := base_api.DefaultResponse{Status: *customerror.NewError(customerror.ComErrorCodeLicenceCheckFail, msg)}
//	ctx.Output.JSON(resp, false, false)
//}
