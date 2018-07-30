package vcr

import (
	"fmt"
	"os"
	"testing"

	"github.com/xybstone/baidubce/auth"
)

const (
	DefaultAccessKeyId     = ""
	DefaultSecretAccessKey = ""
	DefaultDebugHost       = "vcr.bj.baidubce.com"
)

var AccessKeyId string
var SecretAccessKey string
var DebugHost string

func TestInit(t *testing.T) {
	AccessKeyId = DefaultAccessKeyId
	if os.Getenv("ACCESS_KEY_ID") != "" {
		AccessKeyId = os.Getenv("ACCESS_KEY_ID")
	}

	SecretAccessKey = DefaultSecretAccessKey
	if os.Getenv("SECRET_ACCESS_KEY") != "" {
		SecretAccessKey = os.Getenv("SECRET_ACCESS_KEY")
	}

	DebugHost = DefaultDebugHost
	if os.Getenv("DEBUG_HOST") != "" {
		DebugHost = os.Getenv("DEBUG_HOST")
	}

}

func TestNewClient(t *testing.T) {
	c, err := NewClient(auth.NewBceCredentials(AccessKeyId, SecretAccessKey))
	if err != nil {
		t.Errorf("NewClient failed.")
	}

	if c.GetEndpoint() != "http://vcr.bj.baidubce.com" {
		t.Errorf("GetEndpoint failed.")
	}

	if c.GetBaseURL() != "http://vcr.bj.baidubce.com/v1" {
		t.Errorf("GetBaseURL failed.")
	}
}

func TestPutText(t *testing.T) {
	c, err := NewClient(auth.NewBceCredentials(AccessKeyId, SecretAccessKey))
	if err != nil {
		t.Errorf("NewClient failed.")
	}
	c.Host = DebugHost

	o, err := c.PuTText("习近平")
	if err != nil {
		t.Errorf("PuTText failed.")
		t.Errorf(err.Error())
	}
	fmt.Println(o)
}
