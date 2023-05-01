// Code generated by hertz generator.

package backend

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"backend/biz/mw"
	"backend/biz/usecase/account"
	"backend/biz/usecase/applist"
	"backend/biz/usecase/basic_behavior"
	"backend/biz/usecase/component"
	"backend/biz/usecase/crowd"
	"backend/biz/usecase/data_source"
	"backend/biz/usecase/element"
	"backend/biz/usecase/label"
	model2 "backend/biz/usecase/model"
	"backend/biz/usecase/profile"
	"backend/biz/usecase/register"
	"backend/biz/usecase/rule"
	"backend/biz/usecase/upload"
	"backend/biz/usecase/user"
	"backend/cmd/dal/model"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

// Register .
// @router /register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.RegisterReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.RegisterResp)

	rg := register.NewRegister(req)
	if err := rg.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = rg.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// Login .
// @router /login [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	//var err error
	//var req backend.LoginReq
	//err = c.BindAndValidate(&req)
	//if err != nil {
	//	c.String(consts.StatusBadRequest, err.Error())
	//	return
	//}
	//
	//resp := new(backend.LoginResp)
	//
	//lg := login.NewLogin(req)
	//if err := lg.Load(ctx); err != nil {
	//	mErr := microtype.Unwrap(err)
	//	resp.StatusCode = mErr.Code
	//	resp.StatusMsg = mErr.Msg
	//	c.JSON(consts.StatusOK, resp)
	//	return
	//}
	//
	//resp = lg.GetResp()
	//
	//c.JSON(consts.StatusOK, resp)
}

