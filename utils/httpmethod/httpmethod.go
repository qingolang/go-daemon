package httpmethod

import (
	"errors"
	"fmt"
	"godaemon/logger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Get
func Get(url string, params map[string]interface{}) ([]byte, error) {

	if len(params) != 0 {
		num := 0
		for i, v := range params {
			if num == 0 {
				url += "?" + i + "=" + fmt.Sprint(v)
			} else {
				url += "&" + i + "=" + fmt.Sprint(v)
			}
			num++

		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Close = true
	// 重试计数
	retryNum := 0

	c := &http.Client{Timeout: 20 * time.Second}
	// 重试TAG
HTTPRetry:

	resp, err := c.Do(request)
	if err != nil {
		if strings.Contains(err.Error(), "net/http: HTTP/1.x transport connection broken: malformed HTTP response") ||
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "EOF") ||
			strings.Contains(err.Error(), "server closed idle connection") ||
			strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
			if retryNum < 5 {

				retryNum++
				time.Sleep(2 * time.Second)
				logger.Log().Warning("HTTP请求中断 ： %s  现在进行重试策略 2s 重试一次 共重试 5次 当前 %d 次", err.Error(), retryNum)
				err = nil
				goto HTTPRetry
			}
		} else {
			logger.Log().Error("HTTP 发生ERR : %s", err.Error())
			err = nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(" Response http.StatusCOde : " + fmt.Sprint(resp.StatusCode))
	}
	return ioutil.ReadAll(resp.Body)

}
