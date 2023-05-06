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

package mw

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
	"primus/consumer-service/model/auth"
	"primus/pkg/constants"
	"primus/pkg/res"
	"primus/pkg/util/singleton"
)

/*
* 以MW形式通过第三方鉴权来进行判断是否已认证
* @author Damon
* @date   2023/5/5 13:46
 */

type TokenMiddleware struct {

}

//func (mw *TokenMiddleware) AuthMiddleWare() []app.HandlerFunc {
func AuthMiddleWare() app.HandlerFunc {
	//Auth鉴权逻辑可以先进行判断常量：如果不为err，则表示登录成功过，可以直接访问接口，如果为err，则需要重新登录
	/*if constants.AuthState !=nil {

	}*/
	if (constants.AuthState) {
		return func(ctx context.Context, c *app.RequestContext) {
			hlog.Debug("auth success")
			//如果已经授权，不需要返回，直接走目标路由逻辑
			/*c.JSON(http.StatusOK, auth.Response{
				Msg:  &auth.Msg{
					Code: 0,
					Msg:  "authorized",
				},
				Data: "",
			})*/
			//return
		}
	}
	return newHandlerFunc()
}

func newHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		token := c.GetHeader("Authorization")
		hlog.Info("header: ", token)
		name := c.Query("name")//?name="damon"
		hlog.Info("name: ", name)
		//c.Param("id")///:id
		//resp, err := util.NewHTTPClient().Get("http://localhost:2999/v1/checkAuth")
		resp, err := singleton.Get("http://localhost:2999/v1/checkAuth", "")
		if err != nil {
			hlog.Error("request the auth check handler fail: ", err)
			c.Header("Auth", "unauthorized")
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, auth.Response{
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
		if resp.StatusCode() != 200 || err != nil || licenceContent.Msg.Code != 0 {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("Auth", "unauthorized")
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, auth.Response{
				Msg:  &auth.Msg{
					Code: -1,
					Msg:  "unauthorized",
				},
				Data: "",
			})
			return
		} else {
			constants.AuthState = true
			/*c.JSON(http.StatusOK, auth.Response{
				Msg:  &auth.Msg{
					Code: 0,
					Msg:  "authorized",
				},
				Data: "",
			})*/
		}
	}
}
