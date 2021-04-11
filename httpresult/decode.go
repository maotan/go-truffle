package httpresult

import (
	"github.com/goinggo/mapstructure"
)

func Decode(base *BaseResult, resultObject interface{}) {
	if err := mapstructure.Decode(base.Data, resultObject); err != nil {
		panic(NewWarnError(DecodeErrCode, err.Error()))
	}
}
