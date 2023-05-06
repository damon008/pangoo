package util

import (
	"bytes"
	"context"

	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var pHTTPClient *HTTPClient
var httpOnce sync.Once

func NewHTTPClient() *HTTPClient {
	httpOnce.Do(func() {
		pHTTPClient = &HTTPClient{Client: &http.Client{Timeout: 10 * time.Second}}
	})
	return pHTTPClient
}

type HTTPClient struct {
	Client *http.Client
}

type HTTPResp struct {
	Status  int
	Data []byte
}

func (hc *HTTPClient) Post(url string, body interface{}) (*HTTPResp, error) {
	bodyStr, _ := sonic.Marshal(body)
	resp, err := hc.Client.Post(url, "application/json", strings.NewReader(string(bodyStr)))
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post fail, status:%s", resp.Status)
	}

	return &HTTPResp{Status: resp.StatusCode, Data: respBody}, nil
}

func (hc *HTTPClient) Get(url string) (*HTTPResp, error) {
	resp, err := hc.Client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get fail, status:%s", resp.Status)
	}

	return &HTTPResp{Status: resp.StatusCode, Data: respBody}, nil
}

func (hc *HTTPClient) GetWithHeader(url string, token string) (*HTTPResp, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", token)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get fail, status:%s", resp.Status)
	}

	return &HTTPResp{Status: resp.StatusCode, Data: respByte}, nil
}

func (hc *HTTPClient) HttpGetReq(url string, params map[string]string, headers map[string]string) (*HTTPResp, error) {
	//new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail ")
	}
	req.Header.Set("Content-type", "application/json")

	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	/*var client = &http.Client{
		Timeout: 10 * time.Second,
	}*/
	log.Printf("Go GET URL : %s \n", req.URL.String())
	resp, err := hc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get fail, status:%s", resp.Status)
	}

	return &HTTPResp{Status: resp.StatusCode, Data: respByte}, nil
}


func (hc *HTTPClient) HttpPostReq(url string, body map[string]string, params map[string]string, headers map[string]string) (*HTTPResp, error) {
	//add post body
	var bodyJson []byte
	var req *http.Request
	if body != nil {
		var err error
		bodyJson, err = sonic.Marshal(body)
		if err != nil {
			log.Println(err)
			return nil, errors.New("http post body to json failed")
		}
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	if err != nil {
		log.Println(err)
		return nil, errors.New("new request is fail: %v \n")
	}
	req.Header.Set("Content-type", "application/json")
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	/*client := &http.Client{
		Timeout: 10 * time.Second,
	}*/
	log.Printf("Go POST URL : %s \n", req.URL.String())
	resp,err := hc.Client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post fail, status:%s", resp.Status)
	}

	return &HTTPResp{Status: resp.StatusCode, Data: respBody}, nil
}

// ?/path/encode uri , method: get/post/delete/put
func (hc *HTTPClient) Do(cli *client.Client, uri string, method string) (*HTTPResp, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	//req.Header.SetHostBytes(req.URI().Host())
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	err := cli.Do(context.Background(), req, res)
	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}
	return &HTTPResp {
		Status: res.StatusCode(),
		Data:   res.Body(),
	},nil
}

//form submit by post/put/update
func (hc *HTTPClient) DoByForm(cli *client.Client, uri string, method string, formData map[string]string) (*HTTPResp, error) {
	req := &protocol.Request{}
	res := &protocol.Response{}
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)
	//req.SetFormData(formData)
	req.SetMultipartFormData(formData)
	//req.Header.SetHostBytes(req.URI().Host())
	req.Options().Apply([]config.RequestOption{config.WithSD(true)})
	err := cli.Do(context.Background(), req, res)
	if err != nil {
		hlog.Error(err.Error())
		return nil, err
	}
	return &HTTPResp {
		Status: res.StatusCode(),
		Data:   res.Body(),
	},nil
}