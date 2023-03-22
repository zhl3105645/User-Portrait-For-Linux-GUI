from cProfile import label
from cmath import nan
import os
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
import yaml


class EventType:
    AppStart = 1
    AppQuit = 2
    MouseClick = 3
    MouseMove = 4
    KeyClick = 5


class MouseMoveType:
    Begin = 1
    End = 2


class ComponentType:
    Unknow = -1
    Button = 1
    Combo = 2
    Text = 3
    Spin = 4
    Slider = 5
    Calendar = 6
    Lcd = 7
    Progress = 8
    List = 9
    Tree = 10
    Table = 11
    Column = 12
    Action = 13
    Container = 14


# 全局变量

component_type_num = { #  类型数量
    ComponentType.Unknow : 0,
    ComponentType.Button : 0,
    ComponentType.Combo : 0,
    ComponentType.Text : 0,
    ComponentType.Spin : 0,
    ComponentType.Slider : 0,
    ComponentType.Calendar : 0,
    ComponentType.Progress : 0,
    ComponentType.Lcd : 0,
    ComponentType.List : 0,
    ComponentType.Tree : 0,
    ComponentType.Table : 0,
    ComponentType.Column : 0,
    ComponentType.Action : 0,
    ComponentType.Container : 0,
}

component_type_to_str = {  # 中文含义
    ComponentType.Unknow: 'Unknow',
    ComponentType.Button : 'Button',
    ComponentType.Combo : 'Combo',
    ComponentType.Text : 'Text',
    ComponentType.Spin : 'Spin',
    ComponentType.Slider : 'Slider',
    ComponentType.Calendar : 'Calendar',
    ComponentType.Lcd: 'Lcd',
    ComponentType.Progress : 'Progress',
    ComponentType.List : 'List',
    ComponentType.Tree : 'Tree',
    ComponentType.Table : 'Table',
    ComponentType.Column : 'Column',
    ComponentType.Action : 'Action',
    ComponentType.Container : 'Container',
}
component_name_to_front_name_map = {} # 映射到图像的名字

 
data_dir_path = ''   # 数据目录
no_event_interval = 0  # 允许最大事件时间间隔 单位 s

app_run_time = []  # 应用使用时间 单位s
app_operate_time = []  # 应用使用操作时间
app_not_operate_time = []  # 应用使用未操作时间

event_num = []  # 事件次数
event_mouse_click_num = [] 
event_mouse_move_num = []
event_key_click_num = []

key_value_map = []  # 键码点击次数。元素是map，存储 key_value -> key_click_num
all_key_value_map = {} # 键码点击次数

component_name_map = []  # 鼠标点击组件次数。元素是map，存储 component_name -> key_click_num
all_component_name_map = {} 

def read_config():
    global data_dir_path
    global no_event_interval
    with open('./conf.yaml', 'r', encoding='utf-8') as f:
        config = yaml.load(f.read(), Loader=yaml.FullLoader)
        data_dir_path = config['data_dir_path']
        no_event_interval = config['no_event_interval']
    return 


def data_process(file: str):
    # 全局变量
    global component_name_to_front_name_map
    global component_type_num
    global component_type_to_str

    global no_event_interval

    global app_run_time
    global app_operate_time
    global app_not_operate_time

    global event_num
    global event_mouse_click_num
    global event_mouse_move_num
    global event_key_click_num

    global key_value_map
    global component_name_map

    # 局部变量
    begin_time = 0
    end_time = 0

    not_operate_time = 0  # 未操作时间

    mouse_click_num = 0
    mouse_move_num = 0
    key_click_num = 0

    key_value = {}
    component_name = {}

    # 中间变量
    mouse_is_moving = False  # 鼠标是否在移动
    pre_time = 0  # 上一次有效事件的时间

    # 处理
    data = pd.read_csv(file)
    for _, row in data.iterrows():
        cur_time = int(row[1])
        if cur_time <= 0:
            continue  # 过滤非法数据

        typ = int(row[0])
        if typ == EventType.AppStart:
            begin_time = cur_time
        elif typ == EventType.AppQuit:
            end_time = cur_time
        elif typ == EventType.MouseClick:
            com_name = str(row[8])
            if com_name == '':
                continue
            com_typ = 0
            if row[9] != row[9]: # NaN != NaN
                continue
            com_typ = int(row[9])
            com_extra = str(row[10])
            if row[10] != row[10]:
                com_extra = ''

            # 组件类型数目
            component_type_num[com_typ] = component_type_num[com_typ] + 1
            # 生成图像名字
            if com_name not in component_name_to_front_name_map.keys():
                if com_extra != '':
                    component_name_to_front_name_map[com_name] = component_type_to_str[com_typ] +': ' + com_extra
                else:
                    component_name_to_front_name_map[com_name] = component_type_to_str[com_typ] + str(component_type_num[com_typ])

            # 图像名字
            front_name = component_name_to_front_name_map[com_name]
            if front_name in component_name.keys():
                component_name[front_name] = int(component_name[front_name]) + 1
            else:
                component_name[front_name] = 1

            mouse_click_num += 1
        elif typ == EventType.MouseMove:
            moving_typ = int(row[5])
            if moving_typ == MouseMoveType.Begin and mouse_is_moving is False:
                mouse_is_moving = True
            elif moving_typ == MouseMoveType.End and mouse_is_moving is True:
                mouse_is_moving = False
                mouse_move_num += 1
            else:
                continue
        elif typ == EventType.KeyClick:
            com_name = str(row[7])
            if com_name == '':
                continue
            if com_name in key_value.keys():
                key_value[com_name] = int(key_value[com_name]) + 1
            else:
                key_value[com_name] = 1
            key_click_num += 1

        # 感知未操作时间
        if no_event_interval != 0 and pre_time != 0 and (cur_time - pre_time > no_event_interval * 1000):
            not_operate_time = not_operate_time + cur_time - pre_time
        # 更新
        pre_time = cur_time
    
    # 更新
    if end_time == 0 or begin_time == 0 or end_time < begin_time:
        return

    app_run_time.append((end_time - begin_time) / 1000)
    app_operate_time.append((end_time - begin_time - not_operate_time) / 1000)
    app_not_operate_time.append(not_operate_time / 1000)

    event_key_click_num.append(key_click_num)
    event_mouse_click_num.append(mouse_click_num)
    event_mouse_move_num.append(mouse_move_num)
    event_num.append(key_click_num + mouse_click_num + mouse_move_num)

    key_value_map.append(key_value)
    component_name_map.append(component_name)
    
    return


