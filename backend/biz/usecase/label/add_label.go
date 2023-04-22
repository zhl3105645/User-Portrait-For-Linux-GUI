package label

//
//import (
//	"backend/biz/entity/account"
//	"backend/biz/microtype"
//	"backend/biz/model/backend"
//	"backend/cmd/dal/model"
//	"backend/cmd/dal/query"
//	"context"
//	"encoding/json"
//	"strconv"
//)
//
//type AddLabel struct {
//	accountId int64
//	req       backend.AddLabelReq
//
//	//
//	appId int64
//}
//
//func NewAddLabel(accountId int64, req backend.AddLabelReq) *AddLabel {
//	return &AddLabel{
//		accountId: accountId,
//		req:       req,
//	}
//}
//
//func (a *AddLabel) Load(ctx context.Context) error {
//	if !a.check() {
//		return microtype.ParamCheckFailed
//	}
//
//	ac := account.NewAccount(a.accountId, 0, "", "", 0, account.IdQuery)
//	if err := ac.Load(ctx); err != nil {
//		return err
//	}
//
//	a.appId = ac.GetQueryAccount().AppID
//
//	// 转换规则 && 语义描述
//	converts := make([]*ConvertRule, 0, len(a.req.ConvertRules))
//	descMap := make(map[int64]string) // label_value -> desc
//	for _, rule := range a.req.ConvertRules {
//		if rule == nil {
//			continue
//		}
//		x, _ := strconv.ParseFloat(rule.XValue, 64)
//		y, _ := strconv.ParseInt(rule.YValue, 10, 64)
//		convert := &ConvertRule{
//			Op:     getOperator(rule.Operator),
//			XValue: x,
//			YValue: y,
//		}
//
//		converts = append(converts, convert)
//		descMap[y] = rule.YDesc
//	}
//
//	convertStr, err := json.Marshal(converts)
//	if err != nil {
//		return microtype.JsonMarshalFailed
//	}
//
//	descStr, err := json.Marshal(descMap)
//	if err != nil {
//		return microtype.JsonMarshalFailed
//	}
//
//	// 添加
//	mo := &model.Label{
//		LabelID:           0,
//		LabelName:         a.req.LabelName,
//		ModelID:           a.req.ModelID,
//		LabelConvertRule:  string(convertStr),
//		LabelSemanticDesc: string(descStr),
//	}
//
//	err = query.Label.WithContext(ctx).Create(mo)
//	if err != nil {
//		return microtype.LabelCreateFailed
//	}
//
//	return nil
//}
//
//func (a *AddLabel) GetResp() *backend.AddLabelResp {
//	return &backend.AddLabelResp{
//		StatusCode: microtype.SuccessErr.Code,
//		StatusMsg:  microtype.SuccessErr.Msg,
//	}
//}
//
//func (a *AddLabel) check() bool {
//	if a.req.LabelName == "" || a.req.ModelID <= 0 || len(a.req.ConvertRules) < 0 {
//		return false
//	}
//	for _, rule := range a.req.ConvertRules {
//		if rule == nil {
//			continue
//		}
//		if getOperator(rule.Operator) == Unknown {
//			return false
//		}
//		if rule.YDesc == "" {
//			return false
//		}
//		_, err := strconv.ParseFloat(rule.XValue, 64)
//		if err != nil {
//			return false
//		}
//		_, err = strconv.ParseInt(rule.YValue, 10, 64)
//		if err != nil {
//			return false
//		}
//	}
//
//	return true
//}
