create table rule
(
    rule_id   bigint auto_increment comment '规则ID'
        primary key,
    rule_type int    not null comment '规则类型',
    rule_desc text   not null comment '规则描述',
    app_id    bigint not null comment '应用ID',
    constraint a_id4
        foreign key (app_id) references app (app_id)
);

INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (1, 1, '新建文件', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (2, 1, '打开文件、工作文件夹', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (3, 1, '保存文件', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (4, 1, '关闭文件', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (5, 1, '导出、打印文件', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (6, 1, '文件属性、编码', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (7, 1, '退出应用', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (8, 1, '撤销、重做', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (9, 1, '代码剪切、复制、粘贴、删除', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (10, 1, '代码跳转', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (11, 1, '代码选择', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (12, 1, '代码移动', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (13, 1, '(取消)缩进', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (14, 1, '代码注释', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (15, 1, '收起/展开', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (16, 1, '切换只读模式', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (17, 1, '视图设置', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (18, 1, '查找按钮', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (19, 1, '替换按钮', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (20, 1, '查找替换窗口', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (21, 1, '添加、删除、修改书签', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (22, 1, '符号重构', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (23, 1, '项目操作(项目添加文件,从项目删除,查看Makefile,清理构建文件,项目属性)', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (24, 1, '关闭项目', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (25, 1, '第三方打开(文件夹中打开,终端中打开)', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (26, 1, '编译', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (27, 1, '编译器版本', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (28, 1, '编译器选项', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (29, 1, '运行', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (30, 1, '运行参数', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (31, 1, '生成汇编', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (32, 1, '调试', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (33, 1, '中断（调试）,单步跨过,单步进入,单步跳出,执行到光标处,继续执行,停止执行', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (34, 1, '添加、删除、修改监视变量', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (35, 1, '(变量)断点', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (36, 1, 'CPU信息', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (37, 1, 'CPU信息窗口', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (38, 1, '选项', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (39, 1, '选项窗口', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (40, 1, '窗口操作', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (41, 1, '帮助', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (42, 1, '点击文件侧边栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (43, 1, '文件工具(当前工作文件夹路径,定位当前文件,隐藏不支持文件)', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (44, 1, '文件树查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (45, 1, '点击项目侧边栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (46, 1, '项目树查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (47, 1, '点击监视侧边栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (48, 1, '监视变量表查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (49, 1, '点击结构侧边栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (50, 1, '结构工具(排序，继承)', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (51, 1, '结构信息查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (52, 1, '点击编译器底部栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (53, 1, '编译问题查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (54, 1, '点击工具输出底部栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (55, 1, '工具输出查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (56, 1, '点击调试底部栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (57, 1, '调试信息查看(主控台、调用栈等)', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (58, 1, '点击查找底部栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (59, 1, '查找信息查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (60, 1, '点击TODO底部栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (61, 1, 'TODO信息查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (62, 1, '点击书签底部栏', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (63, 1, '书签信息查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (64, 1, '代码区输入', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (65, 1, '代码区查看', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (66, 1, 'Git操作', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (67, 2, '编码', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (68, 2, '测试', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (69, 2, '调试', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (70, 2, '浏览', 2);
INSERT INTO profile.rule (rule_id, rule_type, rule_desc, app_id) VALUES (71, 2, '使用工具', 2);
