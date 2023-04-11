package microtype

import "errors"

type Error struct {
	Code int64
	Msg  string
}

func (e *Error) Error() string {
	return e.Msg
}

var (
	SuccessErr = &Error{Code: 0, Msg: ""}

	UnknownErr = &Error{Code: 10000000, Msg: "未知错误"}

	AccountExist       = &Error{Code: 10001001, Msg: "账户已存在"}
	AccountAddFailed   = &Error{Code: 10001002, Msg: "添加账户失败"}
	AccountNotExist    = &Error{Code: 10001003, Msg: "账户不存在"}
	AccountQueryFailed = &Error{Code: 10001004, Msg: "账户查询失败"}
	AccountPwdFailed   = &Error{Code: 10001005, Msg: "账户密码错误"}

	AppExist            = &Error{Code: 10002001, Msg: "应用已存在"}
	AppParamCheckFailed = &Error{Code: 10002002, Msg: "添加应用参数错误"}
	AppAddFailed        = &Error{Code: 10002003, Msg: "添加应用失败"}
	AppFindAllFailed    = &Error{Code: 10002004, Msg: "查询全部应用失败"}
)

func Unwrap(err error) *Error {
	var e *Error
	if errors.As(err, e) {
		return e
	}

	return UnknownErr
}
