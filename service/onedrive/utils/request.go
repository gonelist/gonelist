package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"gonelist/service/onedrive/auth"
)

// GetData
/**
 * @Description: 请求微软的api
 * @param method 请求方法 GET,POST,DELETE,PUT
 * @param url1 url
 * @param headers 请求头
 * @param data 请求体
 * @return []byte 响应内容
 * @return error
 */
func GetData(method, url1 string, headers map[string]string, data []byte) ([]byte, error) {
	var (
		resp *http.Response
		body []byte
		err  error
	)

	client := auth.GetClient()
	if client == nil {
		logrus.Errorln("cannot get client to start request.")
		return nil, fmt.Errorf("RequestOneURL cannot get client")
	}
	request, err := http.NewRequest(method, url1, bytes.NewReader(data))
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	// 如果超时，重试两次
	for retryCount := 3; retryCount > 0; retryCount-- {
		if resp, err = client.Do(request); err != nil && strings.Contains(err.Error(), "timeout") {
			logrus.WithFields(logrus.Fields{
				"url":  url1,
				"resp": resp,
				"err":  err,
			}).Info("RequestOneUrl 出现错误，开始重试")
			<-time.After(time.Second / 3)
		} else {
			break
		}
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":  url1,
			"resp": resp,
			"err":  err,
		}).Info("请求 graph.microsoft.com 失败, request timeout")
		return body, err
	}

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		logrus.WithField("err", err).Info("读取 graph.microsoft.com 返回内容失败")
		return body, err
	}
	return body, nil
}
