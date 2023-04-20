namespace go backend

struct RegisterReq {
    1: string app_name (api.body="app_name"); //
    2: string account_name (api.body="account_name");
    3: string account_pwd (api.body="account_pwd")
}

struct RegisterResp {
    1: i64 status_code
    2: string status_msg
}

struct LoginReq {
    1: i64 app_id (api.body="app_id");
    2: string account_name (api.body="account_name");
    3: string account_pwd (api.body="account_pwd")
}

struct LoginResp {
    1: i64 status_code
    2: string status_msg
    3: string token
}

struct AppListReq {

}

struct App {
    1: i64 app_id
    2: string app_name
}

struct AppListResp {
    1: i64 status_code
    2: string status_msg
    3: list<App> apps
}

struct AccountReq {

}

struct Account {
    1: string account_name
    2: i64 account_permission
}

struct AccountResp {
    1: i64 status_code
    2: string status_msg
    3: Account account
}

struct AddUserReq {
    1: string username (api.body="username");
}

struct AddUserResp {
    1: i64 status_code
    2: string status_msg
}

struct UserInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
}

struct User {
    1: i64 user_id
    2: string user_name
    3: i64 record_num
}

struct UserInPageResp {
    1: i64 status_code
    2: string status_msg
    3: list<User> users
    4: i64 total // 用户总数
}

struct UserDataUploadReq {
    1: i64 id (api.param="id")
}

struct UserDataUploadResp {
    1: i64 status_code
    2: string status_msg
}

struct ComponentInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
}

struct Component {
    1: i64 component_id
    2: string component_name
    3: i64 component_type
    4: string component_desc
}

struct ComponentInPageResp {
    1: i64 status_code
    2: string status_msg
    3: list<Component> components
    4: i64 total // 组件总数
}

struct GeneReq {

}

struct GeneResp {
    1: i64 status_code
    2: string status_msg
}

struct ElementInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
    4: i64 rule_type
}

struct EventRuleElement {
    1: i64 element_id
    2: i64 rule_id
    3: i64 rule_type
    4: string rule_desc
    5: i64 event_type
    6: i64 mouse_click_type
    7: i64 mouse_click_button
    8: i64 key_click_type
    9: string key_value
    10: string component_name_prefix
}

struct BehaviorRuleElement {
    1: i64 element_id
    2: i64 rule_id
    3: i64 rule_type
    4: string rule_desc
    5: list<RuleElement> event_rules
}

struct RuleElement {
    1: i64 rule_id
    2: string rule_desc
    3: i64 timestamp
}

struct ElementInPageResp {
    1: i64 status_code
    2: string status_msg
    3: list<EventRuleElement> event_elements
    4: list<BehaviorRuleElement> behavior_elements
    5: i64 total
}

struct AddRuleReq {
    1: i64 rule_type (api.body="rule_type")
    2: string rule_desc (api.body="rule_desc")
}

struct AddRuleResp {
    1: i64 status_code
    2: string status_msg
}

struct UpdateRuleReq {
    1: i64 rule_type (api.body="rule_type")
    2: string rule_desc (api.body="rule_desc")
}

struct UpdateRuleResp {
    1: i64 status_code
    2: string status_msg
}

struct DeleteRuleReq {

}

struct DeleteRuleResp {
    1: i64 status_code
    2: string status_msg
}

struct AddElementReq {
    1: i64 rule_id (api.body="rule_id")
    2: i64 event_type (api.body="event_type")
    3: i64 mouse_click_type (api.body="mouse_click_type")
    4: i64 mouse_click_button (api.body="mouse_click_button")
    5: i64 key_click_type (api.body="key_click_type")
    6: string key_value (api.body="key_value")
    7: string component_name_prefix (api.body="component_name_prefix")
    8: list<i64> event_rule_ids (api.body="event_rule_ids")
}

struct AddElementResp {
    1: i64 status_code
    2: string status_msg
}

struct UpdateElementReq {
    1: i64 event_type (api.body="event_type")
    2: i64 mouse_click_type (api.body="mouse_click_type")
    3: i64 mouse_click_button (api.body="mouse_click_button")
    4: i64 key_click_type (api.body="key_click_type")
    5: string key_value (api.body="key_value")
    6: string component_name_prefix (api.body="component_name_prefix")
    7: list<i64> event_rule_ids (api.body="event_rule_ids")
}

