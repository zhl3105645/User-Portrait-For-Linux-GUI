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

	UnknownErr          = &Error{Code: 10000000, Msg: "未知错误"}
	ParamCheckFailed    = &Error{Code: 10000001, Msg: "参数错误"}
	JsonMarshalFailed   = &Error{Code: 10000002, Msg: "json 编码失败"}
	JsonUnMarshalFailed = &Error{Code: 10000003, Msg: "json 解码失败"}

	AccountExist       = &Error{Code: 10001001, Msg: "账户已存在"}
	AccountAddFailed   = &Error{Code: 10001002, Msg: "添加账户失败"}
	AccountNotExist    = &Error{Code: 10001003, Msg: "账户不存在"}
	AccountQueryFailed = &Error{Code: 10001004, Msg: "账户查询失败"}
	AccountPwdFailed   = &Error{Code: 10001005, Msg: "账户密码错误"}

	AppExist            = &Error{Code: 10002001, Msg: "应用已存在"}
	AppParamCheckFailed = &Error{Code: 10002002, Msg: "添加应用参数错误"}
	AppAddFailed        = &Error{Code: 10002003, Msg: "添加应用失败"}
	AppFindAllFailed    = &Error{Code: 10002004, Msg: "查询全部应用失败"}

	UserQueryFailed  = &Error{Code: 10003001, Msg: "用户查询失败"}
	UserExist        = &Error{Code: 10003002, Msg: "用户已存在"}
	UserAddFailed    = &Error{Code: 10003003, Msg: "添加用户失败"}
	UserNameEmpty    = &Error{Code: 10003004, Msg: "用户名为空"}
	UserDeleteFailed = &Error{Code: 10003005, Msg: "用户删除失败"}

	DirOpenFailed = &Error{Code: 10004001, Msg: "打开目录失败"}
	DirReadFailed = &Error{Code: 10004002, Msg: "读取目录文件失败"}

	ComponentQueryFailed  = &Error{Code: 10005001, Msg: "组件查询失败"}
	ComponentCreateFailed = &Error{Code: 10005002, Msg: "组件插入失败"}
	ComponentInGene       = &Error{Code: 10005003, Msg: "组件信息生成中"}

	RuleQueryFailed  = &Error{Code: 10006001, Msg: "规则查询失败"}
	RuleParamFailed  = &Error{Code: 10006002, Msg: "规则参数错误"}
	RuleExist        = &Error{Code: 10006003, Msg: "规则已存在"}
	RuleCreateFailed = &Error{Code: 10006004, Msg: "规则创建失败"}
	RuleUpdateFailed = &Error{Code: 10006005, Msg: "规则更新失败"}
	RuleDeleteFailed = &Error{Code: 10006005, Msg: "规则删除失败"}

	ElementQueryFailed  = &Error{Code: 10007001, Msg: "元素查询失败"}
	ElementCreateFailed = &Error{Code: 10007002, Msg: "元素添加失败"}
	ElementUpdateFailed = &Error{Code: 10007003, Msg: "元素更新失败"}
	ElementDeleteFailed = &Error{Code: 10007004, Msg: "元素删除失败"}

	BasicBehaviorQueryFailed = &Error{Code: 10008001, Msg: "基础行为数据查询失败"}
	BasicBehaviorGene        = &Error{Code: 10008002, Msg: "基础行为数据生成中"}

	RuleGene            = &Error{Code: 10009001, Msg: "规则数据生成中"}
	RuleDataQueryFailed = &Error{Code: 10009002, Msg: "规则数据查询失败"}

	DataSourceCreateFailed = &Error{Code: 10010001, Msg: "数据源创建失败"}
	DataSourceDeleteFailed = &Error{Code: 10010002, Msg: "数据源删除失败"}
	DataSourceQueryFailed  = &Error{Code: 10010003, Msg: "数据源查询失败"}

	DataModelCreateFailed = &Error{Code: 10011001, Msg: "模型创建失败"}
	DataModelQueryFailed  = &Error{Code: 10011002, Msg: "模型查询失败"}
	DataModelDeleteFailed = &Error{Code: 10011003, Msg: "模型删除失败"}

	ModelDataGeneFailed   = &Error{Code: 10012001, Msg: "模型数据生成失败"}
	ModelDataQueryFailed  = &Error{Code: 10012002, Msg: "模型数据查询失败"}
	ModelDataDeleteFailed = &Error{Code: 10012002, Msg: "模型数据删除失败"}

	LabelCreateFailed = &Error{Code: 10013001, Msg: "标签创建失败"}
	LabelQueryFailed  = &Error{Code: 10013002, Msg: "标签查询失败"}
	LabelDeleteFailed = &Error{Code: 10013003, Msg: "标签删除失败"}
	LabelUpdateFailed = &Error{Code: 10013004, Msg: "标签更新失败"}

	LabelDataGeneFailed   = &Error{Code: 10014001, Msg: "标签数据生成失败"}
	LabelDataQueryFailed  = &Error{Code: 10014002, Msg: "标签数据查询失败"}
	LabelDataDeleteFailed = &Error{Code: 10014002, Msg: "标签数据删除失败"}
)

func Unwrap(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}

	return UnknownErr
}
