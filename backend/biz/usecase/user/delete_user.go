package user

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/query"
	"context"
)

type DeleteUser struct {
	userId int64
}

func NewDeleteUser(userId int64) *DeleteUser {
	return &DeleteUser{
		userId: userId,
	}
}

func (d *DeleteUser) Load(ctx context.Context) error {

	// 删除用户
	_, err := query.User.WithContext(ctx).Where(query.User.UserID.Eq(d.userId)).Delete()
	if err != nil {
		return microtype.UserQueryFailed
	}

	return nil
}

func (d *DeleteUser) GetResp() *backend.DeleteUserResp {
	return &backend.DeleteUserResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}