struct UpdateElementResp {
    1: i64 status_code
    2: string status_msg
}

struct DeleteElementReq {

}

struct DeleteElementResp {
    1: i64 status_code
    2: string status_msg
}

struct RulesReq {
    1: i64 rule_type
}

struct RulesResp {
    1: i64 status_code
    2: string status_msg
    3: list<RuleElement> event_rules
}

struct BasicBehaviorInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
}

struct BasicBehaviorInPageResp {
    1: i64 status_code
    2: string status_msg
    3: list<BasicBehavior> basic_behaviors
    4: i64 total
}

struct BasicBehavior {
    1: i64 record_id
    2: i64 user_id
    3: string user_name
    4: string begin_time
    5: string use_time
    6: i64 mouse_click_cnt
    7: i64 mouse_move_cnt
    8: double mouse_move_dis
    9: i64 mouse_wheel_cnt
    10: i64 key_click_cnt
    11: double key_click_speed
    12: i64 shortcut_cnt
}

struct RuleDataInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
}

struct EventRuleData {
    1: list<RuleElement> rule_elements
}

struct BehaviorRuleData {
    1: list<RuleElement> rule_elements
    2: map<string, i64> behavior_duration // 行为时长
}

struct RuleData {
    1: i64 record_id
    2: i64 user_id
    3: string user_name
    4: string begin_time
    5: EventRuleData event_rule_data
    6: BehaviorRuleData behavior_rule_data
}

struct RuleDataInPageResp {
    1: i64 status_code
    2: string status_msg
    3: list<RuleData> rule_data
    4: i64 total
}

struct DataSourceReq {

}

struct DataSource {
    1: i64 value // 一级：数据源类型 二级：数据源值
    2: string label // 数据源含义
    3: list<DataSource> children
}

struct DataSourceResp {
    1: i64 status_code
    2: string status_msg
    3: list<DataSource> data_source
}

struct LearningParam {
    1: i64 http_type (api.body="http_type")// http类型 默认 post
    2: list<BodyParam> body_params (api.body="body_params")// body参数
    3: string http_resp_name (api.body="http_resp_name")
    4: i64 http_resp_data_type (api.body="http_resp_data_type")
    5: string http_addr (api.body="http_addr")
}

struct BodyParam {
    1: string name (api.body="name")// 参数名
    2: i64 model_id (api.body="model_id")// 对应的模型ID
}

struct AddModelReq {
    1: string model_name (api.body="model_name") // 模型名
    2: string model_type (api.body="model_type")// 模型类型：统计，机器学习
    3: i64 calculate_type (api.body="calculate_type")// 统计计算方式：平均数、众数等
    4: i64 source_type (api.body="source_type")// 数据源类型
    5: i64 source_value (api.body="source_value")// 数据源ID
    6: string model_feature (api.body="model_feature") // 模型功能 默认为 label
    7: LearningParam learning_param (api.body="learning_param") // 机器学习服务参数
}

struct AddModelResp {
    1: i64 status_code
    2: string status_msg
}

struct ModelInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
}

struct Axis {
    1: string type
    2: list<string> data
    3: string name
    4: string position
    5: AxisLabel axisLabel
}

struct AxisLabel {
    1: i64 rotate
    2: bool show
}

struct Series {
    1: string type
    2: bool smooth // 是否光滑
    3: list<string> data
    4: i64 yAxisIndex
}

struct ToolTip {
    1: string trigger
    2: AxisPointer axisPointer
    3: string formatter
}

struct AxisPointer {
    1: string type
}

struct ToolBox {
    1: Feature feature
}

struct View {
    1: bool show
}

struct Feature {
    1: View dataView
    2: View saveAsImage
}

struct ChartOption {
    1: list<Axis> xAxis
    2: list<Axis> yAxis
    3: list<Series> series
    4: ToolTip tooltip
    5: ToolBox toolbox
}

struct Model {
    1: string model_name
    2: i64 model_id
    3: i64 model_type
    4: ChartOption option
}

struct ModelInPageResp {
    1: i64 status_code
    2: string status_msg
    3: list<Model> models
    4: i64 total
}

struct DeleteModelReq {

}

struct DeleteModelResp {
    1: i64 status_code
    2: string status_msg
}

struct ConvertRule {
    1: string operator
    2: string x_value
    3: string y_value
    4: string y_desc
}

