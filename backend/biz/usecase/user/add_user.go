package user

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/cmd/dal/model"
	"backend/cmd/dal/query"
	"context"
	"github.com/golang/protobuf/proto"
	"strconv"
)

type AddUser struct {
	accountId int64
	req       backend.AddUserReq
	//
	appId      int64
	userName   string
	userGender Gender
	userAge    int64
	userCareer string
}

func NewUser(accountId int64, req backend.AddUserReq) *AddUser {
	return &AddUser{
		accountId: accountId,
		req:       req,
	}
}

func (u *AddUser) Load(ctx context.Context) error {
	if !u.check() {
		return microtype.ParamCheckFailed
	}
	ac := account.NewAccount(u.accountId, 0, "", "", 0, account.IdQuery)
	if err := ac.Load(ctx); err != nil {
		return err
	}
	u.appId = ac.GetQueryAccount().AppID

	// 插入
	userDO := query.User
	userMO := userDO.WithContext(ctx)

	mo := &model.User{
		UserID:     0,
		UserName:   u.userName,
		AppID:      u.appId,
		UserGender: int64(u.userGender),
		UserAge:    u.userAge,
		UserCareer: proto.String(u.userCareer),
	}

	err := userMO.Create(mo)
	if err != nil {
		return microtype.UserAddFailed
	}

	return nil
}

func (u *AddUser) GetResp() *backend.AddUserResp {
	return &backend.AddUserResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
	}
}

func (u *AddUser) check() bool {
	if u.req.GetUsername() == "" || u.req.GetUserGender() == "" || u.req.GetUserAge() == "" {
		return false
	}
	u.userName = u.req.GetUsername()

	gender, err := strconv.ParseInt(u.req.GetUserGender(), 10, 64)
	if err != nil {
		return false
	}
	u.userGender = Gender(gender)

	age, err := strconv.ParseInt(u.req.GetUserAge(), 10, 64)
	if err != nil {
		return false
	}
	u.userAge = age
	u.userCareer = u.req.UserCareer

	return true
}
