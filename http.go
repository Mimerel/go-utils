package go_utils

import (
	"crypto/tls"
	"fmt"
	"github.com/op/go-logging"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpRequestParams struct {
	Url string
	Method string
	Body string
	Headers map[string]string
	UserAgent string
	Timeout int64
	Insecure bool
	Retry int
	LogPrefix string
	DelayBetweenRetry time.Duration
	Debug bool
	SuccessCode []int
}


var log = logging.MustGetLogger("default")

func HttpExecuteRequest(requestParams *HttpRequestParams) (err error, response *http.Response) {
	if requestParams.Retry <= 0 {
		requestParams.Retry = 1
	}

	for requestParams.Retry > 0 {
		target, err := url.Parse(requestParams.Url)
		if err != nil {
			log.Errorf(requestParams.LogPrefix + "Unable to parse URL / Url format incorrect %+v", err)
			return err, nil
		}

		if requestParams.Debug {
			log.Debug("Parsed url : %v", target)
		}

		var body io.Reader

		if requestParams.Body != "" {
			body = strings.NewReader(requestParams.Body)
		}

		request, err := http.NewRequest(requestParams.Method, target.String(), body)
		if err != nil {
			log.Errorf(requestParams.LogPrefix + "Unable to create request %+v", err)
			return err, nil
		}

		for key, value := range requestParams.Headers {
			request.Header.Set(key, value)
		}

		if requestParams.UserAgent != "" {
			request.Header.Set("User-Agent", requestParams.UserAgent)
		}

		timeout := time.Duration(60) * time.Second
		if requestParams.Timeout != 0 {
			timeout = time.Duration(requestParams.Timeout) * time.Second
		}

		log.Debugf(requestParams.LogPrefix+"Timeout duration set to : %d seconds", timeout/time.Second)

		request.Close = true

		var httpClient http.Client

		insecure := false
		if requestParams.Insecure {
			insecure = true
		}

		if requestParams.DelayBetweenRetry == 0 {
			requestParams.DelayBetweenRetry = 1
		}
		
		log.Debugf(requestParams.LogPrefix+"Request run in insecure mode ? %t", insecure)
		httpClient = http.Client{Timeout: timeout, Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure}}}

		if requestParams.Debug {
			log.Debug("Request : %v",request)
		}

		response, err = httpClient.Do(request)
		if err != nil || intInArray(response.StatusCode, requestParams.SuccessCode) {
			if requestParams.Debug {
				log.Debug("Response : %v",response)
			}
			requestParams.Retry -= 1
			log.Errorf(requestParams.LogPrefix + "Unable to execute Request %d reties left, %+v", requestParams.Retry, err)
			if requestParams.Retry == 0 {
				return fmt.Errorf("Too Many retires -> failed : %v", err), response
			} else {
				time.Sleep(requestParams.DelayBetweenRetry * time.Second)
			}
		} else {
			requestParams.Retry = 0
			if requestParams.Debug {
				log.Debug("Response : %v",response)
			}
		}
	}
	return nil, response
}

func HttpReadResponse(response *http.Response) (err error, body []byte) {
	log.Debugf("Reading Response")
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Unable to read Response %+v", err)
		return err, nil
	}
	response.Body.Close()
	// log.Debugf("Response : %+v", string(body))
	return nil, body
}

func intInArray( value int, array []int) bool {
	for _, v := range array {
		if v == value {
			return false
		}
	}
	if value > 299 {
		return true
	}
	return false
}