struct AddLabelReq {
    1: string label_name
    2: i64 model_id // 模型ID
    3: list<ConvertRule> convert_rules // 规则
}

struct AddLabelResp {
    1: i64 status_code
    2: string status_msg
}

struct LabelInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
}

struct Label {
    1: string label_name
    2: i64 label_id
    3: ChartOption option
}

struct LabelInPageResp {
    1: i64 status_code
    2: string status_msg
    3: list<Label> labels
    4: i64 total
}

struct DeleteLabelReq {

}

struct DeleteLabelResp {
    1: i64 status_code
    2: string status_msg
}

service BackendService {
    // 未登录状态
    // 注册
    RegisterResp Register(1: RegisterReq request) (api.post="/register");
    // 登录
    LoginResp Login(1: LoginReq request) (api.post="/login");
    // 应用列表
    AppListResp AppList(1: AppListReq request) (api.get="/applist");

    // 登录状态 api 开头
    // 获取账号信息
    AccountResp Account(1: AccountReq request) (api.get="/api/account");
    // 添加用户
    AddUserResp AddUser(1: AddUserReq request) (api.post="/api/user");
    // 获取用户信息 分页
    UserInPageResp UserInPage(1: UserInPageReq request) (api.get="/api/users");
    // 导入某一用户数据 上传文件
    UserDataUploadResp UserDataUpload(1: UserDataUploadReq request) (api.post="/api/user/upload/:id");
    // 组件信息 分页
    ComponentInPageResp ComponentInPage(1: ComponentInPageReq request) (api.get="/api/components");
    // 生成全部组件信息
    GeneResp GeneComponent(1: GeneReq request) (api.post="/api/components");
    // 规则元素信息 分页
    ElementInPageResp ElementInPage(1: ElementInPageReq request) (api.get="/api/elements");
    // 添加规则
    AddRuleResp AddRule(1: AddRuleReq request) (api.post="/api/rule");
    // 更新规则
    UpdateRuleResp UpdateRule(1: UpdateRuleReq request) (api.put="/api/rule/:id");
    // 删除规则
    DeleteRuleResp DeleteRule(1: DeleteRuleReq request) (api.delete="/api/rule/:id");
    // 添加规则元素
    AddElementResp AddElement(1: AddElementReq request) (api.post="/api/element");
    // 更新规则元素
    UpdateElementResp UpdateElement(1: UpdateElementReq request) (api.put="/api/element/:id");
    // 删除规则元素
    DeleteElementResp DeleteElement(1: DeleteElementReq request) (api.delete="/api/element/:id");
    // 获取规则
    RulesResp Rules(1: RulesReq request) (api.get="/api/rules");
    // 生成基础行为数据
    GeneResp GeneBasicBehavior(1: GeneReq request) (api.post="/api/gene_basic_behavior");
    // 基础行为数据 分页
    BasicBehaviorInPageResp BasicBehaviorInPage(1: BasicBehaviorInPageReq request) (api.get="/api/basic_behaviors");
    // 生成规则数据
    GeneResp GeneRule(1: GeneReq request) (api.post="/api/gene_rule");
    // 规则数据 分页
    RuleDataInPageResp RuleDataInPage(1: RuleDataInPageReq request) (api.get="/api/rule_data");
    // 数据源
    DataSourceResp DataSources(1: DataSourceReq request) (api.get="/api/data_sources");
    // 添加模型
    AddModelResp AddModel(1: AddModelReq request) (api.post="/api/model");
    // 模型 分页
    ModelInPageResp ModelInPage(1: ModelInPageReq request) (api.get="/api/model");
    // 删除模型
    DeleteModelResp DeleteModel(1: DeleteModelReq request) (api.delete="/api/model/:id");
    // 生成模型数据
    GeneResp GeneModel(1: GeneReq request) (api.post="/api/model/:id");
    // 添加标签
    AddLabelResp AddLabel(1: AddLabelReq request) (api.post="/api/label");
    // 标签 分页
    LabelInPageResp LabelInPage(1: LabelInPageReq request) (api.get="/api/label");
    // 删除标签
    DeleteLabelResp DeleteLabel(1: DeleteLabelReq request) (api.delete="/api/label/:id");
    // 生成模型数据
    GeneResp GeneLabel(1: GeneReq request) (api.post="/api/label/:id");
}