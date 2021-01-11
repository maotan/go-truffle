/**
* @Author: mo tan
* @Description:
* @Date 2021/1/9 18:32
 */
package truffle

type ResultCode struct {
	code  int
	msg	  string
}

var SuccessCode = ResultCode{0, "成功"}
var INNER_ERROR = ResultCode{500, "系统异常"}

// 600 - 699 http异常
var HttpErrorParam = ResultCode{600, "参数错误"}
var HttpServiceNotAvailable = ResultCode{650, "服务不可用"}

// 700 - 799 SQL类异常
var SQLError = ResultCode{700, "SQL操作异常"}

