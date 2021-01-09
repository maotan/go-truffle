/**
* @Author: mo tan
* @Description:
* @Date 2021/1/9 23:21
 */
package feign

import "github.com/go-resty/resty/v2"

func GetRequest(appName string) (res *resty.Request) {
	header := map[string]string{
				"Content-Type": "application/json",
			}
	return GetRequestWithHeader(appName, header)
}

func GetRequestWithHeader(appName string, header map[string]string) (res *resty.Request) {
	restyClient := DefaultFeign.App(appName)
	restyReq := restyClient.R().SetHeaders(header)
	return restyReq
}
