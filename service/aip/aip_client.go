package aip

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/xybstone/baidubce/auth"
	"github.com/xybstone/baidubce/httplib"
	"github.com/xybstone/baidubce/utils"
)

const (
	DefaultLocation   = "aip"
	DefaultAPIVersion = "/rest/2.0/solution/v1"
)

type AipClient struct {
	Credential *auth.BceCredentials
	Location   string
	APIVersion string
	Host       string
	Debug      bool
}

func NewClient(credential *auth.BceCredentials) (AipClient, error) {
	return AipClient{
		Credential: credential,
		Location:   DefaultLocation,
		APIVersion: DefaultAPIVersion,
		Debug:      false,
	}, nil
}

func (c AipClient) GetBaseURL() string {
	return fmt.Sprintf("%s/%s", c.GetEndpoint(), c.APIVersion)
}

func (c AipClient) GetEndpoint() string {
	return fmt.Sprintf("http://%s", c.GetHost())
}

func (c AipClient) GetHost() string {
	if c.Host != "" {
		return c.Host
	}
	return fmt.Sprintf("%s.baidubce.com", c.Location)
}

func (c AipClient) doRequest(req *httplib.Request) ([]byte, error) {
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
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

/* ConclusionType取值 1:合规, 2:不合规, 3:疑似, 4:审核失败 */
const (
	ConclusionTypeQualified   = 1
	ConclusionTypeUnqualified = 2
	ConclusionTypeUnsure      = 3
	ConclusionTypeFailure     = 4
)

type ImageAuditResponse struct {
	ErrorCode      int `json:"error_code"`
	ConclusionType int `json:"conclusionType"`
}

/* 图像审核 */
func (c AipClient) ImageAudit(r io.Reader) (*ImageAuditResponse, error) {
	imageData, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("image", base64.StdEncoding.EncodeToString(imageData))
	timestamp := utils.GetHttpHeadTimeStamp()

	req := &httplib.Request{
		Method: httplib.POST,
		Headers: map[string]string{
			"Host": c.GetHost(),
		},
		Type: "application/x-www-form-urlencoded",
		Body: bytes.NewReader([]byte(v.Encode())),
		Path: c.APIVersion + "/img_censor/user_defined",
	}

	auth.Debug = false
	authorization := auth.Sign(c.Credential, timestamp, req.Method, req.Path, req.Query, req.Headers)
	req.Headers[httplib.AUTHORIZATION] = authorization

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var response ImageAuditResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

