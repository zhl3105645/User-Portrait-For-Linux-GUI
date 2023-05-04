create database profile;

show current roles ;

set role admin;

use profile;

set system:user.name;

desc  test4;


create table test4(id int, name string);
insert into table test4 values (3, "2");
select * from test4;

drop table test4;

show tables;


create table test_partition(id int, name string)
    partitioned by (day string);

select * from profile.test_partition;

insert into table profile.test_partition values (2, "2", "2023-01-02");

create table event(
                      user_id bigint, -- 用户ID
                      begin_time bigint, -- 开始时间
                      event_type int,
                      event_time bigint,
                      mouse_pos string,
                      mouse_click_type int,
                      mouse_click_btn int,
                      mouse_move_type int,
                      key_click_type int,
                      key_code string,
                      component_name string,
                      component_type int,
                      component_extra string
)
    partitioned by (app_id int, day string); // 多级分区 应用ID 天

insert into table event values (1, 12121212, 1, 121212, "(12,23)", 1,1,1,1,'A','com_name', 1, 'd', 2, '2022-01-02');

select * from event;