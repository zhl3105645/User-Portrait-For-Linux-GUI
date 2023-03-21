import os
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
    

def single_running(begin_time, end_time):
    single_running_time.append(end_time - begin_time)
    begin_time = 0
    end_time = 0
    
def cal_component_visit(component_type, name):
    if component_type == -1:
        return
    elif component_type == 1:
        if name not in component_alias:
            component_types['Button'] = component_types['Button'] + 1
            temp_name = 'Button_' + str(component_types['Button'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 2:
        if name not in component_alias:
            component_types['Combo'] = component_types['Combo'] + 1
            temp_name = 'Combo_' + str(component_types['Combo'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 3:
        if name not in component_alias:
            component_types['Text'] = component_types['Text'] + 1
            temp_name = 'Text_' + str(component_types['Text'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 4:
        if name not in component_alias:
            component_types['Spin'] = component_types['Spin'] + 1
            temp_name = 'Spin_' + str(component_types['Spin'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 5:
        if name not in component_alias:
            component_types['Slider'] = component_types['Slider'] + 1
            temp_name = 'Slider_' + str(component_types['Slider'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 6:
        if name not in component_alias:
            component_types['Calendar'] = component_types['Calendar'] + 1
            temp_name = 'Calendar_' + str(component_types['Calendar'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 7:
        if name not in component_alias:
            component_types['Lcd'] = component_types['Lcd'] + 1
            temp_name = 'Lcd_' + str(component_types['Lcd'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 8:
        if name not in component_alias:
            component_types['Progress'] = component_types['Progress'] + 1
            temp_name = 'Progress_' + str(component_types['Progress'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 9:
        if name not in component_alias:
            component_types['List'] = component_types['List'] + 1
            temp_name = 'List_' + str(component_types['List'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 10:
        if name not in component_alias:
            component_types['Tree'] = component_types['Tree'] + 1
            temp_name = 'Tree_' + str(component_types['Tree'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 11:
        if name not in component_alias:
            component_types['Table'] = component_types['Table'] + 1
            temp_name = 'Table_' + str(component_types['Table'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 12:
        if name not in component_alias:
            component_types['Column'] = component_types['Column'] + 1
            temp_name = 'Column_' + str(component_types['Column'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 13:
        if name not in component_alias:
            component_types['Action'] = component_types['Action'] + 1
            temp_name = 'Action_' + str(component_types['Action'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    elif component_type == 14:
        if name not in component_alias:
            component_types['Container'] = component_types['Container'] + 1
            temp_name = 'Container_' + str(component_types['Container'])
            component_names[temp_name] = 0
            component_alias[name] = temp_name
    alias = component_alias[name]
    component_names[alias] = component_names[alias] + 1 
    if last_time == 0 or abs(last_time - current_time) > 1000:
        return
    elif alias not in component_times:
        component_times[alias] = current_time - last_time
    else:
        component_times[alias] = current_time - last_time + component_times[alias]
   
        
        
        
    
    
#读取文件  dir_path换成csv文件目录
dir_path = 'D:\\Download\\wechat\\WeChat Files\\wxid_9mvzfsky1uok22\\FileStorage\\File\\2023-03\\data_testcomponent'
files = os.listdir(dir_path)
csv_files = [os.path.join(dir_path, f) for f in files if f.endswith('.csv')]
data = pd.concat([pd.read_csv(f) for f in csv_files], ignore_index = True)

#程序单次运行
single_running_time = []
begin_time = 0
end_time = 0
last_time = 0
current_time = 0

#鼠标点击，鼠标移动，按键次数记录
mouse_click_cnt = 0
mouse_move_cnt = 0
key_press_cnt = 0

#组件访问记录
component_types = { 'Button' : 0, 'Combo' : 0, 'Text' : 0, 'Spin' : 0,\
                  'Slider' : 0, 'Calendar' : 0, 'Lcd' : 0, 'Progress' : 0,\
                      'List' : 0, 'Tree' : 0, 'Table' : 0, 'Column' : 0, \
                       'Action' : 0, 'Container' : 0 }
component_names = {}
component_alias = {}
component_times = {}




for index, row in data.iterrows():
    current_time = row[1]
    if row[0] == 1:  #程序启动
        begin_time = row[1]
    elif row[0] == 2:  #程序退出
        end_time = row[1]
        single_running(begin_time, end_time)
    elif row[0] == 3:  #鼠标点击
        mouse_click_cnt += 1
    elif row[0] == 4:  #鼠标移动
        mouse_move_cnt += 1
    elif row[0] == 5:
        key_press_cnt += 1
        
    if not np.isnan(row[-2]):  ##组件访问
        cal_component_visit(row[-2], row[-3])
    last_time = row[1]
        

plt.rcParams['font.family'] = 'SimHei'
plt.rcParams['axes.unicode_minus'] = False
fig, axs = plt.subplots(2,2)
fig.set_size_inches(12, 12)

#单次使用时常统计
axs[0, 0].set_title('每次打开应用使用时长')
axs[0, 0].set_ylabel('单次使用时间(ms)')
axs[0, 0].set_xticks([])
axs[0, 0].bar(range(len(single_running_time)),single_running_time)

#事件占比统计
axs[0, 1].set_title('鼠标键盘事件占比')
event_datas = [mouse_click_cnt, mouse_move_cnt, key_press_cnt]
event_labels = ['鼠标点击', '鼠标移动', '键盘点击']
explode = (0, 0.1, 0)
axs[0, 1].pie(event_datas, labels = event_labels, explode = explode, autopct = '%1.1f%%')

#组件访问次数统计
sort_degree = sorted(component_names.items(), key = lambda x: x[1], reverse = True)
com_name = []
com_degree = []
for com in sort_degree:
    com_name.append(com[0])
    com_degree.append(com[1])
axs[1, 0].set_title('组件访问次数（降序）')
axs[1, 0].set_ylabel('组件访问次数')
axs[1, 0].set_xticklabels(com_name[ : 10],rotation = 270, fontsize = 8)
axs[1, 0].bar(com_name[ : 10], com_degree[ : 10])

#组件访问时常统计
sort_time = sorted(component_times.items(), key = lambda x: x[1], reverse = True)
com_name_t = []
com_time = []
for com in sort_time:
    com_name_t.append(com[0])
    com_time.append(com[1])
axs[1, 1].set_title('组件访问时间（降序）')
axs[1, 1].set_ylabel('组件访问时间(ms)')
axs[1, 1].set_xticklabels(com_name_t[ : 10],rotation = 270, fontsize = 8)
axs[1, 1].bar(com_name_t[ : 10], com_time[ : 10])
plt.show()
        