drop Table app;

create TABLE app(
    app_id bigint auto_increment comment '应用ID',
    app_name varchar(256) unique comment '应用名',
    primary key (app_id)
);

drop TABLE account;

create TABLE account(
    account_id bigint auto_increment comment '账号ID',
    account_name varchar(256) not null comment '账号名',
    account_pwd varchar(256) not null comment '账号密码',
    account_permission int default 1 not null comment '账号权限',
    app_id bigint not null comment '应用名',
    primary key (account_id),
    CONSTRAINT a_id foreign key (app_id) references app(app_id)
)