def data_aggr():
    # 全局变量
    global key_value_map
    global all_key_value_map
    global component_name_map
    global all_component_name_map

    # 聚合
    for map in key_value_map:
        for k, v in dict(map).items():
            if k in all_key_value_map.keys():
                all_key_value_map[k] = all_key_value_map[k] + int(v)
            else:
                all_key_value_map[k] = int(v)

    for map in component_name_map:
        for k, v in dict(map).items():
            if k in all_component_name_map.keys():
                all_component_name_map[k] = all_component_name_map[k] + int(v)
            else:
                all_component_name_map[k] = int(v)

    

    return


def draw():
    # 全局变量
    global app_run_time
    global app_operate_time
    global app_not_operate_time

    global event_num
    global event_mouse_click_num
    global event_mouse_move_num
    global event_key_click_num

    global all_key_value_map
    global all_component_name_map

    # 开始绘制
    plt.rcParams['font.family'] = 'SimHei'
    plt.rcParams['axes.unicode_minus'] = False

    fig, axs = plt.subplots(2, 2)
    #fig.set_size_inches(12, 12)

    # 单次使用时长统计
    axs[0, 0].set_title('run time pre run')
    axs[0, 0].set_ylabel('time(s)')
    x1_pos = np.arange(0, 2 * len(app_run_time), 2)
    width = 0.6
    axs[0, 0].bar(x1_pos - width/2, app_operate_time, width=width, label="用户操作时间")
    axs[0, 0].bar(x1_pos + width/2, app_not_operate_time, width=width, label="用户未操作时间")
    axs[0, 0].set_xticks(x1_pos)
    axs[0, 0].set_xticklabels(range(len(app_run_time)))
    axs[0, 0].legend()

    # 事件次数统计
    axs[0, 1].set_title('event count pre run')
    axs[0, 1].set_ylabel('count')
    x2_pos = np.arange(0, 3 * len(event_num), 3)
    width = 0.8
    axs[0, 1].bar(x2_pos - width, event_key_click_num, width=width, label='key click event')
    axs[0, 1].bar(x2_pos, event_mouse_click_num, width=width, label='mouse click event')
    axs[0, 1].bar(x2_pos + width, event_mouse_move_num, width=width, label='mouse move event')
    axs[0, 1].set_xticks(x2_pos)
    axs[0, 1].set_xticklabels(range(len(event_num)))
    axs[0, 1].legend()

    # 键盘点击键统计
    axs[1, 0].set_title('key_code click count(desc)')
    axs[1, 0].set_ylabel('count')
    sort1 = sorted(all_key_value_map.items(), key = lambda x: x[1], reverse=True)
    key_name = []
    key_cnt = []
    for key in sort1:
        key_name.append(key[0])
        key_cnt.append(key[1])
    axs[1, 0].set_xticklabels(key_name[:10],rotation = 270)
    axs[1, 0].bar(key_name[ : 10], key_cnt[:10])
    axs[1, 0].legend()

    # 鼠标点击组件统计
    axs[1, 1].set_title('component click count(desc)')
    axs[1, 1].set_ylabel('count')
    sort2 = sorted(all_component_name_map.items(), key = lambda x: x[1], reverse=True)
    com_name = []
    com_cnt = []
    for com in sort2:
        com_name.append(com[0])
        com_cnt.append(com[1])
    axs[1, 1].set_xticklabels(com_name[ : 10],rotation = 270)
    axs[1, 1].bar(com_name[ : 10], com_cnt[ : 10])
    axs[1, 1].legend()

    plt.show()
    return


def run():
    # 全局变量
    global data_dir_path
    global no_event_interval

    # 读取配置
    print("begin read config from conf.yaml")
    read_config()

    if data_dir_path == '' or no_event_interval == 0:
        print("config read failed.")
        return

    # 读取文件数据
    files = os.listdir(data_dir_path)
    csv_files = [os.path.join(data_dir_path, f) for f in files if f.endswith('.csv')]

    # 数据处理
    print("begin process data")
    for file in csv_files:
        data_process(file)

    # 数据聚合
    print("begin data aggr")
    data_aggr()

    # 图像绘制
    print("begin draw")
    draw()
    return


if __name__ == "__main__":
    run()
