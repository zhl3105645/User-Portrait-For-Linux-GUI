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
);

create TABLE user(
    user_id bigint auto_increment comment '用户ID',
    user_name varchar(256) not null comment '用户名',
    app_id bigint not null comment '应用ID',
    primary key (user_id),
    CONSTRAINT a_id2 foreign key (app_id) references app(app_id)
);

create TABLE component(
    component_id bigint auto_increment comment '组件ID',
    component_name text not null comment '组件名',
    component_type int default -1 not null comment '组件类型',
    component_desc text           null comment '组件描述',
    app_id bigint not null comment '应用ID',
    primary key (component_id),
    CONSTRAINT a_id3 foreign key (app_id) references app(app_id)
);

create TABLE rule(
    rule_id bigint auto_increment comment '规则ID',
    rule_type int not null comment '规则类型',
    rule_desc text not null comment '规则描述',
    app_id bigint not null comment '应用ID',
    primary key (rule_id),
    CONSTRAINT a_id4 foreign key (app_id) references app(app_id)
);

create TABLE rule_element(
     rule_element_id bigint auto_increment comment '规则元素ID',
     rule_element_value text not null comment '规则元素值',
     rule_id bigint not null comment '规则ID',
     primary key (rule_element_id),
     CONSTRAINT r_id foreign key (rule_id) references rule(rule_id)
);

create  Table data_source (
    source_id bigint auto_increment comment '数据源ID',
    source_type int not null comment '数据源类型',
    source_value int not null comment '类型的具体值',
    primary key (source_id)
);

create  Table data_model (
      model_id bigint auto_increment comment '数据源ID',
      model_type int not null comment '数据源类型',
      source_id bigint comment '统计数据源ID',
      ml_param text comment '机器学习服务参数',
      model_feature int not null comment '模型用途',
      primary key (model_id),
      CONSTRAINT s_id foreign key (source_id) references data_source(source_id)
);

create  Table model_data (
    model_data_id bigint auto_increment comment '模型数据ID',
    data double not null comment  '模型数据',
    model_id bigint not null comment '模型ID',
    user_id bigint not null comment '用户ID',
    primary key (model_data_id),
    CONSTRAINT m_id foreign key (model_id) references data_model(model_id),
    CONSTRAINT u_id foreign key (user_id) references user(user_id)
);

create TABLE label(
    label_id bigint auto_increment comment '标签ID',
    label_name varchar(256) comment '标签名',
    source_id bigint not null comment '数据源ID',
    label_convert_rule varchar(256) not null comment '标签数据转换规则',
    label_semantic_desc text comment '标签语义化描述',

    primary key (label_id),
    CONSTRAINT s_id2 foreign key (source_id) references data_source(source_id)
);

create Table label_data (
     label_data_id bigint auto_increment comment '标签数据ID',
     data double not null comment  '标签数据',
     label_id bigint not null comment '标签ID',
     user_id bigint not null comment '用户ID',
     primary key (label_data_id),
     CONSTRAINT l_id foreign key (label_id) references label(label_id),
     CONSTRAINT u_id2 foreign key (user_id) references user(user_id)
);