/**
* @Author: mo tan
* @Description:
* @Date 2021/1/9 23:21
 */
package feign

import (
	"github.com/go-resty/resty/v2"
	"github.com/maotan/go-truffle/httpresult"
)

func Head() *httpresult.BaseResult {
	return nil
}

func Delete() *httpresult.BaseResult {
	return nil
}

func Put() *httpresult.BaseResult {
	return nil
}

func Get(appName string, url string) *httpresult.BaseResult {
	res, err := GetRequest(appName).SetResult(&httpresult.BaseResult{}).Get(url)
	if err != nil {
		panic(httpresult.NewWarnError(500, err.Error()))
	}
	base := res.Result().(*httpresult.BaseResult)
	return base
}

func Post(appName string, url string, body interface{}) *httpresult.BaseResult {
	res, err := GetRequest(appName).SetBody(body).SetResult(&httpresult.BaseResult{}).Post(url)
	if err != nil {
		panic(httpresult.NewWarnError(500, err.Error()))
	}
	base := res.Result().(*httpresult.BaseResult)
	return base
}

func GetRequest(appName string) (res *resty.Request) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	return GetRequestWithHeader(appName, header)
}

func GetRequestWithHeader(appName string, header map[string]string) (res *resty.Request) {
	restyClient := DefaultFeign.App(appName)
	if restyClient.HostURL == "" {
		panic(httpresult.NewWarnError(40401, "service not exist"))
	}
	restyReq := restyClient.R().SetHeaders(header)
	return restyReq
}
