####################### sql操作
# 选取构建聚类相关数据
select
    u.user_id as user_id,
    u.behavior_duration_map as behavior_duration_map
from user as u
order by user_id;


delete from label_data where label_data.user_id >= 13 and label_data.user_id <= 51;
delete from record where record.user_id >= 13 and record.user_id <= 51;
delete from user where user.user_id >= 13 and user.user_id <= 51;