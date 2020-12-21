package truffle

type WarnError struct {
	Code int
	Msg string
	Err error
}

func NewWarnError(code int, msg string)  WarnError{
	w := WarnError{Code: code, Msg: msg}
	return w
}