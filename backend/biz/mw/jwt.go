package mw

import (
	"backend/biz/entity/account"
	"backend/biz/microtype"
	"backend/cmd/dal/model"
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = "identity"
)

const (
	AuthorPrefix = "Bearer "
)

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "test zone",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, utils.H{
				"status_code":    code,
				"token":          token,
				"expire":         expire.Format(time.RFC3339),
				"status_message": "success",
			})
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				AppID       int64  `thrift:"app_id,1" form:"app_id" json:"app_id"`
				AccountName string `thrift:"account_name,2" form:"account_name" json:"account_name"`
				AccountPwd  string `thrift:"account_pwd,3" form:"account_pwd" json:"account_pwd"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			ac := account.NewAccount(0, loginStruct.AppID, loginStruct.AccountName, "", 0, account.NameQuery)
			if err := ac.Load(ctx); err != nil {
				return nil, err
			}

			queryAccount := ac.GetQueryAccount()
			if queryAccount.AccountPwd != loginStruct.AccountPwd {
				return nil, microtype.AccountPwdFailed
			}

			return queryAccount, nil
		},
		IdentityKey: IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &model.Account{
				AccountID: int64(claims[IdentityKey].(float64)),
			}
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.Account); ok {
				return jwt.MapClaims{
					IdentityKey: v.AccountID,
				}
			}
			return jwt.MapClaims{}
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(http.StatusOK, utils.H{
				"status_code": code,
				"status_msg":  message,
			})
			//c.Redirect(consts.StatusMovedPermanently, []byte("/login"))
		},
	})
	if err != nil {
		panic(any(err))
	}
}
