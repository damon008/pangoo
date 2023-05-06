package singleton

import (
	"context"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/middlewares/client/sd"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/registry/nacos"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"time"
)

var httpCli *client.Client
var once sync.Once

// 方法返回唯一的 Singleton 实例
// 这个 Go 代码中，我们定义了一个名为 Singleton 的结构体作为单例模式的对象，并且定义了一个全局变量 instance 来存储唯一的实例。同时使用 sync.Once 来确保在多线程环境下只有一个实例被创建。
// 客户端可以通过调用 GetInstance 函数来获取唯一的实例。第一次调用时，once.Do 函数会检测 instance 是否为 nil，并且在该变量为 nil 时执行闭包函数，创建一个新的 Singleton 实例并将其赋值给 instance 变量。在后续的调用中，由于 instance 已经不是 nil 了，所以不会再执行闭包函数，直接返回之前创建好的实例即可。
// 这种实现方式具有很好的性能，因为只有在第一次调用 GetInstance 函数时才会进行实例化，之后每次调用都只是简单地返回已有的实例，避免了不必要的对象创建和销毁。同时也保证了线程安全，因为 sync.Once 会保证只有一个 goroutine 能够执行闭包函数来创建实例，其他 goroutine 只能等待这个过程结束后直接返回已有的实例。
func GetHttpCli() *client.Client {
	once.Do(func() {
		c, err := client.NewClient(
			//client.WithDialTimeout(1 * time.Second),//连接建立超时时间，默认 1s
			client.WithWriteTimeout(500*time.Millisecond),       //写入数据超时时间，默认值：无限
			client.WithClientReadTimeout(2*time.Second),         //设置读取 response 的最长时间，默认无限长
			client.WithMaxConnWaitTimeout(100*time.Millisecond), //设置等待空闲连接的最大时间，默认不等待
			client.WithMaxIdleConnDuration(3*time.Second),       //空闲连接超时时间,当超时后会关闭该连接，默认10s
			client.WithMaxConnDuration(3*time.Second),           //设置连接存活的最大时长，超过这个时间的连接在完成当前请求后会被关闭，默认无限长
			client.WithMaxConnsPerHost(1000000),
			//client.WithRetryConfig (
			//	retry.WithMaxAttemptTimes(3), // 最大的尝试次数，包括初始调用
			//	retry.WithInitDelay(1 * time.Millisecond), // 初始延迟
			//	retry.WithMaxDelay(6 * time.Millisecond), // 最大延迟，不管重试多少次，策略如何，都不会超过这个延迟
			//	retry.WithMaxJitter(2 * time.Millisecond), // 延时的最大扰动，结合 RandomDelayPolicy 才会有效果
			//	/*
			//	   配置延迟策略，你可以选择下面四种中的任意组合，最后的结果为每种延迟策略的加和
			//	   FixedDelayPolicy 使用 retry.WithInitDelay 所设置的值 ，
			//	   BackOffDelayPolicy 在 retry.WithInitDelay 所设置的值的基础上随着重试次数的增加，指数倍数增长，
			//	   RandomDelayPolicy 生成 [0，2*time.Millisecond）的随机数值 ，2*time.Millisecond 为 retry.WithMaxJitter 所设置的值，
			//	   DefaultDelayPolicy 生成 0 值，如果单独使用则立刻重试，
			//	   retry.CombineDelay() 将所设置的延迟策略所生成的值加和，最后结果即为当前次重试的延迟时间，
			//	   第一次调用失败 -> 重试延迟：1 + 1<<1 + rand[0,2)ms -> 第二次调用失败 -> 重试延迟：min(1 + 1<<2 + rand[0,2) , 6)ms -> 第三次调用成功/失败
			//	*/
			//	retry.WithDelayPolicy(retry.CombineDelay(retry.FixedDelayPolicy, retry.BackOffDelayPolicy, retry.RandomDelayPolicy)),
			//),
		)
		if err != nil {
			hlog.Error("init client err: ", err)
			c = nil
		}
		httpCli = c
	})
	return httpCli
}

func InitHttpCliWithNacos(conf naming_client.INamingClient) {
	c, err := client.NewClient(
		client.WithDialTimeout(1*time.Second),
		client.WithWriteTimeout(500*time.Millisecond),       //写入数据超时时间，默认值：无限
		client.WithClientReadTimeout(2*time.Second),         //设置读取 response 的最长时间，默认无限长
		client.WithMaxConnWaitTimeout(100*time.Millisecond), //设置等待空闲连接的最大时间，默认不等待
		client.WithMaxIdleConnDuration(3*time.Second),       //空闲连接超时时间
		client.WithMaxConnDuration(3*time.Second),           //设置连接存活的最大时长，超过这个时间的连接在完成当前请求后会被关闭，
		client.WithMaxConnsPerHost(1000000),
		//client.WithRetryConfig (
		//	retry.WithMaxAttemptTimes(3), // 最大的尝试次数，包括初始调用
		//	retry.WithInitDelay(1 * time.Millisecond), // 初始延迟
		//	retry.WithMaxDelay(6 * time.Millisecond), // 最大延迟，不管重试多少次，策略如何，都不会超过这个延迟
		//	retry.WithMaxJitter(2 * time.Millisecond), // 延时的最大扰动，结合 RandomDelayPolicy 才会有效果
		//	/*
		//	   配置延迟策略，你可以选择下面四种中的任意组合，最后的结果为每种延迟策略的加和
		//	   FixedDelayPolicy 使用 retry.WithInitDelay 所设置的值 ，
		//	   BackOffDelayPolicy 在 retry.WithInitDelay 所设置的值的基础上随着重试次数的增加，指数倍数增长，
		//	   RandomDelayPolicy 生成 [0，2*time.Millisecond）的随机数值 ，2*time.Millisecond 为 retry.WithMaxJitter 所设置的值，
		//	   DefaultDelayPolicy 生成 0 值，如果单独使用则立刻重试，
		//	   retry.CombineDelay() 将所设置的延迟策略所生成的值加和，最后结果即为当前次重试的延迟时间，
		//	   第一次调用失败 -> 重试延迟：1 + 1<<1 + rand[0,2)ms -> 第二次调用失败 -> 重试延迟：min(1 + 1<<2 + rand[0,2) , 6)ms -> 第三次调用成功/失败
		//	*/
		//	retry.WithDelayPolicy(retry.CombineDelay(retry.FixedDelayPolicy, retry.BackOffDelayPolicy, retry.RandomDelayPolicy)),
		//),
	)
	if err != nil {
		hlog.Error("init client err: ", err)
		c = nil
	}
	c.Use(sd.Discovery(nacos.NewNacosResolver(conf)))
	httpCli = c
}

