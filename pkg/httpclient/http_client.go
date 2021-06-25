package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	tr *http.Transport //
)

const (
	Post   = "POST"
	Delete = "DELETE"
	Get    = "GET"
)

func init() {
	tr = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100}
}

// Header is a type used to store header field name/value pairs when sending HTTPS requests.
type Header map[string]string

// Params is a type used to store parameter name/value pairs when sending HTTPS requests.
type Params map[string]string

// ReqResp hold information relating to an HTTPS request and response.
type ReqResp struct {
	Cli          *http.Client
	Params       Params
	Method       string
	URLAddr      string
	Operation    string
	Timeout      int
	reqParams    string
	Body         interface{}
	ContentType  string
	Resp         *http.Response
	RespText     string
	HeaderFields Header
}

// ReqResp Methods

// CloseBody closes the response body
func (r *ReqResp) CloseBody() {
	if r.Resp != nil {
		if r.Resp.Body != nil {
			e := r.Resp.Body.Close()
			if e != nil {
				log.Warnf("failed to close response body: %s", e.Error())
			}
		}
	}
}

// HTTPreq creates an HTTP client and sends a request. The response is held in ReqResp.RespText
func (r *ReqResp) HTTPreq() error { // nolint gocyclo
	var err error

	if len(r.Method) == 0 {
		r.Method = "GET"
	}
	if len(r.Operation) == 0 {
		return fmt.Errorf("operation is required")
	}
	for k, v := range r.Params {
		if len(v) > 0 {
			r.reqParams += fmt.Sprintf("%s=%s&", k, v)
		}
	}

	httpProto := "http"

	if r.Timeout == 0 {
		r.Timeout = 30
	}
	timeout := time.Duration(int64(time.Second) * int64(r.Timeout))

	r.Cli = &http.Client{Transport: tr}
	r.Cli.Timeout = timeout

	req := fmt.Sprintf("%s://%s/%s", httpProto, r.URLAddr, r.Operation)
	if len(r.reqParams) > 0 {
		req += fmt.Sprintf("?%s", r.reqParams)
	}
	log.Debugf("Request: %s %s", r.Method, req)

	var inputJSON io.ReadCloser
	if r.Method == Post {
		jsonBytes, e := json.Marshal(r.Body)
		if e != nil {
			return fmt.Errorf("failed to convert request body data to JSON, %s", err)
		}
		inputJSON = ioutil.NopCloser(bytes.NewReader(jsonBytes))
	}

	httpReq, err := http.NewRequest(r.Method, req, inputJSON)
	if err != nil {
		return fmt.Errorf("failed to build request, %s", err)
	}

	for k, v := range r.HeaderFields {
		if len(v) > 0 {
			httpReq.Header.Set(k, v)
		}
	}

	retries := 30
	seconds := 1
	start := time.Now()
	for {
		r.Resp, err = r.Cli.Do(httpReq) // nolint bodyclose
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") ||
				strings.Contains(err.Error(), "http2: no cached connection was available") ||
				strings.Contains(err.Error(), "net/http: TLS handshake timeout") ||
				strings.Contains(err.Error(), "i/o timeout") ||
				strings.Contains(err.Error(), "unexpected EOF") ||
				strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
				time.Sleep(time.Second * time.Duration(int64(seconds)))
				retries--
				seconds += seconds
				if seconds > 10 {
					seconds = 2
				}
				if retries > 0 || time.Since(start) > timeout {
					log.Warnf("server failed to respond to %s", r.Operation)
					continue
				}
			}
			return err
		}
		if err := r.GetRespBody(); err != nil {
			return err
		}
		if r.Resp.StatusCode == 200 || (r.Resp.StatusCode == 201 && r.Method == Post) ||
			(r.Resp.StatusCode == 204 && r.Method == Delete) {
			return nil
		}
		return fmt.Errorf("failed: %s, %s", r.Resp.Status, r.RespText)
	}
}

// GetRespBody is used to return the HTTPS response body as a string.
func (r *ReqResp) GetRespBody() error {
	defer r.Resp.Body.Close()
	data, err := ioutil.ReadAll(r.Resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resp: %s", err)
	}
	r.RespText = string(data)
	return nil
}
