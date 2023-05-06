package util

import (
	"context"
	"crypto/tls"
	//"github.com/cloudwego/hertz/pkg/common/config"

	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/retry"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/hertz-contrib/http2/config"
	"github.com/hertz-contrib/http2/factory"
	"net/http"
	"time"
)

const (
	keyPEM  = `<your key PEM>`
	certPEM = `<your cert PEM>`
)

var h2Cli *client.Client

func Http2Cli() {
	cli, _ := client.NewClient()
	cli.SetClientFactory(factory.NewClientFactory(
		config.WithDialTimeout(3*time.Second),
		config.WithReadIdleTimeout(1*time.Second),
		config.WithWriteByteTimeout(time.Second),
		config.WithPingTimeout(time.Minute),
		config.WithMaxIdleConnDuration(2*time.Second),
		config.WithClientDisableKeepAlive(true),     //Close Connection after each request
		config.WithStrictMaxConcurrentStreams(true), // Set the server's SETTINGS_MAX_CONCURRENT_STREAMS to be respected globally.
		config.WithDialer(standard.NewDialer()),     // You can customize dialer here
		config.WithMaxHeaderListSize(0xffffffff),    // Set SETTINGS_MAX_HEADER_LIST_SIZE to unlimited.
		config.WithMaxIdempotentCallAttempts(3),
		config.WithRetryConfig(
			retry.WithMaxAttemptTimes(3),
			retry.WithInitDelay(2*time.Millisecond),
			retry.WithMaxDelay(200*time.Millisecond),
			retry.WithMaxJitter(30*time.Millisecond),
			retry.WithDelayPolicy(retry.FixedDelayPolicy),
		),
		config.WithStrictMaxConcurrentStreams(true), // Set the server's SETTINGS_MAX_CONCURRENT_STREAMS to be respected globally.
		config.WithTLSConfig(&tls.Config{
			SessionTicketsDisabled: false,
			InsecureSkipVerify:     true,
		}),
	))
	h2Cli = cli
}

func TestHttp2Cli() {
	v, _ := sonic.Marshal(map[string]string{
		"hello":    "world",
		"protocol": "h2",
	})

	for {
		time.Sleep(time.Second * 1)
		req, rsp := protocol.AcquireRequest(), protocol.AcquireResponse()
		req.SetMethod("POST")
		req.SetRequestURI("https://127.0.0.1:8888")
		req.SetBody(v)
		err := h2Cli.Do(context.Background(), req, rsp)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[client]: received body: %s\n", string(rsp.Body()))
	}
}

func NewH2Serv() {
	cfg := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}

	cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
	if err != nil {
		fmt.Println(err.Error())
	}
	cfg.Certificates = append(cfg.Certificates, cert)
	h := server.New(server.WithHostPorts(":8888"), server.WithALPN(true), server.WithTLS(cfg))

	// register http2 server factory
	h.AddProtocol("h2", factory.NewServerFactory(
		config.WithReadTimeout(time.Minute), //建立连接后，从服务器读取到可用资源的超时时间
		config.WithDisableKeepAlive(false))) //是否关闭 Keep-Alive模式

	cfg.NextProtos = append(cfg.NextProtos, "h2")

	h.POST("/", func(c context.Context, ctx *app.RequestContext) {
		var j map[string]string
		_ = sonic.Unmarshal(ctx.Request.Body(), &j)
		fmt.Printf("[server]: received request: %+v\n", j)

		r := map[string]string{
			"msg": "hello world",
		}
		for k, v := range j {
			r[k] = v
		}
		ctx.JSON(http.StatusOK, r)
	})

	go Http2Cli()

	h.Spin()
}
