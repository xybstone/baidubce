package vcr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xybstone/baidubce/auth"
	"github.com/xybstone/baidubce/httplib"
	"github.com/xybstone/baidubce/utils"
)

const (
	DefaultLocation   = "vcr.bj"
	DefaultAPIVersion = "v1"
)

type VcrClient struct {
	Credential *auth.BceCredentials
	Location   string
	APIVersion string
	Host       string
	Debug      bool
}

func NewClient(credential *auth.BceCredentials) (VcrClient, error) {
	return VcrClient{
		Credential: credential,
		Location:   DefaultLocation,
		APIVersion: DefaultAPIVersion,
		Debug:      false,
	}, nil
}

func (c VcrClient) GetBaseURL() string {
	return fmt.Sprintf("%s/%s", c.GetEndpoint(), c.APIVersion)
}

func (c VcrClient) GetEndpoint() string {
	return fmt.Sprintf("http://%s", c.GetHost())
}

func (c VcrClient) GetHost() string {
	if c.Host != "" {
		return c.Host
	}
	return fmt.Sprintf("%s.baidubce.com", c.Location)
}

type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("Service returned error: Code=%s, RequestId=%s, Message=%s", e.Code, e.RequestId, e.Message)
}

func (c VcrClient) doRequest(req *httplib.Request) (*http.Response, error) {
	if req.BaseUrl == "" {
		req.BaseUrl = c.GetBaseURL()
	}
	req.Headers[httplib.HOST] = c.GetHost()

	timestamp := utils.GetHttpHeadTimeStamp()
	auth.Debug = c.Debug
	authorization := auth.Sign(c.Credential, timestamp, req.Method, req.Path, req.Query, req.Headers)

	req.Headers[httplib.BCE_DATE] = timestamp
	req.Headers[httplib.AUTHORIZATION] = authorization

	httplib.Debug = c.Debug
	res, err := httplib.Run(req, nil)
	if err != nil {
		return res, err
	}

	if res.StatusCode != 200 && res.StatusCode != 206 {
		errR := &ErrorResponse{}
		if req.Method == httplib.HEAD || req.Method == httplib.DELETE {
			errR.Code = fmt.Sprintf("%d", res.StatusCode)
			errR.Message = res.Status
			errR.RequestId = "EMPTY"
		} else {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return res, err
			}
			j := json.NewDecoder(strings.NewReader(string(body)))
			j.Decode(&errR)
		}
		return res, errR
	}
	return res, err
}

/*************************************************************************************************

VCR Opreation Method

*************************************************************************************************/

/*
 * Name: PuTText
 * URL: /v1/text
 */

type PutTextResponse struct {
	Label   string      `json:"label"`
	Results interface{} `json:"results"`
}

func (c VcrClient) PuTText(text string) (output PutTextResponse, err error) {
	ts := map[string]string{"text": text}
	b, _ := json.Marshal(&ts)
	timestamp := utils.GetHttpHeadTimeStamp()

	req := &httplib.Request{
		Method: httplib.PUT,
		Headers: map[string]string{
			"Host":         c.GetHost(),
			"Content-Type": "application/json",
		},
		Body: bytes.NewReader(b),
		Path: c.APIVersion + "/text",
	}

	auth.Debug = true
	authorization := auth.Sign(c.Credential, timestamp, req.Method, req.Path, req.Query, req.Headers)
	req.Headers[httplib.AUTHORIZATION] = authorization

	res, err := c.doRequest(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	j := json.NewDecoder(strings.NewReader(string(body)))
	j.Decode(&output)
	return
}
