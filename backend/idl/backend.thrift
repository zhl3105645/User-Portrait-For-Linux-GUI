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

struct GeneComponentReq {

}

struct GeneComponentResp {
    1: i64 status_code
    2: string status_msg
}

struct ElementInPageReq {
    1: i64 page_num
    2: i64 page_size
    3: string search
    4: i64 rule_type
}

struct Rule {
    1: i64 rule_id
    2: i64 rule_type
    3: string rule_desc
    4: list<EventRuleElement> event_rule_elements
    5: list<BehaviorRuleElement> behavior_rule_elements
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
    5: list<i64> event_rule_ids
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
    GeneComponentResp GeneComponent(1: GeneComponentReq request) (api.post="/api/components");
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
    UpdateElementResp UpdateElement(1: UpdateElementReq request) (api.put="/api/element/:id");
    DeleteElementResp DeleteElement(1: DeleteElementReq request) (api.delete="/api/element/:id");
}