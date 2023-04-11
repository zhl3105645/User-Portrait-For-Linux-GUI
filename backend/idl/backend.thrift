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

service BackendService {
    // 未登录状态
    RegisterResp Register(1: RegisterReq request) (api.post="/register");
    LoginResp Login(1: LoginReq request) (api.post="/login");
    AppListResp AppList(1: AppListReq request) (api.get="/applist");

    // 登录状态 api 开头
    AccountResp Account(1: AccountReq request) (api.get="/api/account");
}