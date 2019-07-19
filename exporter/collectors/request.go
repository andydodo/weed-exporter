package collectors

import (
	"crypto/tls"
	"errors"
	. "github.com/andy/logger"
	"io/ioutil"
	"net/http"
	"time"
)

type Req struct {
	Url     string
	Method  string
	Timeout time.Duration
	Retry   int
	// add max worker nums
	MaxNums int
}

type ReqApi interface {
	Request() (*http.Response, error)
}

func NewRequest() *Req {
	return &Req{}
}

func (r *Req) Request() ([]byte, error) {
	tr := &http.Transport{
		MaxIdleConnsPerHost: r.MaxNums,
		MaxIdleConns:        20 * r.MaxNums,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Timeout: r.Timeout, Transport: tr}

	req, err := http.NewRequest(r.Method, r.Url, nil)

	if err != nil {
		Logger.Printf("http request url error:(%v)", err)
		return nil, errors.New("http request error")
	}

	resp, err := httpClient.Do(req)

	if err != nil || err != nil && resp.StatusCode != 200 {
		Logger.Printf("http response url error:(%v)", err)
		return nil, errors.New("http response error")
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}

	return nil, errors.New(string(resp.StatusCode))
}