// AppList .
// @router /applist [GET]
func AppList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AppListReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AppListResp)

	al := applist.NewAppList()
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// Account .
// @router /api/account [GET]
func Account(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AccountReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AccountResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := account.NewAccount(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// AddUser .
// @router /api/user [POST]
func AddUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AddUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AddUserResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := user.NewUser(ac.(*model.Account).AccountID, req)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// UserInPage .
// @router /api/users [GET]
func UserInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.UserInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.UserInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := user.NewPageUser(ac.(*model.Account).AccountID, req.GetPageNum(), req.GetPageSize(), req.GetSearch())
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// UserDataUpload .
// @router /api/user/upload/:id [POST]
func UserDataUpload(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.UserDataUploadReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.UserDataUploadResp)

	userIdStr := c.Param("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "userId=%d", userId)

	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["file"]

	ul := upload.NewFileUpload(userId, files)
	if err := ul.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = ul.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// ComponentInPage .
// @router /api/components [GET]
func ComponentInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.ComponentInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.ComponentInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := component.NewPageComponent(ac.(*model.Account).AccountID, req.GetPageNum(), req.GetPageSize(), req.GetSearch())
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// GeneComponent .
// @router /api/components [POST]
func GeneComponent(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.GeneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.GeneResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := component.NewGeneComponent(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// ElementInPage .
// @router /api/elements [GET]
func ElementInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.ElementInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.ElementInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := element.NewPageElement(req.RuleType, ac.(*model.Account).AccountID, req.PageNum, req.PageSize, req.GetSearch())
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// AddRule .
// @router /api/rule [POST]
func AddRule(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AddRuleReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AddRuleResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := rule.NewAddRule(ac.(*model.Account).AccountID, req)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// UpdateRule .
// @router /api/rule_gene/:id [PUT]
func UpdateRule(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.UpdateRuleReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	ruleIdStr := c.Param("id")
	ruleId, err := strconv.ParseInt(ruleIdStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "ruleId=%d", ruleId)

	resp := new(backend.UpdateRuleResp)

	al := rule.NewUpdateRule(ruleId, req.GetRuleDesc())
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// DeleteRule .
// @router /api/rule_gene/:id [DELETE]
func DeleteRule(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.DeleteRuleReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	ruleIdStr := c.Param("id")
	ruleId, err := strconv.ParseInt(ruleIdStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "ruleId=%d", ruleId)

	resp := new(backend.DeleteRuleResp)

	al := rule.NewDeleteRule(ruleId)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// AddElement .
// @router /api/element [POST]
func AddElement(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AddElementReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AddElementResp)

	al := element.NewAddElement(req)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// UpdateElement .
// @router /api/element/:id [PUT]
func UpdateElement(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.UpdateElementReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.UpdateElementResp)

	elementIdStr := c.Param("id")
	elementId, err := strconv.ParseInt(elementIdStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "elementId=%d", elementId)

	el := element.NewUpdateElement(elementId, req)
	if err := el.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = el.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// DeleteElement .
// @router /api/element/:id [DELETE]
func DeleteElement(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.DeleteElementReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.DeleteElementResp)

	elementIdStr := c.Param("id")
	elementId, err := strconv.ParseInt(elementIdStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "elementId=%d", elementId)

	el := element.NewDeleteElement(elementId)
	if err := el.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = el.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// Rules .
// @router /api/rules [GET]
func Rules(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.RulesReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.RulesResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := rule.NewRules(ac.(*model.Account).AccountID, req.RuleType)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// GeneBasicBehavior .
// @router /api/gene_basic_behavior [POST]
func GeneBasicBehavior(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.GeneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.GeneResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := basic_behavior.NewGeneBasicBehavior(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// BasicBehaviorInPage .
// @router /api/basic_behaviors [GET]
func BasicBehaviorInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.BasicBehaviorInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.BasicBehaviorInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := basic_behavior.NewPageBasicBehavior(ac.(*model.Account).AccountID, req.GetPageNum(), req.GetPageSize(), req.GetSearch())
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// GeneRule .
// @router /api/gene_rule [POST]
func GeneRule(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.GeneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.GeneResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := rule.NewGeneRule(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// RuleDataInPage .
// @router /api/rule_data [GET]
func RuleDataInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.RuleDataInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.RuleDataInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := rule.NewPageRuleData(ac.(*model.Account).AccountID, req.GetPageNum(), req.GetPageSize(), req.GetSearch())
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// DataSources .
// @router /api/data_sources [GET]
func DataSources(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.DataSourceReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.DataSourceResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := data_source.NewDataSources(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// AddModel .
// @router /api/model [POST]
func AddModel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AddModelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AddModelResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := model2.NewAddModel(ac.(*model.Account).AccountID, req)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// ModelInPage .
// @router /api/model [GET]
func ModelInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.ModelInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.ModelInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := model2.NewPageModel(ac.(*model.Account).AccountID, req.PageNum, req.PageSize, req.Search)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// DeleteModel .
// @router /api/model/:id [DELETE]
func DeleteModel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.DeleteModelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	resp := new(backend.DeleteModelResp)

	al := model2.NewDeleteModel(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// GeneModel .
// @router /api/model/:id [POST]
func GeneModel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.GeneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	resp := new(backend.GeneResp)

	al := model2.NewGeneModel(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// AddLabel .
// @router /api/label [POST]
func AddLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AddLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AddLabelResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := label.NewAddLabel(ac.(*model.Account).AccountID, req)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// LabelInPage .
// @router /api/label [GET]
func LabelInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.LabelInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.LabelInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := label.NewPageLabel(ac.(*model.Account).AccountID, req.PageNum, req.PageSize, req.Search)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// DeleteLabel .
// @router /api/label/:id [DELETE]
func DeleteLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.DeleteLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	resp := new(backend.DeleteLabelResp)

	al := label.NewDeleteLabel(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// GeneLabel .
// @router /api/label/:id [POST]
func GeneLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.GeneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	resp := new(backend.GeneResp)

	al := label.NewGeneLabel(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// Users .
// @router /api/all_user [GET]
func Users(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.UsersReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.UsersResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := user.NewUsers(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// DeleteUser .
// @router /api/user/:id [DELETE]
func DeleteUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.DeleteUserReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.DeleteUserResp)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	al := user.NewDeleteUser(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// Labels .
// @router /api/labels [GET]
func Labels(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.LabelsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.LabelsResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := label.NewLabels(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// TreeLabels .
// @router /api/tree_label [GET]
func TreeLabels(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.TreeLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.TreeLabelResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := label.NewTreeLabels(ac.(*model.Account).AccountID, 0, true)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// Profile .
// @router /api/profile/:id [GET]
func Profile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.ProfileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.ProfileResp)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	al := profile.NewProfile(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}
	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// AddCrowd .
// @router /api/crowd [POST]
func AddCrowd(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.AddCrowdReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.AddCrowdResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := crowd.NewAddCrowd(ac.(*model.Account).AccountID, req)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// CrowdInPage .
// @router /api/crowd [GET]
func CrowdInPage(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.CrowdInPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.CrowdInPageResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := crowd.NewPageInCrowd(ac.(*model.Account).AccountID, req.PageNum, req.PageSize, req.Search)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// GeneCrowd .
// @router /api/crowd/:id [POST]
func GeneCrowd(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.GeneReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.GeneResp)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	al := crowd.NewGeneCrowd(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// DeleteCrowd .
// @router /api/crowd/:id [DELETE]
func DeleteCrowd(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.DeleteCrowdReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.DeleteCrowdResp)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	al := crowd.NewDeleteCrowd(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// GroupProfile .
// @router /api/group_profile/:id [GET]
func GroupProfile(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.GroupProfileReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.GroupProfileResp)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	ac, _ := c.Get(mw.IdentityKey)
	al := profile.NewGroupProfile(ac.(*model.Account).AccountID, id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// Crowds .
// @router /api/crowds [GET]
func Crowds(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.CrowdsReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.CrowdsResp)

	ac, _ := c.Get(mw.IdentityKey)
	al := crowd.NewCrowds(ac.(*model.Account).AccountID)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}

// SingleLabel .
// @router /api/label/:id [GET]
func SingleLabel(ctx context.Context, c *app.RequestContext) {
	var err error
	var req backend.SingleLabelReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(backend.SingleLabelResp)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.String(consts.StatusOK, err.Error())
		return
	}
	logger.CtxInfof(ctx, "id=%d", id)

	al := label.NewSingleLabel(id)
	if err := al.Load(ctx); err != nil {
		mErr := microtype.Unwrap(err)
		resp.StatusCode = mErr.Code
		resp.StatusMsg = mErr.Msg
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp = al.GetResp()

	c.JSON(consts.StatusOK, resp)
}