func InitHttpCli() {
	c, err := client.NewClient(
		//client.WithDialTimeout(1 * time.Second),//连接建立超时时间，默认 1s
		client.WithWriteTimeout(500*time.Millisecond),       //写入数据超时时间，默认值：无限
		client.WithClientReadTimeout(2*time.Second),         //设置读取 response 的最长时间，默认无限长
		client.WithMaxConnWaitTimeout(100*time.Millisecond), //设置等待空闲连接的最大时间，默认不等待
		client.WithMaxIdleConnDuration(3*time.Second),       //空闲连接超时时间,当超时后会关闭该连接，默认10s
		client.WithMaxConnDuration(3*time.Second),           //设置连接存活的最大时长，超过这个时间的连接在完成当前请求后会被关闭，默认无限长
		client.WithMaxConnsPerHost(1000000),
		//client.WithRetryConfig (
		//	retry.WithMaxAttemptTimes(3), // 最大的尝试次数，包括初始调用
		//	retry.WithInitDelay(1 * time.Millisecond), // 初始延迟
		//	retry.WithMaxDelay(6 * time.Millisecond), // 最大延迟，不管重试多少次，策略如何，都不会超过这个延迟
		//	retry.WithMaxJitter(2 * time.Millisecond), // 延时的最大扰动，结合 RandomDelayPolicy 才会有效果
		//	/*
		//	   配置延迟策略，你可以选择下面四种中的任意组合，最后的结果为每种延迟策略的加和
		//	   FixedDelayPolicy 使用 retry.WithInitDelay 所设置的值 ，
		//	   BackOffDelayPolicy 在 retry.WithInitDelay 所设置的值的基础上随着重试次数的增加，指数倍数增长，
		//	   RandomDelayPolicy 生成 [0，2*time.Millisecond）的随机数值 ，2*time.Millisecond 为 retry.WithMaxJitter 所设置的值，
		//	   DefaultDelayPolicy 生成 0 值，如果单独使用则立刻重试，
		//	   retry.CombineDelay() 将所设置的延迟策略所生成的值加和，最后结果即为当前次重试的延迟时间，
		//	   第一次调用失败 -> 重试延迟：1 + 1<<1 + rand[0,2)ms -> 第二次调用失败 -> 重试延迟：min(1 + 1<<2 + rand[0,2) , 6)ms -> 第三次调用成功/失败
		//	*/
		//	retry.WithDelayPolicy(retry.CombineDelay(retry.FixedDelayPolicy, retry.BackOffDelayPolicy, retry.RandomDelayPolicy)),
		//),
	)
	if err != nil {
		hlog.Error("init client err: ", err)
		c = nil
	}
	httpCli = c
}

func Get(uri string, query string) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodGet)
	req.SetRequestURI(uri)
	req.SetQueryString(query) //a=1&b=2
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}

// Use SetQueryString to set query parameters
func Post(uri string, query string) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.SetRequestURI(uri)
	req.SetQueryString(query) //a=1&b=2
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}

// Send "www-url-encoded" request
func FormData4Post(uri string, body map[string]string) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.SetRequestURI(uri)
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	req.SetMultipartFormData(body)
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}

// Send "multipart/form-data" request
func MultipartFormData4Post(uri string, body map[string]string) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.SetRequestURI(uri)
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	req.SetMultipartFormData(body)
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}

// body = "{"key1":"bbbb"}"
func JsonStr4Post(uri string, body string) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI(uri)
	//req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	//jsonByte, _ := sonic.Marshal(body)
	req.SetBody([]byte(body))
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}

// Send body "Json" request
func Json4Post(uri string, body interface{}) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI(uri)
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	jsonByte, _ := sonic.Marshal(body)
	req.SetBody(jsonByte)
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}

// ?/path/encode uri , method: get/post/delete/put
func HttpDo(uri string, method string) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	//req.Header.SetHostBytes(req.URI().Host())
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}

// form submit by post/put/update
func DoByForm(uri string, method string, formData map[string]string) (*protocol.Response, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	//req.SetFormData(formData)
	req.SetMultipartFormData(formData)
	//req.Header.SetHostBytes(req.URI().Host())
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	//err := httpCli.Do(context.Background(), req, res)
	if err := httpCli.Do(context.Background(), req, res); err != nil {
		return res, err
	}
	return res, nil
}
