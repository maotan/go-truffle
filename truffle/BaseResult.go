package truffle

type BaseResult struct {
	Code int
	Msg string
	TipMsg string
	Data interface{}
}

func Success(data interface{}) BaseResult {
	base := BaseResult{Code: 0, Data:data}
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