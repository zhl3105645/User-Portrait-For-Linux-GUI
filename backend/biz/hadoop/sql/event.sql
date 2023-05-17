create database profile;

create table event(
      record_id bigint,
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

use profile;

insert into table event values
(2, 2, 121212, "(12,23)", 1,1,1,1,'\\','com_name', 1, '[错误] ' 'd' ' undeclared (first use in this function)', 2, '2022-01-02'),
(3, 3, 121212, "(12,23)", 1,1,1,1,'A','com_name', 1, 'd', 2, '2022-01-02');

select * from event;

select  count(distinct record_id) from event;

select count(*) from event;

truncate table event;