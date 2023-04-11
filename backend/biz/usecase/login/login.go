package login

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"context"
)

type Login struct {
	req backend.LoginReq
}

func NewLogin(req backend.LoginReq) *Login {
	return &Login{
		req: req,
	}
}

func (l *Login) Load(ctx context.Context) error {

	return nil
}

func (l *Login) GetResp() *backend.LoginResp {
	resp := &backend.LoginResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Token:      "",
	}

	return resp
}
