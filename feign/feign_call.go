/**
* @Author: mo tan
* @Description:
* @Date 2021/1/9 23:21
 */
package feign

import (
	"github.com/go-resty/resty/v2"
	"github.com/maotan/go-truffle/truffle"
)

func GetRequest(appName string) (res *resty.Request) {
	header := map[string]string{
				"Content-Type": "application/json",
			}
	return GetRequestWithHeader(appName, header)
}

func GetRequestWithHeader(appName string, header map[string]string) (res *resty.Request) {
	restyClient := DefaultFeign.App(appName)
	if restyClient.HostURL == ""{
		panic(truffle.NewWarnError(40401, "service not exist"))
	}
	restyReq := restyClient.R().SetHeaders(header)
	return restyReq
}
