/**
* @Author: mo tan
* @Description:
* @Date 2021/1/1 21:58
 */
package httpresult

type WarnError struct {
	Code int
	Msg  string
	Err  error
}

func NewWarnError(code int, msg string) WarnError {
	w := WarnError{Code: code, Msg: msg}
	return w
}
