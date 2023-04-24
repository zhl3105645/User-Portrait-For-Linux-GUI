package label

import (
	"backend/biz/microtype"
	"backend/biz/model/backend"
	"context"
)

type PageLabel struct {
	accountId int64
	pageNum   int64
	pageSize  int64
	search    string

	//
	appId int64
	res   []*backend.Label
	total int64
}

func NewPageLabel(accountId int64, pageNum int64, pageSize int64, search string) *PageLabel {
	return &PageLabel{
		accountId: accountId,
		pageSize:  pageSize,
		pageNum:   pageNum,
		search:    search,
	}
}

func (p *PageLabel) Load(ctx context.Context) error {
	//// appId
	//ac := account.NewAccount(p.accountId, 0, "", "", 0, account.IdQuery)
	//if err := ac.Load(ctx); err != nil {
	//	return err
	//}
	//
	//p.appId = ac.GetQueryAccount().AppID
	//
	//// 标签
	//labelDO := query.Label
	//labelMO := labelDO.WithContext(ctx)
	//modelDO := query.DataModel
	//
	//offset := (p.pageNum - 1) * p.pageSize
	//labels, count, err := labelMO.Join(modelDO, modelDO.ModelID.EqCol(labelDO.ModelID)).
	//	Where(modelDO.AppID.Eq(p.appId), labelDO.LabelName.Like("%"+p.search+"%")).
	//	FindByPage(int(offset), int(p.pageSize))
	//if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
	//	return microtype.LabelQueryFailed
	//}
	//
	//if len(labels) == 0 {
	//	return nil
	//}
	//p.total = count
	//
	//// 标签数据
	//labelIds := make([]int64, 0, len(labels))
	//for _, l := range labels {
	//	if l == nil {
	//		continue
	//	}
	//	labelIds = append(labelIds, l.LabelID)
	//}
	//dataDO := query.LabelDatum
	//dataMO := dataDO.WithContext(ctx)
	//labelData, err := dataMO.Where(dataDO.LabelID.In(labelIds...)).Find()
	//if err != nil && !errors.Is(err, gorm.ErrEmptySlice) {
	//	return microtype.LabelDataQueryFailed
	//}
	//// label_id -> data
	//labelDataMap := make(map[int64][]*model.LabelDatum)
	//for _, data := range labelData {
	//	if data == nil {
	//		continue
	//	}
	//	if v, ok := labelDataMap[data.LabelID]; !ok || len(v) == 0 {
	//		labelDataMap[data.LabelID] = make([]*model.LabelDatum, 0)
	//	}
	//	labelDataMap[data.LabelID] = append(labelDataMap[data.LabelID], data)
	//}
	//
	//// 整合数据
	//p.res = make([]*backend.Label, 0, len(labels))
	//for _, label := range labels {
	//	if label == nil {
	//		continue
	//	}
	//
	//	descMap := make(map[int64]string)
	//	err := json.Unmarshal([]byte(label.LabelSemanticDesc), &descMap)
	//	if err != nil {
	//		return microtype.JsonUnMarshalFailed
	//	}
	//
	//	r := &backend.Label{
	//		LabelName: label.LabelName,
	//		LabelID:   label.LabelID,
	//		Option:    chart.GetLabelOption(chart.Label, labelDataMap[label.LabelID], descMap),
	//	}
	//	p.res = append(p.res, r)
	//}

	return nil
}

func (p *PageLabel) GetResp() *backend.LabelInPageResp {
	return &backend.LabelInPageResp{
		StatusCode: microtype.SuccessErr.Code,
		StatusMsg:  microtype.SuccessErr.Msg,
		Labels:     p.res,
		Total:      p.total,
	}
}
