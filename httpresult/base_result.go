/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package httpresult

type BaseResult struct {
	Code    int
	Msg     string
	TipMsg  string
	TraceId string
	Data    interface{}
}

func Success(data interface{}) BaseResult {
	base := BaseResult{Code: 0, Data: data}
	return base
}

func SuccessWithTip(data interface{}, tipMsg string) BaseResult {
	base := BaseResult{Code: 0, Data: data, TipMsg: tipMsg}
	return base
}

func Fail(code int, msg string) BaseResult {
	base := BaseResult{Code: code, Msg: msg}
	return base
}

func FailWithTip(code int, msg string, tipMsg string) BaseResult {
	base := BaseResult{Code: code, Msg: msg, TipMsg: tipMsg}
	return base
}

const (
	// 成功编码
	SuccessCode = 0
	// 参数错误编码
	ParamErrCode = 400
	// 系统异常编码
	SystemErrCode = 500
)
